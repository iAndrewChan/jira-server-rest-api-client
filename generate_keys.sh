#!/bin/bash

# keypair and cert are for Oauth authentication

# generate private and public key pair
openssl genrsa -out jira_privatekey.pem 1024

# I want to fill in the information automatically, I need to look at the openssl req manual
openssl req \
    -newkey rsa:1024 \
    -x509 \
    -key jira_privatekey.pem \
    -out jira_publickey.cer \
    -days 365 \
    -config openssl_conf.ini

# extract private key
openssl pkcs8 -topk8 -nocrypt -in jira_privatekey.pem -out jira_privatekey.pcks8

# extract public key from cert
openssl x509 -pubkey -noout -in jira_publickey.cer  > jira_publickey.pem
