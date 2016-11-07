#!/bin/bash
set -e
# ps -a | grep "sudo python -m SimpleHTTPServer" | head -n 1 | cut -d' ' -f2 | xargs kill -KILL
# sudo lsof -i tcp:80 | grep "Python" 

if (( $EUID != 0 )); then
    echo "Please run as root"
    exit
fi




DOMAIN=backend.yoursite.com

LETSENCRYPT_ACCOUNT_KEY=tls/letsencrypt_account.key
SERVER_KEY=tls/server.key
SERVER_PASS=tls/server.pass.key

CERT_SIGNING_REQ=tls/domain.csr
SIGNED_CERT=tls/signed.crt

# generate letsencrypt account key
openssl genrsa 2048 > $LETSENCRYPT_ACCOUNT_KEY

# generate server key
openssl genrsa -passout pass:x -out $SERVER_PASS 2048
openssl rsa -passin pass:x -in $SERVER_PASS -out $SERVER_KEY

# create csr
openssl req -new -sha256 -key $SERVER_KEY -subj "/CN=$DOMAIN" > $CERT_SIGNING_REQ

# challenge for letsencrypt to prove we own the domain
mkdir -p letsencrypt_challenge/.well-known/acme-challenge/



./renew-cert.sh




# wget -O - https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem > intermediate.pem
# cat /tmp/signed.crt intermediate.pem > /path/to/chained.pem
# service nginx reload


# 0 0 1 * * /path/to/renew_cert.sh 2>> /var/log/acme_tiny.log