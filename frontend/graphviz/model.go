package graphviz

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"go.uber.org/multierr"
)

// Node is generic graph node
type Node map[string]interface{}

// IsValid checks that node is valid
func (n Node) IsValid() bool {
	if n == nil {
		return false
	}
	if _, ok := n["id"]; !ok {
		return false
	}
	// TODO: check that value in map is scalar
	return true
}

// Edge is generic graph edge
type Edge map[string]interface{}

// IsValid checks that edge is valid
func (e Edge) IsValid() bool {
	if e == nil {
		return false
	}
	if _, ok := e["from"]; !ok {
		return false
	}
	if _, ok := e["to"]; !ok {
		return false
	}
	// TODO: check that value in map is scalar
	return true
}

// Graph is generic graph structure
type Graph struct {
	Nodes []Node
	Edges []Edge
}

// NewGraphFromJSONLReader parses JSONL from reader into a generic graph
func NewGraphFromJSONLReader(r io.Reader) (g Graph, finalErr error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()

		var node Node
		var nodeErr error
		if nodeErr = json.Unmarshal(line, &node); nodeErr == nil && node.IsValid() {
			g.Nodes = append(g.Nodes, node)
		}

		var edge Edge
		var edgeErr error
		if edgeErr = json.Unmarshal(line, &edge); edgeErr == nil && edge.IsValid() {
			g.Edges = append(g.Edges, edge)
		}

		// can not get either, keep errors
		if !node.IsValid() && !edge.IsValid() {
			finalErr = multierr.Combine(finalErr, errors.New("both edge and node are invalid"))
		}
		if nodeErr != nil && edgeErr != nil {
			finalErr = multierr.Combine(finalErr, nodeErr, edgeErr)
		}
	}

	if err := scanner.Err(); err != nil {
		finalErr = multierr.Combine(finalErr, fmt.Errorf("got error from stdout from scanner: %w", err))
	}

	return
}
