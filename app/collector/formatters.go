package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/nikolaydubina/import-graph/pkg/codecov"
	"github.com/nikolaydubina/import-graph/pkg/gitstats"
	"github.com/nikolaydubina/import-graph/pkg/go/goreportcard"
	"github.com/nikolaydubina/import-graph/pkg/go/testrunner"
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
	LastCommit          string `json:"git_last_commit,omitempty"`  // applying formatting to days
	DaysSinceLastCommit uint   `json:"git_last_commit_days_since"` // num full days
	NumContributors     uint   `json:"git_num_contributors"`
}

func NewGitStats(r *gitstats.GitStats) *GitStats {
	if r == nil {
		return nil
	}
	return &GitStats{
		LastCommit:          r.LastCommit.Format("2006-01-02"),
		DaysSinceLastCommit: uint(math.Floor(r.DaysSinceLastCommit)),
		NumContributors:     r.NumContributors,
	}
}

type GoReportCardStats struct {
	Average   json.Number            `json:"goreportcard_average"`
	Grade     goreportcard.GradeEnum `json:"goreportcard_grade"`
	NumFiles  uint                   `json:"goreportcard_files"`
	NumIssues uint                   `json:"goreportcard_issues"`
}

func NewGoReportCardStats(r *goreportcard.Report) *GoReportCardStats {
	if r == nil {
		return nil
	}
	return &GoReportCardStats{
		Average:   json.Number(fmt.Sprintf("%.2f", r.Average)),
		Grade:     r.Grade,
		NumFiles:  r.NumFiles,
		NumIssues: r.NumIssues,
	}
}

type FileStats struct {
	HasBenchmarks bool `json:"files_has_benchmarks"`
	HasTests      bool `json:"files_has_tests"`
}

type GoTestStats struct {
	HasTests               bool    `json:"gotest_has_tests"`
	AllTestsPassed         bool    `json:"gotest_all_tests_passed"`
	NumPackages            uint    `json:"gotest_num_packages"`
	NumPackagesWithTests   uint    `json:"gotest_num_packages_with_tests"`
	NumPackagesTestsPassed uint    `json:"gotest_num_packages_tests_passed"`
	AvgPackageCoverage     float64 `json:"gotest_package_coverage_avg"`
}

func NewGoTestStats(r *testrunner.GoModuleTestRunResult) *GoTestStats {
	return &GoTestStats{
		HasTests:               r.HasTests,
		AllTestsPassed:         r.AllTestsPassed,
		NumPackages:            r.NumPackages,
		NumPackagesWithTests:   r.NumPackagesWithTests,
		NumPackagesTestsPassed: r.NumPackagesTestsPassed,
		AvgPackageCoverage:     r.AvgPackageCoverage,
	}
}
