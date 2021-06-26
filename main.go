package main

import (
	"context"
	_ "embed"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"

	"github.com/nikolaydubina/import-graph/pkg/awesomelists"
	"github.com/nikolaydubina/import-graph/pkg/codecov"
	"github.com/nikolaydubina/import-graph/pkg/collector"
	cgithub "github.com/nikolaydubina/import-graph/pkg/github"
	"github.com/nikolaydubina/import-graph/pkg/gitstats"
	"github.com/nikolaydubina/import-graph/pkg/gofilescanner"
	"github.com/nikolaydubina/import-graph/pkg/gomodgraph"
	"github.com/nikolaydubina/import-graph/pkg/goreportcard"
	"github.com/nikolaydubina/import-graph/pkg/gotestrunner"
	"github.com/nikolaydubina/import-graph/pkg/gourlresolver"
	"github.com/nikolaydubina/import-graph/pkg/gourlresolver/basiccache"
)

func main() {
	var runType string
	flag.StringVar(&runType, "i", "gomod", "type of input (e.g. gomod)")
	flag.Parse()

	ctx := context.Background()
	ghtoken := os.Getenv("GITHUB_IMPORT_GRAPH_TOKEN")
	if ghtoken == "" {
		log.Println("WARN: $GITHUB_IMPORT_GRAPH_TOKEN is empty, might not be able to fetch GitHub data")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghtoken},
	)
	tc := oauth2.NewClient(ctx, ts)

	gitClient := gitstats.GitCmdLocalClient{
		Path: ".import-graph/git-repos/",
	}

	switch runType {
	case "gomod":
		g, err := gomodgraph.GoModGraphParser{}.Parse(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		goModGraphCollector := collector.GoModuleGraphStatsCollector{
			ModuleCollector: collector.GoModuleStatsCollector{
				URLResolver: basiccache.GoCachedResolver{
					URLResolver: gourlresolver.GoURLResolver{HTTPClient: http.DefaultClient},
					Storage:     sync.Map{},
				},
				GitStorage: gitClient,
				GitStatsFetcher: gitstats.GitStatsFetcher{
					GitLogFetcher: &gitClient,
				},
				TestRunner: gotestrunner.GoCmdTestRunner{},
				CodecovClient: codecov.HTTPClient{
					HTTPClient: http.DefaultClient,
					BaseURL:    "api.codecov.io",
				},
				GoReportCardClient: goreportcard.GoReportCardHTTPClient{
					HTTPClient: http.DefaultClient,
					BaseURL:    "goreportcard.com",
				},
				FileScanner:         gofilescanner.FileScanner{},
				AwesomeListsChecker: awesomelists.AwesomeListsChecker{HTTPClient: http.DefaultClient},
				GitHubSummarizer: cgithub.GitHubSummarizer{
					GitHubClient: github.NewClient(tc),
				},
			},
		}
		goModGraphCollector.CollectStatsWrite(g, os.Stdout)
	default:
		log.Fatalln("unknown type of run")
	}
}
