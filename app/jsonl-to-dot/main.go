// This script transforms JSONL graph from stdin to Graphviz in stdout
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nikolaydubina/import-graph/pkg/graphviz"
)

// loadColorConfigFromURL can load from local storage too like file:///myconfig.json
func loadColorConfigFromURL(path string) (*graphviz.ColorConfig, error) {
	if path == "" {
		return nil, errors.New("empty path")
	}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := &http.Client{Transport: t}

	res, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("can not load colorscheme file at path %s: %w", path, err)
	}

	colorschemeBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can not read file: %w", err)
	}

	var cconf graphviz.ColorConfig
	if err := json.Unmarshal(colorschemeBytes, &cconf); err != nil {
		return nil, fmt.Errorf("can not unmarshal: %w", err)
	}
	return &cconf, nil
}

func main() {
	var (
		colorSchemeFilePath string
	)
	flag.StringVar(&colorSchemeFilePath, "color-scheme", "", "optional path to colorscheme file (can be e.g. file://basic-colors.json)")
	flag.Parse()

	g, err := graphviz.NewGraphFromJSONLReader(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	if cconf, err := loadColorConfigFromURL(colorSchemeFilePath); cconf != nil && err == nil {
		if err := graphviz.NewGraphvizColorRenderer(*cconf).Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
			log.Println(err)
		}
		return
	} else {
		log.Println(err)
	}

	if err := graphviz.NewGraphvizBasicRenderer().Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
		log.Println(err)
	}
}
