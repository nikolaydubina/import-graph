package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/nikolaydubina/import-graph/frontend/graphviz"
	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo"
	"github.com/nikolaydubina/import-graph/iggo/gomodgraph"
	"github.com/nikolaydubina/import-graph/iggo/testrunner"
	"github.com/nikolaydubina/import-graph/iggo/urlresolver"
	"github.com/nikolaydubina/import-graph/iggo/urlresolver/basiccache"
)

type OutputType struct{ V string }

var (
	OutputTypeJSONL = OutputType{V: "jsonl"}
	OutputTypeDot   = OutputType{V: "dot"}
)

func main() {
	var (
		outputType OutputType
	)
	flag.StringVar(&outputType.V, "output", "jsonl", "output type (jsonl, dot)")
	flag.Parse()

	gitStorage := gitstats.GitProcessStorage{
		Path: ".import-graph/git-repos/",
	}
	goModGraphParser := &gomodgraph.GoModGraphParser{}
	goModGraphCollector := iggo.GoModuleGraphStatsCollector{
		ModuleCollector: iggo.GoModuleStatsCollector{
			URLResolver: &basiccache.GoCachedResolver{
				URLResolver: &urlresolver.GoURLResolver{HTTPClient: http.DefaultClient},
				Storage:     sync.Map{},
			},
			GitStorage: &gitStorage,
			GitStatsFetcher: &gitstats.GitStatsFetcher{
				GitStorage: &gitStorage,
			},
			TestRunner: testrunner.GoCmdTestRunner{},
		},
	}
	graphVizRenderer, err := graphviz.NewGraphvizRenderer()
	if err != nil {
		log.Fatalln(err)
	}

	g, err := goModGraphParser.Parse(os.Stdin)
	if err != nil {
		log.Println(err)
	}

	gCollected, err := goModGraphCollector.CollectStats(*g)
	if err != nil {
		log.Println(err)
	}

	switch outputType {
	case OutputTypeJSONL:
		if err := gCollected.WriteJSONL(os.Stdout); err != nil {
			log.Println(err)
		}
	case OutputTypeDot:
		var buf bytes.Buffer

		if err := gCollected.WriteJSONL(&buf); err != nil {
			log.Fatalln(err)
		}

		g, err := graphviz.NewGraphFromJSONLReader(&buf)
		if err != nil {
			log.Fatalln(err)
		}

		if err := graphVizRenderer.Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
			log.Println(err)
		}
	default:
		log.Fatalf("unexpected output type %s", outputType)
	}
}
