package iggo

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Node struct {
	Name string
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}
