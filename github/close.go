package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const closeJson string = "{\"state\": \"closed\"}"

func CloseIssue(repo, number, oauthToken string) error {
	closeJsonReader := strings.NewReader(closeJson)
	req, err := http.NewRequest("PATCH", CloseIssuesURL+
		repo+"/issues/"+number, closeJsonReader)
	// get var only works for search, not edit +"?state=closed",
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
		return fmt.Errorf("close issue %s failed: %s\nrequest url: %v\n%s",
			number, resp.Status, req.URL, respBody)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)

	resp.Body.Close()
	return nil
}

//!-
