package iggo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"go.uber.org/multierr"

	"github.com/nikolaydubina/import-graph/gitstats"
	"github.com/nikolaydubina/import-graph/iggo/testrunner"
)

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ModuleStats is stats about single module
type ModuleStats struct {
	ID         string `json:"id"` // unique key among all nodes, for Go this is module name
	ModuleName string `json:"module_name"`

	// Data bellow will be filled by appropriate routines
	gitstats.GitStats
	testrunner.GoModuleTestRunResult
}

type Graph struct {
	Modules []ModuleStats `json:"nodes"`
	Edges   []Edge        `json:"edges"`
}

// WriteJSONL which is default export format
func (g *Graph) WriteJSONL(w io.Writer) error {
	encoder := json.NewEncoder(w)
	var finalErr error

	for _, n := range g.Modules {
		finalErr = multierr.Combine(encoder.Encode(n))
	}

	for _, e := range g.Edges {
		finalErr = multierr.Combine(encoder.Encode(e))
	}

	return finalErr
}

// GoModGraphParser builds graph from output of `go mod graph`
// This is conveneint if caller can call `go mod graph` by himself.
type GoModGraphParser struct{}

func (c *GoModGraphParser) Parse(input io.Reader) (*Graph, error) {
	scanner := bufio.NewScanner(input)

	nodeAdded := map[string]bool{}
	var graph Graph
	for scanner.Scan() {
		from, to := processLine(scanner.Text())
		graph.Edges = append(graph.Edges, Edge{From: from, To: to})

		if !nodeAdded[from] {
			graph.Modules = append(graph.Modules, ModuleStats{ModuleName: from})
			nodeAdded[from] = true
		}
		if !nodeAdded[to] {
			graph.Modules = append(graph.Modules, ModuleStats{ModuleName: to})
			nodeAdded[to] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("got error from stdout go mod graph scanner: %w", err)
	}

	return &graph, nil
}

// processLine parses single line of go mod graph output
func processLine(line string) (from, to string) {
	vNames := strings.Split(strings.TrimSpace(line), " ")
	if len(vNames) < 2 {
		return "", ""
	}
	return getNameFromVersioned(vNames[0]), getNameFromVersioned(vNames[1])
}

func getNameFromVersioned(versioned string) string {
	parts := strings.Split(versioned, "@")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
