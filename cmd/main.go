package main

import (
	"flag"
	"io"
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

type Renderer interface {
	Render(w io.Writer) error
}

func main() {
	var (
		buildGraphFromCurrentDir bool
		outputType               string
	)

	flag.BoolVar(&buildGraphFromCurrentDir, "graph-from-current-dir", false, "true means to build graph from current dir")
	flag.StringVar(&outputType, "output", "jsonl", "output type (jsonl, dot)")

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
	goModGraphCollector := iggo.GoModuleGraphStatsCollector{
		ModuleCollector: iggo.GoModuleStatsCollector{
			URLResolver: &iggorc.GoCachedResolver{Resolver: &iggo.GoResolver{HTTPClient: http.DefaultClient}, Storage: sync.Map{}},
			GitStorage:  &gitStorage,
			GitStatsFetcher: &gitstats.GitStatsFetcher{
				GitStorage: &gitStorage,
			},
			TestRunner: iggo.GoCmdTestRunner{},
		},
	}

	g, err := goModGraphBuilder.BuildGraph()
	if err != nil {
		log.Println(err)
	}

	gCollected, err := goModGraphCollector.CollectStats(*g)
	if err != nil {
		log.Println(err)
	}

	if outputType == "dot" {
		panic("TODO: dot")
	}

	if err := gCollected.WriteJSONL(os.Stdout); err != nil {
		log.Println(err)
	}
}
