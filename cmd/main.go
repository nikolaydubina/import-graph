package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

// check codecov
// check readme mentions alpha/beta
// version of package is stable same as godoc
// readme has go-report card
// readme reports code coverage
// try run Makefile lint
// try run linting
// benchmarks detected

// GitHub -> stars, organization, contributor profiles
// if github get github page (Python/JS headless browser?)
// if github check owner
// if github owner is organization: match against lists
// if github owner is user: collect stats on other repos; (try fetch linkedin?)

func main() {
	goCollector := GoModuleStatsCollector{}

	modGraphBuilder := iggo.GoProcessGraphBuilder{}
	modURLResolver := iggorc.GoCachedResolver{Resolver: &iggo.GoResolver{HTTPClient: http.DefaultClient}, Storage: sync.Map{}}
	gitStorage := gitstats.GitProcessStorage{
		Path: ".import-graph/git-repos/",
	}
	gitStatsFetcher := gitstats.GitStatsFetcher{
		GitStorage: &gitStorage,
	}
	testRunner := iggo.GoProcessTestRunner{}

	g, err := modGraphBuilder.BuildGraph()
	if err != nil {
		log.Println(err)
	}
	for _, n := range g.Nodes {
		log.Printf("%s: start\n", n.Name)
		result, err := goCollector.Collect()

		moduleStats := ModuleStats{}

		gitURL, err := modURLResolver.ResolveGitURL(n.Name)
		if err != nil {
			log.Printf("%s: %s", n.Name, err)
			continue
		}
		log.Printf("%s: url resolved\n", n.Name)
		if err := gitStatsFetcher.GitStorage.Fetch(*gitURL); err != nil {
			log.Printf("%s: %s", n.Name, err)
			continue
		}
		log.Printf("%s: git is cloned\n", n.Name)

		gitStats, err := gitStatsFetcher.GetGitStats(*gitURL)
		if err != nil {
			log.Printf("%s: %s", n.Name, err)
		}
		log.Printf("%s: git stats computed\n", n.Name)
		moduleStats.GitStats = gitStats

		testStats, err := testRunner.RunModuleTets(gitStorage.DirPath(*gitURL))
		if err != nil {
			log.Printf("%s: %s", n.Name, err)
		}
		moduleStats.GoModuleTestRunResult = testStats

		log.Printf("%s: done with result: %+v\n", n.Name, moduleStats)
	}
}
