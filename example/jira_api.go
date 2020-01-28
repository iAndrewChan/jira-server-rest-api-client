package main

import (
	"crypto/tls"
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
	const issueurl = "/rest/api/2/issue"
	const projectKey = "EX"

	// accept any certificate, for TESTING only
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	issueid := "10100"
	accptr := &jira.Account{
		Username: "user",
		Password: "user",
		Client:   &http.Client{Transport: transport},
	}

	si := jira.BasicIssue{
		URL:         host + issueurl,
		ProjectKey:  projectKey,
		Summary:     "Summary example",
		Description: "Description example",
		IssueType:   "Story",
	}

	jira.CreateIssue(si, accptr)

	ic := jira.IssueCreate{
		Issue: jira.BasicIssue{
			URL:         host + issueurl,
			ProjectKey:  projectKey,
			Summary:     "Summary example",
			Description: "Description example",
			IssueType:   "Story",
		},
		// Add fields for new issue
		Fields: map[string]interface{}{
			"customfield_10208": jira.CFSingleChoiceList{Value: "priority"},
			"labels":            []string{"l1", "l2"},
		},
		Debug: false,
	}

	jira.CreateIssue(ic, accptr)

	iu := jira.IssueU{
		URL: host + issueurl + "/" + issueid,
		// Add fields to change
		Fields: map[string]interface{}{
			"summary":           "new summary",
			"description":       "new description",
			"customfield_10208": jira.CFSingleChoiceList{Value: "priority"},
		},
	}

	iu.UpdateIssueStrict(accptr)

	issueClone := jira.IssueClone{
		URL:              host + issueurl,
		SourceIssueURL:   host + issueurl + "/" + issueid,
		TargetProjectKey: "EX2",
		Debug:            false,
	}

	issueClone.Clone(accptr)
}
