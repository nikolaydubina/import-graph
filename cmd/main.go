package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

type GoModGraphBuilder interface {
	BuildGraph() (*iggo.Graph, error)
}

func main() {
	var (
		buildGraphFromCurrentDir bool
	)
	flag.BoolVar(&buildGraphFromCurrentDir, "graph-from-current-dir", false, "true means to build graph from current dir")

	var goModGraphBuilder GoModGraphBuilder
	if buildGraphFromCurrentDir {
		goModGraphBuilder = &iggo.GoCmdModGraphBuilder{
			GoModGraphParser: iggo.GoModGraphParser{},
		}
	} else {
		goModGraphBuilder = &iggo.GoReaderModGraphBuilder{
			Reader:           os.Stdin,
			GoModGraphParser: iggo.GoModGraphParser{},
		}
	}

	gitStorage := gitstats.GitProcessStorage{
		Path: ".import-graph/git-repos/",
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
