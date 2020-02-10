package main

import (
	jira "jira-restapi-docker/api"
	"net/http"
)

/*
	In the jira app:

	Anything with `Setting >` requires admin to see and access in JIRA

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

	const host = "https://localhost:8080"
	const restAPI = "/rest/api/2/"
	const issueURL = restAPI + "issue"
	const searchURL = restAPI + "search"
	const projectKey = "EX"

	issueid := "10100"
	accptr := &jira.Account{
		Username: "user",
		Password: "user",
		Client:   &http.Client{},
	}

	bi := jira.BasicIssue{
		URL:         host + issueURL,
		ProjectKey:  projectKey,
		Summary:     "Summary example",
		Description: "Description example",
		IssueType:   "Story",
	}

	jira.CreateIssue(bi, accptr)

	ic := jira.IssueCreate{
		Issue: jira.BasicIssue{
			URL:         host + issueURL,
			ProjectKey:  projectKey,
			Summary:     "Summary example",
			Description: "Description example",
			IssueType:   "Story",
		},
		// Add fields for new issue
		Fields: map[string]interface{}{
			"customfield_10208": jira.CFSingleChoiceList{Value: "medium"},
			"labels":            []string{"l1", "l2"},
		},
		Debug: false,
	}

	jira.CreateIssue(ic, accptr)

	iu := jira.IssueU{
		URL: host + issueURL + "/" + issueid,
		// Add fields to change
		Fields: map[string]interface{}{
			"summary":           "new summary",
			"description":       "new description",
			"customfield_10208": jira.CFSingleChoiceList{Value: "priority"},
		},
	}

	iu.UpdateIssueStrict(accptr)

	issueClone := jira.IssueClone{
		URL:              host + issueURL,
		SourceIssueURL:   host + issueURL + "/" + issueid,
		TargetProjectKey: "EX2",
		Debug:            false,
	}

	issueClone.Clone(accptr)

	si := jira.SearchIssue{
		URL: host + searchURL,
		Query: map[string]interface{}{
			"startAt":    0,
			"maxResults": 5,
			"fields":     [2]string{"id", "labels"}, // fields to include
		},
		Debug: false,
	}

	si.PrintSearchResult(accptr, "l1,l2")
}
