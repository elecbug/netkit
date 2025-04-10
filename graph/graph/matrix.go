package graph

// `Matrix` represents the adjacency matrix of a graph.
// Each element in the matrix corresponds to the distance between two nodes.
// If two nodes are not directly connected, the value is set to `INF`.
type Matrix [][]Distance
