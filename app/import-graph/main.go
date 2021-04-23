package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/nikolaydubina/import-graph/app/collector"

	"github.com/nikolaydubina/import-graph/pkg/codecov"
	"github.com/nikolaydubina/import-graph/pkg/gitstats"
	"github.com/nikolaydubina/import-graph/pkg/go/filescanner"
	"github.com/nikolaydubina/import-graph/pkg/go/gomodgraph"
	"github.com/nikolaydubina/import-graph/pkg/go/goreportcard"
	"github.com/nikolaydubina/import-graph/pkg/go/testrunner"
	"github.com/nikolaydubina/import-graph/pkg/go/urlresolver"
	"github.com/nikolaydubina/import-graph/pkg/go/urlresolver/basiccache"
	"github.com/nikolaydubina/import-graph/pkg/graphviz"

	_ "embed"
)

//go:embed basic-colors.json
var defaultColorsConfig []byte

type OutputType struct{ V string }

var (
	OutputTypeJSONL    = OutputType{V: "jsonl"}
	OutputTypeDot      = OutputType{V: "dot"}
	OutputTypeDotColor = OutputType{V: "dot-color"}
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
			GoReportCardClient: &goreportcard.GoReportCardHTTPClient{
				HTTPClient: http.DefaultClient,
				BaseURL:    "goreportcard.com",
			},
			FileScanner: &filescanner.FileScanner{},
		},
	}

	g, err := goModGraphParser.Parse(os.Stdin)
	if err != nil {
		log.Println(err)
	}
	if g == nil {
		log.Fatalln("parsed graph is nil")
	}

	switch outputType {
	case OutputTypeJSONL:
		goModGraphCollector.CollectStatsWrite(*g, os.Stdout)
	case OutputTypeDot:
		var buf bytes.Buffer

		goModGraphCollector.CollectStatsWrite(*g, &buf)
		g, err := graphviz.NewGraphFromJSONLReader(&buf)
		if err != nil {
			log.Fatalln(err)
		}

		if err := graphviz.NewGraphvizBasicRenderer().Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
			log.Println(err)
		}
	case OutputTypeDotColor:
		var buf bytes.Buffer

		goModGraphCollector.CollectStatsWrite(*g, &buf)
		g, err := graphviz.NewGraphFromJSONLReader(&buf)
		if err != nil {
			log.Fatalln(err)
		}

		var cconf graphviz.ColorConfig
		if err := json.Unmarshal(defaultColorsConfig, &cconf); err != nil {
			log.Println(err)
		}

		if err := graphviz.NewGraphvizColorRenderer(cconf).Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
			log.Println(err)
		}
	default:
		log.Fatalf("unexpected output type %s", outputType)
	}
}
