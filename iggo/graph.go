package iggo

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// GoProcessGraphBuilder constructs graph using Go toolchain
type GoProcessGraphBuilder struct{}

// FetchGraph loads graph using go in separate process
func (c *GoProcessGraphBuilder) BuildGraph() (*Graph, error) {
	cmd := exec.Command("go", "mod", "graph")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("can not get stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("can not start go command: %w", err)
	}
	scanner := bufio.NewScanner(stdout)

	var nodeAdded map[string]bool
	var graph Graph
	for scanner.Scan() {
		from, to := processLine(scanner.Text())
		graph.Edges = append(graph.Edges, Edge{From: from, To: to})
		if !nodeAdded[from] {
			graph.Nodes = append(graph.Nodes, Node{Name: from})
		}
		if !nodeAdded[to] {
			graph.Nodes = append(graph.Nodes, Node{Name: to})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("got error from stdout go mod graph scanner: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command did not finish successfully: %w", err)
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
