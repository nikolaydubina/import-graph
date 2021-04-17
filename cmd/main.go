package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

// check test data?
// check code coverage?
// check readme mentions alpha/beta
// version of package is stable same as godoc
// readme has go-report card
// readme reports code coverage
// try run Makefile lint
// try run linting
// try run go test
// test files detected
// go test ./...
// benchmarks detected

// GitHub -> stars, organization, contributor profiles
// if github get github page (Python/JS headless browser?)
// if github check owner
// if github owner is organization: match against lists
// if github owner is user: collect stats on other repos; (try fetch linkedin?)

func main() {
	gc := iggo.GoProcessGraphBuilder{}
	gf := iggorc.GoCachedResolver{Resolver: &iggo.GoResolver{HTTPClient: http.DefaultClient}, Storage: sync.Map{}}
	gitStatsFetcher := gitstats.GitStatsFetcher{
		GitStorage: &gitstats.GitProcessStorage{
			Path: ".import-graph/git-repos/",
		},
	}

	g, err := gc.BuildGraph()
	if err != nil {
		log.Println(err)
	}
	for _, n := range g.Nodes {
		gitURL, err := gf.ResolveGitURL(n.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := gitStatsFetcher.GitStorage.Fetch(*gitURL); err != nil {
			log.Println(err)
			continue
		}

		gitStats, err := gitStatsFetcher.GetGitStats(*gitURL)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%+v\n", gitStats)
	}
}
