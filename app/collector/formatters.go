package collector

import (
	"errors"

	"github.com/nikolaydubina/import-graph/pkg/codecov"
	"github.com/nikolaydubina/import-graph/pkg/gitstats"
)

type CodecovStats struct {
	RepoURL  string  `json:"codecov_url"` // need flat value so can not use url.URL
	NumFiles uint    `json:"codecov_files"`
	NumLines uint    `json:"codecov_lines"`
	Coverage float64 `json:"codecov_coverage"`
}

func NewCodecovStats(r *codecov.RepoStats) (*CodecovStats, error) {
	if r == nil {
		return nil, errors.New("codecov object is nil")
	}
	if r.LatestCommmit == nil {
		return nil, errors.New("latest commit is not found in codecov")
	}
	stats := CodecovStats{
		RepoURL:  r.RepoURL.String(),
		NumFiles: r.LatestCommmit.Report.Totals.NumFiles,
		NumLines: r.LatestCommmit.Report.Totals.NumLines,
		Coverage: r.LatestCommmit.Report.Totals.Coverage,
	}
	return &stats, nil
}

type GitStats struct {
	LastCommit            string  `json:"git_last_commit,omitempty"` // applying formatting to days
	DaysSinceLastCommit   float64 `json:"git_last_commit_days_since"`
	YearsSinceLastCommit  float64 `json:"git_last_commit_years_since"`
	MonthsSinceLastCommit float64 `json:"git_last_commit_months_since"`
	NumContributors       uint    `json:"git_num_contributors"`
}

func NewGitStats(r *gitstats.GitStats) *GitStats {
	if r == nil {
		return nil
	}
	return &GitStats{
		LastCommit:            r.LastCommit.Format("2006-01-02"),
		DaysSinceLastCommit:   r.DaysSinceLastCommit,
		MonthsSinceLastCommit: r.MonthsSinceLastCommit,
		YearsSinceLastCommit:  r.YearsSinceLastCommit,
		NumContributors:       r.NumContributors,
	}
}
