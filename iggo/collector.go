package iggo

import (
	"log"
	"net/http"
	"sync"

	"github.com/nikolaydubina/import-graph/gitstats"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

type GoModuleStatsCollector struct {
	GraphBuilder    GoProcessGraphBuilder
	URLResolver     iggorc.GoCachedResolver
	GitStatsFetcher gitstats.GitStatsFetcher
	TestRunner      GoProcessTestRunner
	GitStorage      *gitstats.GitProcessStorage
}

func NewGoModStatsCollector() GoModuleStatsCollector {
	gitStorage := gitstats.GitProcessStorage{
		Path: ".import-graph/git-repos/",
	}
	return GoModuleStatsCollector{
		GraphBuilder: GoProcessGraphBuilder{},
		URLResolver:  iggorc.GoCachedResolver{Resolver: &GoResolver{HTTPClient: http.DefaultClient}, Storage: sync.Map{}},
		GitStorage:   &gitStorage,
		GitStatsFetcher: gitstats.GitStatsFetcher{
			GitStorage: &gitStorage,
		},
		TestRunner: GoProcessTestRunner{},
	}
}

type ModuleStats struct {
	gitstats.GitStats
	GoModuleTestRunResult
}

func (c *GoModuleStatsCollector) CollectStats() (ModuleStats, error) {
	g, err := c.GraphBuilder.BuildGraph()
	if err != nil {
		log.Println(err)
	}
	for _, n := range g.Nodes {
		log.Printf("%s: start\n", n.Name)
		moduleStats := ModuleStats{}

		gitURL, err := c.URLResolver.ResolveGitURL(n.Name)
		if err != nil {
			log.Printf("%s: %s", n.Name, err)
			continue
		}
		log.Printf("%s: url resolved\n", n.Name)
		if err := c.GitStorage.Fetch(*gitURL); err != nil {
			log.Printf("%s: %s", n.Name, err)
			continue
		}
		log.Printf("%s: git is cloned\n", n.Name)

		gitStats, err := c.GitStatsFetcher.GetGitStats(*gitURL)
		if err != nil {
			log.Printf("%s: %s", n.Name, err)
		}
		log.Printf("%s: git stats computed\n", n.Name)
		moduleStats.GitStats = gitStats

		testStats, err := c.TestRunner.RunModuleTets(c.GitStorage.DirPath(*gitURL))
		if err != nil {
			log.Printf("%s: %s", n.Name, err)
		}
		moduleStats.GoModuleTestRunResult = testStats

		log.Printf("%s: done with result: %+v\n", n.Name, moduleStats)
	}
}
