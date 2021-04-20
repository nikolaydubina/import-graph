package collector

import (
	"encoding/json"
	"fmt"
	"io"

	"go.uber.org/multierr"

	"github.com/nikolaydubina/import-graph/pkg/codecov"
	"github.com/nikolaydubina/import-graph/pkg/gitstats"
	"github.com/nikolaydubina/import-graph/pkg/go/gomodgraph"
	"github.com/nikolaydubina/import-graph/pkg/go/testrunner"
	"github.com/nikolaydubina/import-graph/pkg/go/urlresolver/basiccache"
)

// ModuleStats is stats about single module
type ModuleStats struct {
	ID         string `json:"id"` // unique key among all nodes, for Go this is module name
	ModuleName string `json:"-"`  // this is in id anyways

	// Collected data
	CanGetGitStats     bool `json:"can_get_gitstats"`
	CanGetCodecovStats bool `json:"can_get_codecov"`
	CanRunTests        bool `json:"can_run_tests"`

	GitHubURL string `json:"github_url,omitempty"`
	GitURL    string `json:"git_url,omitempty"`

	*GitStats                         `json:",omitempty"`
	*CodecovStats                     `json:",omitempty"`
	*testrunner.GoModuleTestRunResult `json:",omitempty"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Graph struct {
	Modules []ModuleStats
	Edges   []Edge
}

// WriteJSONL which is default export format
func (g *Graph) WriteJSONL(w io.Writer) error {
	encoder := json.NewEncoder(w)
	var finalErr error

	for _, n := range g.Modules {
		finalErr = multierr.Combine(encoder.Encode(n))
	}

	for _, e := range g.Edges {
		finalErr = multierr.Combine(encoder.Encode(e))
	}

	return finalErr
}

// GoModuleStatsCollector is collecting all the details about single Go module
// Does not fail if encounters errors, but still collects thoese errors.
type GoModuleStatsCollector struct {
	GitStorage      *gitstats.GitProcessStorage
	URLResolver     *basiccache.GoCachedResolver
	GitStatsFetcher *gitstats.GitStatsFetcher
	TestRunner      testrunner.GoCmdTestRunner
	CodecovClient   *codecov.HTTPClient
}

// CollectStats fetches all possible information about Go module
func (c *GoModuleStatsCollector) CollectStats(moduleName string) (moduleStats ModuleStats, errFinal error) {
	moduleStats = ModuleStats{
		ID:         moduleName,
		ModuleName: moduleName,
	}

	gitURL, err := c.URLResolver.ResolveGitURL(moduleName)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not resolve URL: %w", err))
		return
	}
	moduleStats.GitURL = gitURL.String()
	if err := c.GitStorage.Fetch(gitURL); err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not fetch git: %w", err))
		return
	}
	gitDirPath := c.GitStorage.DirPath(gitURL)

	// Git Stats
	gitStats, err := c.GitStatsFetcher.GetGitStats(gitDirPath)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not get git stats: %w", err))
	}
	if gitStats != nil {
		moduleStats.GitStats = NewGitStats(gitStats)
		moduleStats.CanGetGitStats = true
	}

	// GitHub URL
	gitHubURL, err := c.URLResolver.ResolveGitHubURL(moduleName)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not resolve URL: %w", err))
	}
	moduleStats.GitHubURL = gitHubURL.String()

	// codecov.io
	codecovStats, err := c.CodecovClient.GetRepoStatsFromGitHubURL(gitHubURL)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not get codecov stats: %w", err))
	}
	if cc, err := NewCodecovStats(codecovStats); cc != nil && err == nil {
		moduleStats.CodecovStats = cc
		moduleStats.CanGetCodecovStats = true
	} else if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not format codecov stats: %w", err))
	}

	// Run tests
	testStats, err := c.TestRunner.RunModuleTets(gitDirPath)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not run tests: %w", err))
	}
	if testStats != nil {
		moduleStats.CanRunTests = true
	}
	moduleStats.GoModuleTestRunResult = testStats

	return
}

// GoModuleGraphStatsCollector collects data about Go modules and their relationships
type GoModuleGraphStatsCollector struct {
	ModuleCollector GoModuleStatsCollector
}

// CollectStats returns new Graph with collected data
// Keeps as much data as possible. Does no stop on errors, but keep track of them.
func (c *GoModuleGraphStatsCollector) CollectStats(gmod gomodgraph.Graph) (g Graph, finalErr error) {
	for _, n := range gmod.Modules {
		moduleWithStats, err := c.ModuleCollector.CollectStats(n.ModuleName)
		if err != nil {
			finalErr = multierr.Combine(finalErr, fmt.Errorf("can not get module stats for module %s: %w", n.ModuleName, err))
		}
		g.Modules = append(g.Modules, moduleWithStats)
	}

	for _, e := range gmod.Edges {
		g.Edges = append(g.Edges, Edge{From: e.From, To: e.To})
	}

	return
}
