#!/usr/bin/env bash
# https://www.accuweaver.com/2014/09/19/make-chrome-accept-a-self-signed-certificate-on-osx/

# https://gist.github.com/jessedearing/2351836

# Run using "sudo"

echo "Generating an SSL private key to sign your certificate..."
openssl genrsa -des3 -out debtus-local.key 1024

echo "Generating a Certificate Signing Request..."
openssl req -new -key debtus-local.key -out debtus-local.csr

echo "Removing pass-phrase from key (for nginx)..."
cp debtus-local.key debtus-local.key.org
openssl rsa -in debtus-local.key.org -out debtus-local.key
rm debtus-local.key.org

echo "Generating certificate..."
openssl x509 -req -days 365 -in debtus-local.csr -signkey debtus-local.key -out debtus-local.crt

echo "Copying certificate (debtstracker-local.crt) to /etc/ssl/certs/"
mkdir -p  /etc/ssl/certs
cp debtus-local.crt /etc/ssl/certs/

echo "Copying key (debtstracker-local.key) to /etc/ssl/private/"
mkdir -p  /etc/ssl/private
cp debtus-local.key /etc/ssl/private/