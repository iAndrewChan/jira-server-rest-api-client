package api

import "fmt"

/*
	In the GUI I found two ways:
		* Project > Select issue > More action > Clone
		* Issue > Search for issues > More > Clone

	Cloning an issue is not supported in the REST api
	https://community.atlassian.com/t5/Answers-Developer-Questions/Jira-REST-Api-Clone-an-Issue/qaq-p/561045
	https://community.atlassian.com/t5/Jira-questions/Clone-Issue-in-java-with-rest-api/qaq-p/932948

	Fields had to delete:

	"timespent"
	"timeoriginalestimate"
	"aggregatetimespent
	"resolution"
	"timetracking"
	"customfield_10105" (rank field)
	"customfield_10106" (story points)
	"attachment
	"aggregatetimeestimate"
	"resolutiondate"
	"workratio"
	"lastViewed"
	"watches"
	"creator"
	**"subtasks"**
	"created"
	"customfield_10000" (development summary)
	"aggregateprogress"
	"environment"
	"timeestimate"
	"aggregatetimeoriginalestimate"
	"versions"
	*"duedate"*
	"progress"
	**"comment"**
	*"issuelinks"*
	*"votes"*
	*"worklog"*
	"updated"
	*"status"*
	reporter

	---

	customfield_10104 (sprint field)
		* https://community.atlassian.com/t5/Agile-questions/Jira-How-can-I-find-out-the-sprint-ID-of-a-particular-sprint/qaq-p/408823
		* In the GUI: Project > Backlog > mouse over the 3 dots in the sprint section.
		* At the corner of the screen `url/sprintId=x`

	Fields could keep:
	issuetype
	components
	description
	project
	fixVersions
	custom field 10100
	customfields (created by user)
	summary
	priority
	labels
	assignee
*/

// IssueClone used for cloning a issue
type IssueClone struct {
	URL              string
	SourceIssueURL   string
	TargetProjectKey string
	Debug            bool
}

// Clone from an issue to the specified project
func (i IssueClone) Clone(accptr *Account) {

	excludeFields := [36]string{
		"customfield_10000", "customfield_10105", "customfield_10106", "customfield_10104",
		"timespent", "timeoriginalestimate", "aggregatetimespent", "resolution", "timetracking",
		"attachment", "aggregatetimeestimate", "resolutiondate", "reporter",
		"workratio", "lastViewed", "watches", "creator", "subtasks", "project",
		"created", "aggregateprogress", "environment", "timeestimate",
		"aggregatetimeoriginalestimate", "versions", "duedate",
		"progress", "comment", "issuelinks", "votes", "worklog",
		"updated", "status", "summary", "description", "issuetype",
	}

	issueFields := GetIssueFields(accptr, i.SourceIssueURL)

	summary := issueFields["summary"]
	description := issueFields["description"]
	issueTypeM := issueFields["issuetype"].(map[string]interface{})
	issueType := issueTypeM["name"]

	filter(&issueFields, excludeFields)

	if i.Debug {
		fmt.Println(">>>>Issue fields kept>>>>")
		for k := range issueFields {
			fmt.Println(k)
		}
		fmt.Println("<<<<Issue fields kept<<<<")
	}

	ic := IssueCreate{
		Issue: BasicIssue{
			URL:         i.URL,
			ProjectKey:  i.TargetProjectKey,
			Summary:     summary.(string),
			Description: description.(string),
			IssueType:   issueType.(string),
		},
		// Add fields for new issue
		Fields: issueFields,
		Debug:  true,
	}

	CreateIssue(ic, accptr)

}

func filter(input *map[string]interface{}, exclude [36]string) {
	for _, field := range exclude {
		delete(*input, field)
	}
}
