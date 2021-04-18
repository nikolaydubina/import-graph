package renderers

import (
	"io"
	"text/template"

	"github.com/nikolaydubina/import-graph/iggo"
)

type GraphvizRenderer struct{}

type Params struct {
	Graph iggo.Graph
}

func (g *GraphvizRenderer) Render(graph iggo.Graph, w io.Writer) error {
	basicDotTemplate := template.Must(template.New("basicDotTemplate").Parse(basicDotTemplateStr))
	return basicDotTemplate.Execute(w, Params{
		Graph: graph,
	})
}

const basicDotTemplateStr = `digraph G {
	concentrate=True;
	rankdir=TB;
	node [shape=record];
	
	{{range $i, $o := $.Nodes}}{{$o.ID}} [label="{{$o.ID}} :| {{$.o}}"];
	{{end}}

	{{range $i, $o := $.Edges}}{{$o.From}} -> {{$o.To}};
	{{end}}
}
`
