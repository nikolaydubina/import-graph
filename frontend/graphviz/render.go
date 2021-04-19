package graphviz

import (
	"fmt"
	"io"
	"sort"
	"strings"
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
	basicTemplate, err := template.New("basicDotTemplate").Funcs(template.FuncMap{
		"nodeLabelBasic": RenderBasicLabel,
	}).Parse(basicTemplate)
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

// RenderBasicLabel makes graphviz string for a single node
// This is pretty complex to write in Go template language due to map structure.
func RenderBasicLabel(n Node) string {
	rows := []string{}
	for k, v := range n {
		if k == "id" {
			continue
		}

		valStr := fmt.Sprintf("%v", v)
		if v, ok := n[k].(int64); ok {
			valStr = fmt.Sprintf("%d", v)
		} else if v, ok := n[k].(float64); ok {
			valStr = fmt.Sprintf("%.2f", v)
		}

		rows = append(rows, fmt.Sprintf(`{%v\l | %s\r}`, k, valStr))
	}

	// this will sort by key, since key is first
	sort.Strings(rows)

	return fmt.Sprintf("{ %s | %s }", n["id"], strings.Join(rows, " | "))
}
