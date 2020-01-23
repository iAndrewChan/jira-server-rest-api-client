package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Account - basic auth account details
type Account struct {
	Username string
	Password string
}

// Issue - common to all requests for issues
type Issue interface {
	validate()
	payloadAsJSON() string
}

// IssueC common to requests that create issues
type IssueC interface {
	createIssue(accptr *Account)
}

// RequestBody - base JSON object
type RequestBody struct {
	Fields map[string]interface{} `json:"fields"`
}

// Validate - interface validate()
func Validate(i Issue) {
	i.validate()
}

// BuildJSONPayload - interface payLoadAsJSON()
func BuildJSONPayload(i Issue) string {
	return i.payloadAsJSON()
}

// CreateIssue - interface CreateIssue()
func CreateIssue(i IssueC, accptr *Account) {
	i.createIssue(accptr)
}

// BuildPayload -
func BuildPayload(i Issue) *strings.Reader {

	Validate(i)
	issuePayload := BuildJSONPayload(i)

	payload := strings.NewReader(issuePayload)
	return payload
}

// Encode payload
func Encode(payload RequestBody, debug bool) string {

	var JSONPayload []byte
	var err error

	if debug {
		// formatting the payload to have indent for pretty print
		JSONPayload, err = json.MarshalIndent(payload, "", " ")
	} else {
		JSONPayload, err = json.Marshal(payload)
	}
	if err != nil {
		panic("Error: could not marshal payload")
	}

	JSONPayloadStr := string(JSONPayload)

	if debug {
		fmt.Println("payload" + JSONPayloadStr)
	}

	return JSONPayloadStr
}

// SendRequest - Make a request to the given url with payload
func SendRequest(accptr *Account, url string, requestType string, payload *strings.Reader, debug bool) []byte {

	var req *http.Request
	var err error
	client := &http.Client{}

	if payload != nil {
		req, err = http.NewRequest(requestType, url, payload)
	} else {
		// payload nil is of type *Reader, we need to pass a stardard nil
		req, err = http.NewRequest(requestType, url, nil)
	}
	if err != nil {
		panic("Error: problem with building the payload as NewRequest" + err.Error())
	}

	req.SetBasicAuth((*accptr).Username, (*accptr).Password)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic("Error: problem with client request: " + err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	return body
}
