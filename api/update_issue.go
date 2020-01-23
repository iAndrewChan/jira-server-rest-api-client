package api

import (
	"encoding/json"
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
	payload := RequestBody{i.Fields}
	return Encode(payload, i.Debug)
}

// getIssueFields retrieves the following information
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
func (i IssueU) getIssueFields(accptr *Account) map[string]interface{} {

	v := make(map[string]interface{})

	body := SendRequest(accptr, i.URL, "GET", nil)
	err := json.Unmarshal(body, &v)
	if err != nil {
		panic("Error: unmarshalling problem " + err.Error())
	}

	fields := v["fields"].(map[string]interface{})

	return fields
}

// UpdateIssueStrict - update issue only if the fields already exist
func (i IssueU) UpdateIssueStrict(accptr *Account) {

	issueFields := i.getIssueFields(accptr)

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
