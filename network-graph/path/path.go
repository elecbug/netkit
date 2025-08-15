package path

import "github.com/elecbug/go-dspkg/network-graph/node"

type Path struct {
	distance int
	nodes    []node.ID
	isInf    bool
}

func NewPath(nodes ...node.ID) *Path {
	return &Path{
		distance: len(nodes) - 1, // Assuming distance is the number of edges
		isInf:    len(nodes) == 0,
		nodes:    nodes,
	}
}
