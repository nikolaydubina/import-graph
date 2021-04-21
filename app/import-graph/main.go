package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/nikolaydubina/import-graph/app/collector"

	"github.com/nikolaydubina/import-graph/pkg/codecov"
	"github.com/nikolaydubina/import-graph/pkg/gitstats"
	"github.com/nikolaydubina/import-graph/pkg/go/gomodgraph"
	"github.com/nikolaydubina/import-graph/pkg/go/testrunner"
	"github.com/nikolaydubina/import-graph/pkg/go/urlresolver"
	"github.com/nikolaydubina/import-graph/pkg/go/urlresolver/basiccache"
	"github.com/nikolaydubina/import-graph/pkg/graphviz"
)

type OutputType struct{ V string }

var (
	OutputTypeJSONL = OutputType{V: "jsonl"}
	OutputTypeDot   = OutputType{V: "dot"}
)

func main() {
	var (
		outputType OutputType
		runTests   bool
	)
	flag.StringVar(&outputType.V, "output", "jsonl", "output type (jsonl, dot)")
	flag.BoolVar(&runTests, "test", false, "set to run tests")
	flag.Parse()

	var testRunner *testrunner.GoCmdTestRunner
	if runTests {
		testRunner = &testrunner.GoCmdTestRunner{}
	}

	gitClient := gitstats.GitCmdLocalClient{
		Path: ".import-graph/git-repos/",
	}
	goModGraphParser := &gomodgraph.GoModGraphParser{}
	goModGraphCollector := collector.GoModuleGraphStatsCollector{
		ModuleCollector: collector.GoModuleStatsCollector{
			URLResolver: &basiccache.GoCachedResolver{
				URLResolver: &urlresolver.GoURLResolver{HTTPClient: http.DefaultClient},
				Storage:     sync.Map{},
			},
			GitStorage: &gitClient,
			GitStatsFetcher: &gitstats.GitStatsFetcher{
				GitLogFetcher: &gitClient,
			},
			TestRunner: testRunner,
			CodecovClient: &codecov.HTTPClient{
				HTTPClient: http.DefaultClient,
				BaseURL:    "api.codecov.io",
			},
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
