package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/nikolaydubina/import-graph/iggo"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

// get repo code
// get stats about git repo (last updated; num people updated)
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
	g, _ := gc.BuildGraph()
	for _, n := range g.Nodes {
		gitHubURL, err := gf.ResolveGitHubURL(n.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(gitHubURL)
	}
	return
}
