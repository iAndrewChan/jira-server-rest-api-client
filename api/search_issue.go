package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SearchIssue - search issues
type SearchIssue struct {
	URL   string
	Query map[string]interface{}
	Debug bool
}

func (i *SearchIssue) validate() {}

func (i *SearchIssue) payloadAsJSON() string {
	payload := RequestBody{i.Query}
	return Encode(payload, i.Debug)
}

// SearchWithLabel - search issues based on label
func (i *SearchIssue) searchWithLabel(labels string) {
	var b strings.Builder

	b.WriteString("labels in (")
	b.WriteString(labels)
	b.WriteString(")")

	i.Query["jql"] = b.String()
}

// PrintSearchResult - print the search result
func (i *SearchIssue) PrintSearchResult(accptr *Account, labels string) {

	i.searchWithLabel(labels)

	payload := BuildPayload(i)

	if i.Debug {
		fmt.Println(">>>>Payload>>>>")
		fmt.Println(payload)
		fmt.Println("<<<<Payload<<<<")
	}
	body := SendRequest(accptr, i.URL, "POST", payload)

	m := make(map[string]interface{})
	json.Unmarshal(body, &m)

}
