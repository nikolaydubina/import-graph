package gitstats

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type GitStorage interface {
	DirPath(gitURL url.URL) string
}

// GitStats contains information about single git repository computed using git only
type GitStats struct {
	LastCommit            time.Time `json:"last_commit,omitempty"`
	DaysSinceLastCommit   float64   `json:"last_commit_days_since"`
	YearsSinceLastCommit  float64   `json:"last_commit_years_since"`
	MonthsSinceLastCommit float64   `json:"last_commit_months_since"`
	NumContributors       uint      `json:"num_contributors"`
}

// GitStatsFetcher computes git stats after fetching using provided storage
type GitStatsFetcher struct {
	GitStorage GitStorage
}

func (g *GitStatsFetcher) GetGitStats(gitDirPath string) (*GitStats, error) {
	logs, err := NewGitLog(gitDirPath)
	if err != nil {
		return nil, fmt.Errorf("can not get git logs: %w", err)
	}
	if len(logs) == 0 {
		return nil, errors.New("git log is empty")
	}

	stats := GitStats{
		LastCommit:            logs[0].AuthorDate,
		DaysSinceLastCommit:   logs.DaysSinceLastCommit(),
		MonthsSinceLastCommit: logs.DaysSinceLastCommit() / 28.0,
		YearsSinceLastCommit:  logs.DaysSinceLastCommit() / 28.0 / 12.0,
		NumContributors:       logs.NumContributors(),
	}
	return &stats, nil
}
