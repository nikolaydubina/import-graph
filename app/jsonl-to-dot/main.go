// This script transforms JSONL graph from stdin to Graphviz in stdout
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	_ "embed"

	"github.com/nikolaydubina/import-graph/pkg/graphviz"
)

//go:embed basic-colors.json
var defaultColorsConfig []byte

func main() {
	var (
		colorFlag bool
	)
	flag.BoolVar(&colorFlag, "color", false, "set to make colored")
	flag.Parse()

	g, err := graphviz.NewGraphFromJSONLReader(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	if colorFlag {
		var cconf graphviz.ColorConfig
		if err := json.Unmarshal(defaultColorsConfig, &cconf); err != nil {
			log.Println(err)
		}

		if err := graphviz.NewGraphvizColorRenderer(cconf).Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
			log.Println(err)
		}
	} else {
		if err := graphviz.NewGraphvizBasicRenderer().Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
			log.Println(err)
		}
	}
}
