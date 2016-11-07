#!/bin/bash
LETSENCRYPT_ACCOUNT_KEY=tls/letsencrypt_account.key
SERVER_KEY=tls/server.key
SERVER_PASS=tls/server.pass.key

CERT_SIGNING_REQ=tls/domain.csr
SIGNED_CERT=tls/signed.crt

if (( $EUID != 0 )); then
    echo "Please run as root"
    exit
fi

openssl genrsa -passout pass:x -out $SERVER_PASS 2048
openssl rsa -passin pass:x -in $SERVER_PASS -out $SERVER_KEY
openssl req -new -key $SERVER_KEY -out $CERT_SIGNING_REQ
openssl x509 -req -sha256 -days 365 -in $CERT_SIGNING_REQ -signkey $SERVER_KEY -out $SIGNED_CERT