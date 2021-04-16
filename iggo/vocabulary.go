package iggo

type Edge struct {
	From string
	To   string
}

type Node struct {
	Name string
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}
