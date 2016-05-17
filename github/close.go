package github

import (
	"fmt"
	"net/http"
)

func CloseIssue(repo, number string) error {
	req, err := http.NewRequest("PATCH", CloseIssuesURL+
		"/"+repo+"/issues/"+number+"?state=closed", nil)
	if err != nil {
		return err
	}
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
		resp.Body.Close()
		return fmt.Errorf("close issue %s failed: %s", number, resp.Status)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)

	resp.Body.Close()
	return nil
}

//!-
