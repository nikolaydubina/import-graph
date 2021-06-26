package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v35/github"
)

// GitHubSummarizer collects summary about github repo
type GitHubSummarizer struct {
	GitHubClient *github.Client
}

// GetSummary collects summary about github repo
func (c *GitHubSummarizer) GetSummary(ctx context.Context, ghURL url.URL) (*GitHubSummary, error) {
	owner, repo := ParseGitHubURL(ghURL)
	ghRepo, _, err := c.GitHubClient.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("can not get github repo %s %s: %w", owner, repo, err)
	}
	if ghRepo == nil {
		return nil, errors.New("retrieved repository info is nil")
	}

	summary := GitHubSummary{
		NumStartsRepo: ghRepo.StargazersCount,
	}

	return &summary, nil
}

// GitHubSummary is used in collector result
type GitHubSummary struct {
	NumStartsRepo *int `json:"github_repo_stars,omitempty"`
}

func ParseGitHubURL(repoURL url.URL) (owner, repoName string) {
	parts := []string{}
	// Filtering out empty strings
	for _, p := range strings.Split(repoURL.EscapedPath(), "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
