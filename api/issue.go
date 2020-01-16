package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Account - user details and REST API request config
type Account struct {
	Username string
	Password string
	Issue    Issue
}

// Issue - REST API calls to jira issue operations
type Issue struct {
	URL         string
	ProjectKey  string
	Summary     string
	Description string
	IssueType   string
}

// validate - when an invalid issuetype value is given, panic
func (i Issue) validate() {

	switch i.IssueType {
	case "Bug":
	case "Story":
	default:
		panic("invalid value " + i.IssueType)
	}
}

type requestBody struct {
	Fields issueJSON `json:"fields"`
}

type issueJSON struct {
	Project     project   `json:"project"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	IssueType   issueType `json:"issuetype"`
}

type project struct {
	Key string `json:"key"`
}

type issueType struct {
	Name string `json:"name"`
}

// buildJSONPayload - build a json request object from Issue
func buildJSONPayload(issue Issue, debug bool) string {

	var JSONPayload []byte
	var err error

	payload := requestBody{
		Fields: issueJSON{
			Project: project{
				Key: issue.ProjectKey,
			},
			Summary:     issue.Summary,
			Description: issue.Description,
			IssueType: issueType{
				Name: issue.IssueType,
			},
		},
	}

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

// CreateIssue - create jira issue
func (acc Account) CreateIssue() {

	acc.Issue.validate()
	issuePayload := buildJSONPayload(acc.Issue, false)

	method := "POST"

	payload := strings.NewReader(issuePayload)

	client := &http.Client{}
	req, err := http.NewRequest(method, acc.Issue.URL, payload)
	if err != nil {
		fmt.Println(err)
	}

	req.SetBasicAuth(acc.Username, acc.Password)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
