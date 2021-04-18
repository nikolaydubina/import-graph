package iggo

import (
	"fmt"

	"go.uber.org/multierr"

	"github.com/nikolaydubina/import-graph/gitstats"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

// GoModuleStatsCollector is a concrete class for collecting details about single Go module
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

// ModuleStats is stats about single module
type ModuleStats struct {
	ID         string `json:"id"` // used for JSONL graph rendering
	ModuleName string `json:"module_name"`
	gitstats.GitStats
	GoModuleTestRunResult
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
