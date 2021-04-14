package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
)

type GoCollector struct{}

func processLine(line string) (from, to string) {
	vals := strings.Split(strings.TrimSpace(line), " ")
	if len(vals) < 2 {
		return "", ""
	}
	return getNameFromVersioned(vals[0]), getNameFromVersioned(vals[1])
}

func getNameFromVersioned(versioned string) string {
	vals := strings.Split(versioned, "@")
	return vals[0]
}

func (c *GoCollector) FetchGraphInvoke() (*Graph, error) {
	cmd := exec.Command("go", "mod", "graph")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
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
		log.Println(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return &graph, nil
}
