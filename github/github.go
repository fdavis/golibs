// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 110.
//!+

// Package github provides a Go API for the GitHub issue tracker.
// See https://developer.github.com/v3/search/#search-issues.
package github

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	IssuesURL     = "https://api.github.com/search/issues"
	EditIssuesURL = "https://api.github.com/repos/"
	editJson      = "{\"title\": \"%s\", \"body\": \"%s\"}"
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

// changed CreatedAt to pointer to fix omitempty issue with json
// http://stackoverflow.com/questions/32643815/golang-json-omitempty-with-time-time-field
type Issue struct {
	Number    int        `json:"number,omitempty"`
	HTMLURL   string     `json:"html_url,omitempty"`
	Title     string     `json:"title,omitempty"`
	State     string     `json:"state,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Body      string     `json:"body,omitempty"`
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

const (
	defaultEditContents   string = "Title: %s\nBody:\n%s"
	defaultCreateContents string = "Title: <your title here>\nBody:\n<your body here>"
)

func editIssue(title, body string) (string, string) {

	tmpfile, _ := ioutil.TempFile("", "gogit")
	myFile := tmpfile.Name()

	if title != "" && body != "" {
		s := fmt.Sprintf(defaultEditContents, title, body)
		tmpfile.WriteString(s)
	} else {
		tmpfile.WriteString(defaultCreateContents)
	}

	tmpfile.Close()

	// FIXME: use env $EDITOR instead of vim
	cmd := exec.Command("vim", myFile)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()

	tmpfile, _ = os.Open(myFile)
	data, _ := ioutil.ReadAll(tmpfile)

	tmpfile.Close()
	os.Remove(myFile)

	return extractTitleBody(string(data[:]))
}

//!-
func extractTitleBody(text string) (string, string) {
	split := strings.SplitN(text, "\n", 3)
	titleLine, body := split[0], split[2]
	title := strings.SplitN(titleLine, " ", 2)[1]
	return title, body
}
