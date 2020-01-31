package api

import (
	"bytes"
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
	Client   *http.Client
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

// IssueRequestBody - base issue JSON object
type IssueRequestBody struct {
	Fields map[string]interface{} `json:"fields"`
}

// RequestBody abstraction for different request body structs
type RequestBody struct {
	body interface{}
}

// Validate - interface validate()
func Validate(i Issue) {
	i.validate()
}

// CreateIssue - interface CreateIssue()
func CreateIssue(i IssueC, accptr *Account) {
	i.createIssue(accptr)
}

// BuildPayload -
func BuildPayload(i Issue) *strings.Reader {

	Validate(i)
	issuePayload := i.payloadAsJSON()

	payload := strings.NewReader(issuePayload)
	return payload
}

// Encode payload to JSON
func Encode(req RequestBody, debug bool) string {

	payload := req.body

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
		fmt.Println(">>>>JSON payload string>>>>")
		fmt.Println(JSONPayloadStr)
		fmt.Println("<<<<JSON payload string<<<<")
	}

	return JSONPayloadStr
}

// SendRequest - Make a request to the given url with payload
func SendRequest(accptr *Account, url string, requestType string, payload *strings.Reader) []byte {

	var req *http.Request
	var err error

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

	res, err := (*accptr).Client.Do(req)
	if err != nil {
		panic("Error: problem with client request: " + err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var prettyJSONBody bytes.Buffer
	json.Indent(&prettyJSONBody, body, "", " ")

	fmt.Println(">>>>Response>>>>")
	fmt.Println(string(prettyJSONBody.Bytes()))
	fmt.Println("<<<<Response<<<<")

	return body
}

// GetIssueFields retrieves the following information
// * issue id
// * issue key
// * self - url to itself
// * fields: the field key contains all the fields in an issue, for example:
// 	* customfields_*
// 	* summary
// 	* subtask etc.
//
// Custom fields exist in issuetype Epic:
// * customfield_10100 - Epic to assign the issue to
// * customfield_10101 - Epic status
// * customfield_10102 - (short) name of the epic
// * customfield_10103 - Epic colour field
//
// Custom fields exst in all issuetypes:
// * customfield_10000 - Development summary
// * customfield_10104 - sprint field
// * customfield_10105 - rank field
// * customfield_10106 - story points (exists in issuetype epics and story)
func GetIssueFields(accptr *Account, url string) map[string]interface{} {

	v := make(map[string]interface{})

	body := SendRequest(accptr, url, "GET", nil)
	err := json.Unmarshal(body, &v)
	if err != nil {
		panic("Error: unmarshalling problem " + err.Error())
	}

	fields := v["fields"].(map[string]interface{})

	return fields
}
