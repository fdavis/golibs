package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getIssueTitleBody(repo, number string) (string, string, error) {
	req, err := http.NewRequest("GET", EditIssuesURL+
		repo+"/issues/"+number, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set(
		"Accept", "application/vnd.github.v3.text-match+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		respBody := string(bodyBytes[:])
		resp.Body.Close()
		return "", "", fmt.Errorf("get issue failed: %s\nrequest url: %v\n%s",
			resp.Status, req.URL, respBody)
	}
	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return "", "", err
	}
	resp.Body.Close()
	return result.Title, result.Body, nil
}

const editJsonTemplate string = "{\"title\": \"%s\", \"body\": \"%s\"}"

func UpdateIssue(repo, number, oauthToken string) error {
	title, body, _ := getIssueTitleBody(repo, number)
	title, body = editIssue(title, body)

	newIssue := Issue{Title: title, Body: body}
	editJson, _ := json.Marshal(newIssue)
	editJsonReader := strings.NewReader(string(editJson[:]))

	req, err := http.NewRequest("PATCH", EditIssuesURL+
		repo+"/issues/"+number, editJsonReader)
	//fmt.Printf("url params: %s\n", EditIssuesURL+repo+"/issues")
	//fmt.Printf("json params: %s\n", editJson)
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

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		respBody := string(bodyBytes[:])
		resp.Body.Close()
		return fmt.Errorf("edit issue failed: %s\nrequest url: %v\n%s",
			resp.Status, req.URL, respBody)
	}
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Body)

	resp.Body.Close()
	return nil
}

//!-
