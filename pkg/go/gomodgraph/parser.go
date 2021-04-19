package gomodgraph

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Edge struct {
	From string
	To   string
}

type Node struct {
	ModuleName string
}

type Graph struct {
	Modules []Node
	Edges   []Edge
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
			graph.Modules = append(graph.Modules, Node{ModuleName: from})
			nodeAdded[from] = true
		}
		if !nodeAdded[to] {
			graph.Modules = append(graph.Modules, Node{ModuleName: to})
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
