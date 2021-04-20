package graphviz

import (
	"encoding/json"
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
// TODO: consider adding colors in background https://stackoverflow.com/questions/17765301/graphviz-dot-how-to-change-the-colour-of-one-record-in-multi-record-shape
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

// RenderValue coerces to json.Number and tries to avoid adding decimal points to integers
func RenderValue(v interface{}) string {
	if v, ok := v.(json.Number); ok {
		if vInt, err := v.Int64(); err == nil {
			return fmt.Sprintf("%d", vInt)
		}
		if v, err := v.Float64(); err == nil {
			return fmt.Sprintf("%.2f", v)
		}
	}
	return fmt.Sprintf("%v", v)
}

// RenderBasicLabel makes graphviz string for a single node
// This is pretty complex to write in Go template language due to map structure.
func RenderBasicLabel(n Node) string {
	rows := []string{}
	for k, v := range n {
		if k == "id" {
			continue
		}

		if strings.HasSuffix(k, "_url") {
			// URLs tend to be big and clutter dot outputs
			continue
		}

		rows = append(rows, fmt.Sprintf(`{%v\l | %s\r}`, k, RenderValue(v)))
	}

	// this will sort by key, since key is first
	sort.Strings(rows)

	return fmt.Sprintf("{ %s | %s }", n["id"], strings.Join(rows, " | "))
}
