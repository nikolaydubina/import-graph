package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo"
	iggorc "github.com/nikolaydubina/import-graph/iggo/resolver_cached"
)

func main() {
	var (
		outputType string
	)

	flag.StringVar(&outputType, "output", "jsonl", "output type (jsonl, dot)")

	gitStorage := gitstats.GitProcessStorage{
		Path: ".import-graph/git-repos/",
	}
	goModGraphParser := &iggo.GoModGraphParser{}
	goModGraphCollector := iggo.GoModuleGraphStatsCollector{
		ModuleCollector: iggo.GoModuleStatsCollector{
			URLResolver: &iggorc.GoCachedResolver{
				URLResolver: &iggo.GoURLResolver{HTTPClient: http.DefaultClient},
				Storage:     sync.Map{},
			},
			GitStorage: &gitStorage,
			GitStatsFetcher: &gitstats.GitStatsFetcher{
				GitStorage: &gitStorage,
			},
			TestRunner: iggo.GoCmdTestRunner{},
		},
	}

	g, err := goModGraphParser.Parse(os.Stdin)
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
