package main

import (
	"fmt"
	"encoding/json"
	"os/exec"
	"io/ioutil"
	"strings"
	"os"
	"path/filepath"
	"bufio"
	"log"
	"bytes"
	"errors"
)


var logger *log.Logger
var production bool
var pwd string


const message_SendCode = "send_code"
const message_SendInput = "send_input"

const message_Ouput = "output"
const message_WaitingForInput = "waiting_for_input"
const message_EndOfOutput = "end_of_output"
const message_Error = "error"



type Message struct {
	Type string             `json:"type"`
    Body string 			`json:"body"`
}


// trimEOL cuts unixy style \n and windowsy style \r\n suffix from the string
func trimEOL(b []byte) []byte {
	lns := len(b)
	if lns > 0 && b[lns-1] == '\n' {
		lns--
		if lns > 0 && b[lns-1] == '\r' {
			lns--
		}
	}
	return b[:lns]
}


type ManagedStdout struct {
}

func (mgdStdout *ManagedStdout) Write(buf []byte) (n int, err error) {
	sendReply(Message{message_Ouput, string(trimEOL(buf))})
	return len(buf), nil
}


type ManagedStderr struct {
}

func (mgdStderr *ManagedStderr) Write(buf []byte) (n int, err error) {
	sendReply(Message{message_Error, string(buf)})
	return len(buf), nil
}

func makeDockerCmd(cmdStr string) (*exec.Cmd) {
	// Run the Docker process
	if !production {
		os.Setenv("DOCKER_TLS_VERIFY", "1")
		os.Setenv("DOCKER_HOST", "tcp://192.168.99.100:2376")
		os.Setenv("DOCKER_CERT_PATH", "/Users/liamz/.docker/machine/machines/default")
		os.Setenv("DOCKER_MACHINE_NAME", "default")
	}
	os.Setenv("PYTHONUNBUFFERED", "hellyeah!")


	// gtimeout 2
	// -u means Python is unbuffered
	// args := strings.Split("2 docker run -m 128m --ulimit nofile=10:10 --network=none -i --rm -v "+pwd+":/usr/src/myapp:ro -w /usr/src/myapp -e PYTHONPATH='./libs/' python:2.7 python -u "+codefileName, " ")
	// cmd := exec.Command("gtimeout", args...)

	cmd_pieces := strings.Split(cmdStr, " ")
	cmd := exec.Command("docker", cmd_pieces...)
	return cmd
}

func handleMessage(msg Message, input *chan string) error {
	logger.Println("Handling "+msg.Type+" message")
	if msg.Type == message_SendInput {
		*input <- msg.Body;

	} else if msg.Type == message_SendCode {
		// TODO hack to fix global
		codeWorkingDir, _ := os.Getwd()
		codeWorkingDir = codeWorkingDir+"/code-mount"

		// Write code into a temporary file
		codeFile, err := ioutil.TempFile(codeWorkingDir, "codeFile")
		defer os.Remove(codeFile.Name())
		_, codefileName := filepath.Split(codeFile.Name())

		if err != nil {
			return err
		}
		logger.Println("Creating file "+codeFile.Name())
		if _, err := codeFile.WriteString(msg.Body); err != nil {
			return err
		}
		codeFile.Sync()


		// -m 128m
		timeoutCmd := "timeout 10"
		// if !production {
		// 	timeoutCmd = "gtimeout 10"
		// }
		cmd := makeDockerCmd("run --ulimit nofile=10:10 --network=none -i --rm -v "+codeWorkingDir+":/usr/src/myapp:ro -w /usr/src/myapp -e PYTHONPATH=libs/ python:2.7 "+timeoutCmd+" python -u "+codefileName)
		cmd.Stdout = &ManagedStdout{}
		cmd.Stderr = &ManagedStderr{}

		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		inputFinished := make(chan bool)
		defer close(inputFinished)


		go func() {
			const STD_FLUSH_CHAR = '\n'
			for {
				select {
				case data := <-*input:
					stdin.Write(append([]byte(data), STD_FLUSH_CHAR))
				case <-inputFinished:
					return;
				}
			}
		}()


		err = cmd.Start()
		if err != nil {
			return err
		}

		cmd.Wait()

		reply := Message{message_EndOfOutput, ""}
		sendReply(reply)

		inputFinished <- true
	}
	return nil	
}


func sendReply(msg Message) {
	replyJson, err := json.Marshal(msg)
	
	if err == nil {
		fmt.Printf("%s\n", replyJson)
	} else {
		handleProtocolErrors("error encoding reply", err)
	}
}

func handleProtocolErrors(cause string, err error) {
	fmt.Errorf(err.Error())
	reply := Message{message_Error, cause + " " + err.Error()}
	sendReply(reply)
}



func main() {
	// if productionEnv := os.Getenv("PRODUCTION"); productionEnv == "1" {
	// 	production = true
	// } else {
	// 	production = false
	// }

	// TODO HACK big hack
	production = true
	// production = false

	logfile, err := os.Open("./log.txt") // For read access.
	if err != nil {
		fmt.Println("Can't open log.txt for logging - create it if need be")
		panic(err)
	}
	defer logfile.Close()
	logger = log.New(logfile, "jogger: ", log.Ldate | log.Ltime /*| log.Lshortfile*/)


	pwd, _ := os.Getwd()
	isDockerDaemonRunning := makeDockerCmd("run --rm -v "+pwd+":/usr/src/myapp:ro -w /usr/src/myapp python:2.7 python tests/1_sample.py")
	errBuf := new(bytes.Buffer)
	isDockerDaemonRunning.Stderr = errBuf
	isDockerDaemonRunning.Run()

	if errStr := errBuf.String(); errStr != "" {
		handleProtocolErrors("error with Docker daemon", errors.New(errStr))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	input := make(chan string)
	defer close(input)

	for scanner.Scan() {
		// Protocol messages are delimited by \n
		msgJson := scanner.Text()
		var msg Message

		err := json.Unmarshal([]byte(msgJson), &msg)
		if err != nil {
			handleProtocolErrors("error decoding message", err)
			continue
		}

		go func() {
			err = handleMessage(msg, &input)
			if err != nil {
				handleProtocolErrors("error handling message", err)
			}
		}()
	}
}


