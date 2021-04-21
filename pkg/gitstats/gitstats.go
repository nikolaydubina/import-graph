package gitstats

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type gitLogFetcher interface {
	GetGitLog(gitURL url.URL) (GitLog, error)
}

// GitStatsFetcher computes git stats after fetching using provided storage
type GitStatsFetcher struct {
	GitLogFetcher gitLogFetcher
}

// GitStats contains information about single git repository computed using local git only
type GitStats struct {
	LastCommit          time.Time `json:"last_commit,omitempty"`
	DaysSinceLastCommit float64   `json:"last_commit_days_since"`
	NumContributors     uint      `json:"num_contributors"`
}

func (g *GitStatsFetcher) GetGitStats(gitURL url.URL) (*GitStats, error) {
	logs, err := g.GitLogFetcher.GetGitLog(gitURL)
	if err != nil {
		return nil, fmt.Errorf("can not get git logs: %w", err)
	}
	if len(logs) == 0 {
		return nil, errors.New("git log is empty")
	}

	stats := GitStats{
		LastCommit:          logs[0].AuthorDate,
		DaysSinceLastCommit: logs.DaysSinceLastCommit(),
		NumContributors:     logs.NumContributors(),
	}
	return &stats, nil
}
