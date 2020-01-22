package api

// SimpleIssue - REST API calls to jira issue operations
type SimpleIssue struct {
	URL         string
	ProjectKey  string
	Summary     string
	Description string
	IssueType   string
	Debug       bool
}

// validate - when an invalid issuetype value is given, panic
func (i SimpleIssue) validate() {

	switch i.IssueType {
	case "Bug":
	case "Story":
	default:
		panic("invalid value " + i.IssueType)
	}
}

func (i SimpleIssue) payloadAsJSON() string {

	issue := map[string]interface{}{
		"project":     Project{i.ProjectKey},
		"summary":     i.Summary,
		"description": i.Description,
		"issuetype":   IssueType{i.IssueType},
	}

	payload := RequestBody{issue}

	return Encode(payload, i.Debug)
}

// createIssue - create jira issue
func (i SimpleIssue) createIssue(accptr *Account) {

	payload := BuildPayload(i)
	SendRequest(accptr, i.URL, "POST", payload, false)
}
