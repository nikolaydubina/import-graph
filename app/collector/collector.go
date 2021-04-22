package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

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
	GitStorage      *gitstats.GitCmdLocalClient
	URLResolver     *basiccache.GoCachedResolver
	GitStatsFetcher *gitstats.GitStatsFetcher
	TestRunner      *testrunner.GoCmdTestRunner
	CodecovClient   *codecov.HTTPClient
}

// CollectStats fetches all possible information about Go module
func (c *GoModuleStatsCollector) CollectStats(moduleName string) (ModuleStats, error) {
	moduleStats := ModuleStats{
		ID:         moduleName,
		ModuleName: moduleName,
	}
	var errFinal error

	gitURL, err := c.URLResolver.ResolveGitURL(moduleName)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not resolve URL: %w", err))
	}
	moduleStats.GitURL = gitURL.String()

	gitHubURL, err := c.URLResolver.ResolveGitHubURL(moduleName)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not resolve URL: %w", err))
	}
	moduleStats.GitHubURL = gitHubURL.String()

	if err := c.GitStorage.Clone(gitURL); err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not fetch git: %w", err))
	}

	if c.GitStatsFetcher != nil {
		if st, err := c.GitStatsFetcher.GetGitStats(gitURL); err != nil {
			errFinal = multierr.Combine(errFinal, fmt.Errorf("can not get git stats: %w", err))
		} else {
			moduleStats.GitStats = NewGitStats(st)
			moduleStats.CanGetGitStats = true
		}
	}

	if c.CodecovClient != nil {
		if resp, err := c.CodecovClient.GetRepoStatsFromGitHubURL(gitHubURL); err != nil {
			errFinal = multierr.Combine(errFinal, fmt.Errorf("can not get codecov stats: %w", err))
		} else {
			if st, err := NewCodecovStats(resp); err != nil {
				errFinal = multierr.Combine(errFinal, fmt.Errorf("can not format codecov stats: %w", err))
			} else {
				moduleStats.CanGetCodecovStats = true
				moduleStats.CodecovStats = st
			}
		}
	}

	if c.TestRunner != nil {
		if st, err := c.TestRunner.RunModuleTets(c.GitStorage.DirPath(gitURL)); err != nil {
			errFinal = multierr.Combine(errFinal, fmt.Errorf("can not run tests: %w", err))
		} else {
			moduleStats.CanRunTests = true
			moduleStats.GoModuleTestRunResult = st
		}
	}

	return moduleStats, errFinal
}

// GoModuleGraphStatsCollector collects data about Go modules and their relationships
type GoModuleGraphStatsCollector struct {
	ModuleCollector GoModuleStatsCollector
}

// CollectStats returns new Graph with collected data
// Keeps as much data as possible. Does no stop on errors, but keep track of them.
func (c *GoModuleGraphStatsCollector) CollectStats(gmod gomodgraph.Graph) (Graph, error) {
	var g Graph
	var finalErr error

	for i, n := range gmod.Modules {
		moduleWithStats, err := c.ModuleCollector.CollectStats(n.ModuleName)
		infoStr := ""
		if err != nil {
			finalErr = multierr.Combine(finalErr, fmt.Errorf("can not get module stats for module %s: %w", n.ModuleName, err))
			infoStr = fmt.Sprintf(" with error: %s", err)
		}
		g.Modules = append(g.Modules, moduleWithStats)
		log.Printf("[%d/%d] %s: done%s\n", i+1, len(gmod.Modules), n.ModuleName, infoStr)
	}

	for _, e := range gmod.Edges {
		g.Edges = append(g.Edges, Edge{From: e.From, To: e.To})
	}

	return g, finalErr
}

// CollectStatsWrite is version that serializes output as soon as it is computed
func (c *GoModuleGraphStatsCollector) CollectStatsWrite(gmod gomodgraph.Graph, w io.Writer) {
	encoder := json.NewEncoder(w)

	for _, n := range gmod.Modules {
		m, err := c.ModuleCollector.CollectStats(n.ModuleName)
		if err != nil {
			log.Println(fmt.Errorf("%s got error: %w", n.ModuleName, err))
		}
		if err := encoder.Encode(m); err != nil {
			log.Println(err)
		}
	}

	for _, e := range gmod.Edges {
		if err := encoder.Encode(Edge{From: e.From, To: e.To}); err != nil {
			log.Println(e)
		}
	}
}
