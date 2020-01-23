package main

import (
	jira "jira-restapi-docker/api"
)

/*
	In the jira app:

	To find users on the system:
	Setting > User management

	To find the project key:
	Project > Project Settings > Details

	To get the issue id:
	Issues > Search for issues > Select the issue > mouse over the edit
	In the corner of the browser screen you will see a `host/secure/EditIssue!default.jspa?id=xxxxx`

	To get the custom field:
	Setting > Issues > Custom Fields > Custom Field setting > Edit, the id in the url

	To get the custom field options:
	Setting > Issues > Custom Fields > Custom Field setting > Configure > Options > Edit Options
*/

func main() {

	const host = "http://localhost:8080"
	const issueurl = "/rest/api/2/issue"
	const projectKey = "EX"
	accptr := &jira.Account{"user", "user"}

	si := jira.SimpleIssue{
		URL:         host + issueurl,
		ProjectKey:  projectKey,
		Summary:     "Summary example",
		Description: "Description example",
		IssueType:   "Story",
	}

	jira.CreateIssue(si, accptr)

	const id = "10208"
	const CFOption = "priority"

	icf := jira.IssueCF{
		Issue: si,
		CFID:  "customfield_" + id,
		CF:    jira.CustomField{CFOption},
	}

	jira.CreateIssue(icf, accptr)

	issueid := "10100"
	iu := jira.IssueU{
		URL: host + issueurl + "/" + issueid,
		Fields: map[string]interface{}{
			"summary":           "new summary",
			"description":       "new description",
			"customfield_10208": map[string]string{"value": "priority"},
		},
	}

	iu.UpdateIssueStrict(accptr)

}
