package graphviz

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"sort"
	"strings"
	"text/template"

	_ "embed"
)

//go:embed templates/color.dot
var colorTemplate string

type ColorConfigVal struct {
	ValToColor map[string]color.RGBA
}

type ColorConfig map[string]ColorConfigVal

func (c ColorConfig) RenderColorVal(k string, v interface{}) color.Color {
	valC, ok := c[k]
	if !ok {
		return color.White
	}

	var key string
	if vs, ok := v.(string); ok {
		key = vs
	} else {
		vs, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		key = string(vs)
	}

	if c, ok := valC.ValToColor[key]; ok {
		return c
	}

	return color.White
}

// GraphvizColorRenderer contains methods to tranform input to Graphviz format
// TODO: consider adding colors in background https://stackoverflow.com/questions/17765301/graphviz-dot-how-to-change-the-colour-of-one-record-in-multi-record-shape
type GraphvizColorRenderer struct {
	Template    *template.Template
	ColorConfig ColorConfig
}

// NewGraphvizColorRenderer initializes template for reuse
func NewGraphvizColorRenderer(conf ColorConfig) *GraphvizColorRenderer {

	ret := GraphvizColorRenderer{
		ColorConfig: conf,
	}
	ret.Template = template.Must(template.New("colorDotTemplate").Funcs(template.FuncMap{
		"nodeLabelTableColored": ret.RenderLabelTableColored,
	}).Parse(colorTemplate))

	return &ret
}

// Render writes graph in Graphviz format to writer
func (g *GraphvizColorRenderer) Render(params TemplateParams, w io.Writer) error {
	return g.Template.Execute(w, params)
}

func RenderColor(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("#%x%x%x%x", uint8(r), uint8(g), uint8(b), uint8(a))
}

// RenderLabelTableColored makes graphviz string for a single node with colored table
func (c *GraphvizColorRenderer) RenderLabelTableColored(n Node) string {
	rows := []string{}
	for k, v := range n {
		if k == "id" || strings.HasSuffix(k, "_url") {
			continue
		}

		row := fmt.Sprintf(`
			<tr>
				<td border="1" ALIGN="LEFT">%s</td>
				<td border="1" ALIGN="RIGHT" bgcolor="%s">%s</td>
			</tr>`,
			k,
			RenderColor(c.ColorConfig.RenderColorVal(k, v)),
			RenderValue(v),
		)

		rows = append(rows, row)
	}

	// this will sort by key, since key is first
	sort.Strings(rows)

	return strings.Join(
		[]string{
			"<<table border=\"0\" cellspacing=\"0\" CELLPADDING=\"6\">",
			fmt.Sprintf(`
				<tr>
					<td port="port0" border="1" colspan="2" ALIGN="CENTER" bgcolor="%s">%s</td>
				</tr>`,
				RenderColor(color.RGBA{R: 200, G: 200, B: 200, A: 200}),
				RenderValue(n["id"]),
			),
			strings.Join(rows, "\n"),
			"</table>>",
		},
		"\n",
	)
}
