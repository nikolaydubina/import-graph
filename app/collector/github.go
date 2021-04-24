package collector

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-github/v35/github"

	iggithub "github.com/nikolaydubina/import-graph/pkg/github"
)

type GitHubSummarizer struct {
	GitHubClient *github.Client
}

func (c *GitHubSummarizer) GetSummary(ctx context.Context, ghURL url.URL) (*GitHubSummary, error) {
	owner, repo := iggithub.ParseGitHubURL(ghURL)
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

type GitHubSummary struct {
	NumStartsRepo *int `json:"github_repo_stars,omitempty"`
}
