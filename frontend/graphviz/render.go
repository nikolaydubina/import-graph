package graphviz

import (
	"fmt"
	"io"
	"text/template"

	_ "embed"
)

//go:embed templates/basic.dot
var basicTemplate string

// GraphvizRenderer contains methods to tranform input to Graphviz format
type GraphvizRenderer struct {
	Template *template.Template
}

// NewGraphvizRenderer initializes template for reuse
func NewGraphvizRenderer() (*GraphvizRenderer, error) {
	basicTemplate, err := template.New("basicDotTemplate").Parse(basicTemplate)
	if err != nil {
		return nil, fmt.Errorf("can not init template: %w", err)
	}
	ret := GraphvizRenderer{
		Template: basicTemplate,
	}
	return &ret, nil
}

// TemplateParams is data passed to template
type TemplateParams struct {
	Graph Graph
}

// Render writes graph in Graphviz format to writer
func (g *GraphvizRenderer) Render(params TemplateParams, w io.Writer) error {
	return g.Template.Execute(w, params)
}
