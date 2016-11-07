#!/bin/bash
export PATH=$(pwd):$PATH
export GOPATH=$PWD

# Test the client
# ---------------

# go build && ./jogger


# Test the Websocket
# ------------------

go build && websocketd --address=localhost --port=8080 --devconsole --loglevel=access  --ssl --sslcert=./tls/signed.crt --sslkey=./tls/server.key ./jogger 