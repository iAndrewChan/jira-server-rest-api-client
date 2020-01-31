package api

import (
	"fmt"
)

// IssueU update issue fields
type IssueU struct {
	URL    string
	Fields map[string]interface{}
	Debug  bool
}

func (i IssueU) validate() {}

func (i IssueU) payloadAsJSON() string {
	payload := RequestBody{IssueRequestBody{i.Fields}}
	return Encode(payload, i.Debug)
}

// UpdateIssueStrict - update issue only if the fields already exist
func (i IssueU) UpdateIssueStrict(accptr *Account) {

	issueFields := GetIssueFields(accptr, i.URL)

	for key := range i.Fields {
		if _, present := issueFields[key]; !present {
			panic("cannot update field: " + key + " does not exist")
		}
	}

	payload := BuildPayload(i)

	if i.Debug {
		fmt.Println("====Payload====")
		fmt.Println(payload)
		fmt.Println("====Payload====")
	}

	SendRequest(accptr, i.URL, "PUT", payload)
}
