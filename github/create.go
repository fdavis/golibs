package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func CreateIssue(repo, oauthToken string) error {
	title, body := editIssue("", "")
	newIssue := Issue{Title: title, Body: body}
	// FIXME: json encode title/body strings
	createJson, _ := json.Marshal(newIssue)
	createJsonReader := strings.NewReader(string(createJson[:]))
	//fmt.Printf("url params: %s\n", EditIssuesURL+repo+"/issues")
	//fmt.Printf("json params: %s\n", createJson)
	req, err := http.NewRequest("POST", EditIssuesURL+
		repo+"/issues", createJsonReader)
	if err != nil {
		return err
	}
	req.Header.Set(
		"Authorization", "token "+oauthToken)
	req.Header.Set(
		"Accept", "application/vnd.github.v3.text-match+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	//!-
	// For long-term stability, instead of http.Get, use the
	// variant below which adds an HTTP request header indicating
	// that only version 3 of the GitHub API is acceptable.
	//
	//   req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
	//   if err != nil {
	//       return nil, err
	//   }
	//   req.Header.Set(
	//       "Accept", "application/vnd.github.v3.text-match+json")
	//   resp, err := http.DefaultClient.Do(req)
	//!+

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)

	// 201 is github created
	if resp.StatusCode != 201 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		respBody := string(bodyBytes[:])
		resp.Body.Close()
		return fmt.Errorf("create issue failed: %s\nrequest url: %v\n%s",
			resp.Status, req.URL, respBody)
	}
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Body)

	resp.Body.Close()
	return nil
}

//!-
