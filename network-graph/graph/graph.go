package graph

import (
	"github.com/elecbug/go-dspkg/network-graph/node"
)

type Graph struct {
	nodes map[node.ID]bool
	edges map[node.ID]map[node.ID]bool
}

func New() *Graph {
	return &Graph{
		nodes: make(map[node.ID]bool),
		edges: make(map[node.ID]map[node.ID]bool),
	}
}
