package iggo

import (
	"fmt"

	"go.uber.org/multierr"

	"github.com/nikolaydubina/import-graph/gitstats"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

// GoModuleStatsCollector is collecting all the details about single Go module
// Does not fail if encounters errors, but still collects thoese errors.
//
// TODO:
// check codecov
// check readme mentions alpha/beta
// version of package is stable same as godoc
// readme has go-report card
// readme reports code coverage
// try run Makefile lint
// try run linting
// benchmarks detected
//
// GitHub -> stars, organization, contributor profiles
// if github get github page (Python/JS headless browser?)
// if github check owner
// if github owner is organization: match against lists
// if github owner is user: collect stats on other repos; (try fetch linkedin?)
type GoModuleStatsCollector struct {
	GitStorage      *gitstats.GitProcessStorage
	URLResolver     *iggorc.GoCachedResolver
	GitStatsFetcher *gitstats.GitStatsFetcher
	TestRunner      GoCmdTestRunner
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
func (c *GoModuleGraphStatsCollector) CollectStats(g Graph) (finalGraph Graph, finalErr error) {
	for _, n := range g.Modules {
		moduleWithStats, err := c.ModuleCollector.CollectStats(n.ModuleName)
		if err != nil {
			finalErr = multierr.Combine(finalErr, fmt.Errorf("can not get module stats for module %s: %w", n.ModuleName, err))
		}
		finalGraph.Modules = append(finalGraph.Modules, moduleWithStats)
	}

	finalGraph.Edges = append(finalGraph.Edges, g.Edges...)

	return
}
