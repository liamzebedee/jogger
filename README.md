Jogger
======

Jogger is like IPython without all the wheels. It provides a web microservice which runs Python code in a sandboxed environment. Features:

 - Docker sandbox with no networking
 - libraries can be added in libs
 - code is given a fixed timeframe to run (no malicious infinite loops)
 - code is sandboxed from networking
 - websocket communications where reading from stdin is supported

GPL v3 license.

## Docs
 - `libs` contains the libraries the Python code can access
 - `tests` contain some **manual tests** for the features listed above
 - `uploaded_code` is for temporary code that is uploaded

## Install
 - Docker
 - Go
 - Websocketd

```
brew install homebrew boot2docker
brew install go
https://github.com/joewalnes/websocketd/wiki
docker build -t jogger .
```

go build && websocketd --port=8080 --devconsole ./jogger 
go build && websocketd --port=51000 --devconsole ./jogger 


### Generate SSL certs (locally)
To use Websockets Secure (WSS), we have to generate our own self-signed certs.

```
openssl genrsa -passout pass:x -out ./tls/server.pass.key 2048
openssl rsa -passin pass:x -in ./tls/server.pass.key -out ./tls/server.key
openssl req -new -key ./tls/server.key -out ./tls/server.csr
openssl x509 -req -sha256 -days 365 -in ./tls/server.csr -signkey ./tls/server.key -out ./tls/server.crt
```

At the generate cert step, you have to specify localhost as the FQDN. **Then visit https://localhost:8080/ in the browser to say you trust the certificate** (without a CA certifying it).

### Generate SSL certs (production)
For the domain `backend.yoursite.com`:


```
# generate letsencrypt account key
openssl genrsa 4096 > tls/account.key

# create csr
openssl req -new -sha256 -key tls/server.key -subj "/CN=backend.yoursite.com" > tls/domain.csr

# challenge for letsencrypt to prove we own the domain
mkdir -p letsencrypt_challenge/.well-known/acme-challenge/

python acme-tiny/acme_tiny.py --account-key ./tls/account.key --csr ./tls/domain.csr --acme-dir $PWD/letsencrypt_challenge/.well-known/acme-challenge/ > ./tls/signed.crt
```


**Auto-renew (every 90 days):**

```
#!/usr/bin/sh
python acme-tiny/acme_tiny.py --account-key ./account.key --csr ./domain.csr --acme-dir $PWD/letsencrypt_challenge/.well-known/acme-challenge/ > ./signed.crt

wget -O - https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem > intermediate.pem
cat /tmp/signed.crt intermediate.pem > /path/to/chained.pem
service nginx reload
```

```
#example line in your crontab (runs once per month)
0 0 1 * * /path/to/renew_cert.sh 2>> /var/log/acme_tiny.log
```



