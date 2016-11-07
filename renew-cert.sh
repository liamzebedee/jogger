#!/bin/bash
LETSENCRYPT_ACCOUNT_KEY=tls/letsencrypt_account.key
CERT_SIGNING_REQ=tls/domain.csr
SIGNED_CERT=tls/signed.crt

if (( $EUID != 0 )); then
    echo "Please run as root"
    exit
fi


TIMEOUT=timeout
if !(which $TIMEOUT); then 
	# macOS
	TIMEOUT=gtimeout
fi

function start_letsencrypt_verification_server {
	cd letsencrypt_challenge
	$TIMEOUT 6 python -m SimpleHTTPServer 80
	cd ..
}

start_letsencrypt_verification_server &

python acme-tiny/acme_tiny.py --account-key $LETSENCRYPT_ACCOUNT_KEY --csr $CERT_SIGNING_REQ --acme-dir $PWD/letsencrypt_challenge/.well-known/acme-challenge/ > $SIGNED_CERT