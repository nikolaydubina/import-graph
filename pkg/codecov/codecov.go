// Package codecov is a client to interact with codecov.io
package codecov

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nikolaydubina/import-graph/pkg/github"
)

// HTTPClient is codecov HTTP based client
type HTTPClient struct {
	BaseURL    string // e.g. "api.codecov.io"
	HTTPClient *http.Client
}

// Totals in HTTP response
type Totals struct {
	NumFiles uint    `json:"files"`
	NumLines uint    `json:"lines"`
	Coverage float64 `json:"coverage"`
}

// Report in HTTP response
type Report struct {
	Totals Totals `json:"totals"`
}

// CommitStats in HTTP response
type CommitStats struct {
	Report Report `json:"report"`
}

// RepoStats codecov stats about single repo
type RepoStats struct {
	Language      string       `json:"language"`      // e.g. "go"
	Branch        string       `json:"branch"`        // e.g. "main"
	Name          string       `json:"name"`          // name of repository
	LatestCommmit *CommitStats `json:"latest_commit"` // can be null for repos registered but no data yet
	RepoURL       url.URL      `json:"-"`             // computed
}

func getRepoURL(owner string, repoName string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("https://app.codecov.io/gh/%s/%s", owner, repoName))
}

// GetRepoStats makes HTTP call to codecov and parses response
func (c HTTPClient) GetRepoStats(owner string, repoName string) (*RepoStats, error) {
	if owner == "" || repoName == "" {
		return nil, errors.New("owner or repo is empty stirng")
	}
	url := fmt.Sprintf("https://%s/internal/github/%s/repos/%s/", c.BaseURL, owner, repoName)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can not get make request: %w", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("can not read body: %w", err)
	}

	var stats RepoStats
	if err := json.Unmarshal(buf.Bytes(), &stats); err != nil {
		return nil, fmt.Errorf("can not unmarshal response: %w", err)
	}
	if rURL, err := getRepoURL(owner, repoName); err == nil && rURL != nil {
		stats.RepoURL = *rURL
	}

	return &stats, nil
}

// GetRepoStatsFromGitHubURL is convenience wrapper
func (c HTTPClient) GetRepoStatsFromGitHubURL(ghURL url.URL) (*RepoStats, error) {
	owner, repo := github.ParseGitHubURL(ghURL)
	return c.GetRepoStats(owner, repo)
}
