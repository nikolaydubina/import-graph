package iggo

import (
	"encoding/json"
	"fmt"
	"io"

	"go.uber.org/multierr"

	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo/gomodgraph"
	"github.com/nikolaydubina/import-graph/iggo/testrunner"
	"github.com/nikolaydubina/import-graph/iggo/urlresolver/basiccache"
)

// ModuleStats is stats about single module
type ModuleStats struct {
	ID         string `json:"id"` // unique key among all nodes, for Go this is module name
	ModuleName string `json:"module_name"`

	// Data bellow will be filled by appropriate routines
	gitstats.GitStats
	testrunner.GoModuleTestRunResult
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
	}

	if err := c.GitStorage.Fetch(*gitURL); err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not fetch git: %w", err))
	}
	gitDirPath := c.GitStorage.DirPath(*gitURL)

	gitStats, err := c.GitStatsFetcher.GetGitStats(gitDirPath)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not get git stats: %w", err))
	}
	moduleStats.GitStats = gitStats

	testStats, err := c.TestRunner.RunModuleTets(gitDirPath)
	if err != nil {
		errFinal = multierr.Combine(errFinal, fmt.Errorf("can not run tests: %w", err))
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
