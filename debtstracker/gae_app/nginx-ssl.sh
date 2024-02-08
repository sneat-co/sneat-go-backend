#!/usr/bin/env bash
# https://www.accuweaver.com/2014/09/19/make-chrome-accept-a-self-signed-certificate-on-osx/

# https://gist.github.com/jessedearing/2351836

# Run using "sudo"

echo "Generating an SSL private key to sign your certificate..."
openssl genrsa -des3 -out debtstracker-local.key 1024

echo "Generating a Certificate Signing Request..."
openssl req -new -key debtstracker-local.key -out debtstracker-local.csr

echo "Removing pass-phrase from key (for nginx)..."
cp debtstracker-local.key debtstracker-local.key.org
openssl rsa -in debtstracker-local.key.org -out debtstracker-local.key
rm debtstracker-local.key.org

echo "Generating certificate..."
openssl x509 -req -days 365 -in debtstracker-local.csr -signkey debtstracker-local.key -out debtstracker-local.crt

echo "Copying certificate (debtstracker-local.crt) to /etc/ssl/certs/"
mkdir -p  /etc/ssl/certs
cp debtstracker-local.crt /etc/ssl/certs/

echo "Copying key (debtstracker-local.key) to /etc/ssl/private/"
mkdir -p  /etc/ssl/private
cp debtstracker-local.key /etc/ssl/private/