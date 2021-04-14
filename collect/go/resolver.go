package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type GoResolver struct {
	HTTPClient *http.Client
}

func (c *GoResolver) ResolveGitHubURL(name string) (*url.URL, error) {
	if strings.HasPrefix(name, "github.com/") {
		return url.Parse("https://" + name)
	}
	resp, err := c.fetchData(name)
	if err != nil {
		return nil, fmt.Errorf("can not make GET to Go module name: %w", err)
	}
	gitURL, err := parseResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("can not parse response: %w", err)
	}
	if gitURL.Host != "github.com" {
		return nil, fmt.Errorf("git is not on GitHub: %v", gitURL)
	}
	return gitURL, nil
}

func (c *GoResolver) ResolveGitURL(name string) (*url.URL, error) {
	if strings.HasPrefix(name, "github.com/") {
		return url.Parse("https://" + name)
	}
	resp, err := c.fetchData(name)
	if err != nil {
		return nil, fmt.Errorf("can not make GET to Go module name: %w", err)
	}
	return parseResponse(resp)
}

func (c *GoResolver) fetchData(name string) (string, error) {
	resp, err := c.HTTPClient.Get("https://" + name + "?go-get=1")
	if err != nil {
		return "", fmt.Errorf("can not make GET to Go module name: %w", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

// parseResponse reads string body returned from HTTP request and extracts git URL
// Reference: https://golang.org/ref/mod#vcs-find
// Input: <html><head><meta name="go-import" content="sourcegraph.com/sqs/pbtypes git https://github.com/sqs/pbtypes"></head><body></body></html>
// Output: https://github.com/sqs/pbtypes nil
func parseResponse(resp string) (*url.URL, error) {
	idxGoImport := strings.Index(resp, "go-import")
	if idxGoImport == -1 {
		return nil, errors.New("can not find go-import metadata")
	}

	idxContent := strings.Index(resp[idxGoImport:], "content")
	if idxContent == -1 {
		return nil, errors.New("can not find content after go-import")
	}
	idxStart := idxGoImport + idxContent + len("content=\"")
	idxEnd := int(idxStart)
	for i, v := range resp[idxStart:] {
		if v == '"' {
			idxEnd = idxStart + i
			break
		}
	}

	vals := strings.Split(strings.TrimSpace(resp[idxStart:idxEnd]), " ")
	if len(vals) != 3 {
		return nil, fmt.Errorf("unexpected num of vals in string: %s", resp[idxStart:idxEnd])
	}

	if vals[1] != "git" {
		return nil, fmt.Errorf("not git repo, vcs is: %s", vals[1])
	}
	return url.Parse(vals[2])
}
