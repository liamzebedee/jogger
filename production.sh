#!/bin/bash
export PATH=$(pwd):$PATH
export GOPATH=$PWD
export PRODUCTION=1

PORT=8080

# Port 80 is used so we can serve the Let's Encrypt challenge
go build && websocketd  --port=$PORT --devconsole --loglevel=access --ssl --sslcert=./tls/signed.crt --sslkey=./tls/server.key ./jogger