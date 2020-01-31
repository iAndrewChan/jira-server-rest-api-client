package api

import "fmt"

// BasicIssue - REST API calls to jira issue operations
type BasicIssue struct {
	URL         string
	ProjectKey  string
	Summary     string
	Description string
	IssueType   string
	Debug       bool
}

// Project -
type Project struct {
	Key string `json:"key"`
}

// IssueType -
type IssueType struct {
	Name string `json:"name"`
}

// validate - when an invalid issuetype value is given, panic
func (i BasicIssue) validate() {

	switch i.IssueType {
	case "Bug":
	case "Story":
	default:
		panic("invalid value " + i.IssueType)
	}
}

func (i BasicIssue) payloadAsJSON() string {

	issue := map[string]interface{}{
		"project":     Project{i.ProjectKey},
		"summary":     i.Summary,
		"description": i.Description,
		"issuetype":   IssueType{i.IssueType},
	}

	payload := RequestBody{IssueRequestBody{issue}}

	return Encode(payload, i.Debug)
}

// createIssue makes a basic jira issue
func (i BasicIssue) createIssue(accptr *Account) {

	payload := BuildPayload(i)

	if i.Debug {
		fmt.Println(">>>>Payload>>>>")
		fmt.Println(payload)
		fmt.Println("<<<<Payload<<<<")
	}
	SendRequest(accptr, i.URL, "POST", payload)
}
