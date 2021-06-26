package gourlresolver

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GoURLResolver find Git and GitHub URLs for a Go module
type GoURLResolver struct {
	HTTPClient *http.Client
}

// ResolveGitHubURL finds GitHub URL
func (c GoURLResolver) ResolveGitHubURL(name string) (url.URL, error) {
	if strings.HasPrefix(name, "github.com/") {
		return resolvePointerURL(url.Parse("https://" + normalizeGitURLPath(name)))
	}
	resp, err := c.fetchData(name)
	if err != nil {
		return url.URL{}, fmt.Errorf("can not make GET to Go module name: %w", err)
	}
	gitURL, err := parseResponse(resp)
	if err != nil {
		return url.URL{}, fmt.Errorf("can not parse response: %w", err)
	}
	if gitURL.Host != "github.com" {
		return url.URL{}, fmt.Errorf("git is not on GitHub: %v", gitURL)
	}
	return *gitURL, nil
}

// ResolveGitURL finds git URL
func (c GoURLResolver) ResolveGitURL(name string) (url.URL, error) {
	if strings.HasPrefix(name, "github.com/") {
		return resolvePointerURL(url.Parse("https://" + normalizeGitURLPath(name)))
	}
	resp, err := c.fetchData(name)
	if err != nil {
		return url.URL{}, fmt.Errorf("can not make GET to Go module name: %w", err)
	}
	return resolvePointerURL(parseResponse(resp))
}

func normalizeGitURLPath(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join(parts[:3], "/")
}

func (c GoURLResolver) fetchData(name string) (string, error) {
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

// convenience function
func resolvePointerURL(u *url.URL, err error) (url.URL, error) {
	if u != nil {
		return *u, err
	}
	return url.URL{}, err
}
