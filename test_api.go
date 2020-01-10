package main

import (
		"fmt"
		"crypto/rsa"
		"crypto/x509"
		"encoding/pem"

		"github.com/dghubble/oauth1"
)

const jiraPrivateKey = `
-----BEGIN PRIVATE KEY-----
<enter private key here>
-----END PRIVATE KEY-----`

func main() {


	// create client
	var jiraHome string = "http://localhost:8080/jira"
	var jiraOauth string = jiraHome + "/plugins/servlet/oauth/"

	var authorizeEndpoint = oauth1.Endpoint{
		RequestTokenURL: jiraOauth + "request-token",
		AuthorizeURL:    jiraOauth + "authorize",
		AccessTokenURL:  jiraOauth + "access-token",
	}
	// decode PEM ASCII into PEM blocks
	block, _ := pem.Decode([]byte(jiraPrivateKey))
	if block == nil {
		panic("failed to parse PEM block")
	}

	// parse uncrypted private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse key: " + err.Error())
	}

	// get key with required type from interface
	privKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		panic("key is not type *rsa.PrivateKey")
	}

	// CallbackURL oob - out of band ?
	config := oauth1.Config{
		ConsumerKey: "OauthKey",
		ConsumerSecret: "SampleSecret",
		CallbackURL: "oob",
		Endpoint: authorizeEndpoint,
		Signer: &oauth1.RSASigner{
			PrivateKey: privKey,
		},
	}

	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		panic("Unable to get request token: " + err.Error())
	}
	fmt.Println(requestToken, requestSecret)
	
}