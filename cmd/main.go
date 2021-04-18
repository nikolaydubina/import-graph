package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

func main() {
	gitStorage := gitstats.GitProcessStorage{
		Path: ".import-graph/git-repos/",
	}
	goModGraphBuilder := iggo.GoCmdModGraphBuilder{
		GoModGraphParser: iggo.GoModGraphParser{},
	}
	goModCollector := iggo.GoModuleStatsCollector{
		URLResolver: &iggorc.GoCachedResolver{Resolver: &iggo.GoResolver{HTTPClient: http.DefaultClient}, Storage: sync.Map{}},
		GitStorage:  &gitStorage,
		GitStatsFetcher: &gitstats.GitStatsFetcher{
			GitStorage: &gitStorage,
		},
		TestRunner: iggo.GoCmdTestRunner{},
	}

	g, err := goModGraphBuilder.BuildGraph()
	if err != nil {
		log.Println(err)
	}

	encoder := json.NewEncoder(os.Stdout)

	// write nodes
	for _, n := range g.Nodes {
		result, err := goModCollector.CollectStats(n.Name)
		if err != nil {
			log.Println(err)
			continue
		}

		if err := encoder.Encode(result); err != nil {
			log.Println(err)
		}
	}

	// write edges
	for _, e := range g.Edges {
		if err := encoder.Encode(e); err != nil {
			log.Println(err)
		}
	}
}
