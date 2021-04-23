// This script transforms JSONL graph from stdin to Graphviz in stdout
package main

import (
	"log"
	"os"

	"github.com/nikolaydubina/import-graph/pkg/graphviz"
)

func main() {
	g, err := graphviz.NewGraphFromJSONLReader(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	if err := graphviz.NewGraphvizRenderer().Render(graphviz.TemplateParams{Graph: g}, os.Stdout); err != nil {
		log.Println(err)
	}
}
