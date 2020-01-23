package api

import "fmt"

/*
	# Making a custom field madatory on a project

	In the GUI:
		* Project > Project settings > Click required

	# Creating an issue and include a custom field

	The error could not find valid 'id' or 'value in the Parent Option object means
	REST API is looking for an object similar to `"customfield_11700": { "value": "Compliance" }` in this case.

	You should also check whether there are options added in the custom field if they exist through the GUI.
		* Setting > Issues > Custom Field > Configure > Options
		* Refer to https://developer.atlassian.com/server/jira/platform/jira-rest-api-examples/ for custom field data formats

	Example for a single choice list
		```json
		"{{customexternalfield}}": {"value": "medium"}
		"{{customexternalfield}}": {"value": "priority"}
		"{{customexternalfield}}": {"value": "low"}
		```

	I tried to create a unselected option for the custom field, but was not able to.

	# Get custom field options

	A custom field option id is unique for the option it contains

	In the GUI this can be found in:
		* Setting > Issues > Custom Field > Configure > Options

	To get the custom field id you need to:
		* Either add a custom field to an existing issue or add the custom field when creating the issue
		* In the GUI: Project > Select Issue > three dot icon > edit > configure field > tick the desired custom field
		* Make a request to get issue, and then search for the custom field id to custom field option details (id, data format)

	# Creating a custom field of type select list (single choice)

	type: select for the Select List (single choice list)
	searcherKey: multi select searcher

	https://community.atlassian.com/t5/Answers-Developer-Questions/Need-a-list-of-custom-field-types/qaq-p/467438
		* Finding more about custom field through JIRA Java API doc, mentioned in the next section

	https://docs.atlassian.com/software/jira/docs/api/8.5.3/ (latest docs)

		com.atlassian.jira.issue.customfields > CustomFieldType (interface) > "All known implementing classes"
		* we can see typenameCFType format
		* e.g. select, multiselect, number, string, text etc.

		com.atlassian.jira.issue.customfields > CustomFieldSearcher (interface) > "All known implementing classes"
		* e.g. SelectSearcher, MultiSelectSearcher, Group Picker searcher etc.
*/

// IssueCreate - Issue with custom field
type IssueCreate struct {
	Issue  BasicIssue
	Fields map[string]interface{}
	Debug  bool
}

// CFSingleChoiceList -
type CFSingleChoiceList struct {
	Value string `json:"value"`
}

func (i IssueCreate) validate() {
	i.Issue.validate()
	// TODO check possible fields in a issue for a project.
}

func (i IssueCreate) payloadAsJSON() string {

	issue := map[string]interface{}{
		"project":     Project{i.Issue.ProjectKey},
		"summary":     i.Issue.Summary,
		"description": i.Issue.Description,
		"issuetype":   IssueType{i.Issue.IssueType},
	}

	for key, value := range i.Fields {
		issue[key] = value
	}

	payload := RequestBody{issue}

	return Encode(payload, i.Debug)
}

// createIssue - create jira issue with custom field
func (i IssueCreate) createIssue(accptr *Account) {

	payload := BuildPayload(i)

	if i.Debug {
		fmt.Println(">>>>Payload>>>>")
		fmt.Println(payload)
		fmt.Println("<<<<Payload<<<<")
	}

	SendRequest(accptr, i.Issue.URL, "POST", payload)
}
