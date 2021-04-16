package gitstats

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type GitProcessStorage interface {
	Fetch(gitURL url.URL) error
	DirPath(gitURL url.URL) string
}

// GitStats contains information about single git repository computed using git only
type GitStats struct {
	URL                   url.URL   `json:"-"`
	LastCommit            time.Time `json:"last_commit"`
	MonthsSinceLastCommit uint      `json:"months_since_last_commit"`
	NumContributors       uint      `json:"num_contributors"`
}

// GitStatsFetcher computes git stats after fetching using provided storage
type GitStatsFetcher struct {
	GitStorage GitProcessStorage
}

func (g *GitStatsFetcher) GetGitStats(gitURL url.URL) (GitStats, error) {
	stats := GitStats{
		URL: gitURL,
	}

	logs, err := NewGitLog(g.GitStorage.DirPath(gitURL))
	if err != nil {
		return stats, fmt.Errorf("can not get git logs: %w", err)
	}
	if len(logs) == 0 {
		return stats, errors.New("git log is empty")
	}

	stats.LastCommit = logs[0].AuthorDate
	stats.MonthsSinceLastCommit = logs.MonthsSinceLastCommit()
	stats.NumContributors = logs.NumContributors()

	return stats, nil
}
