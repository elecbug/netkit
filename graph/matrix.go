package graph

// `Matrix` represents the adjacency matrix of a graph.
// Each element in the matrix corresponds to the distance between two nodes.
// If two nodes are not directly connected, the value is set to `INF`.
type Matrix [][]Distance

func newMatrix(nodeCount int) Matrix {
	result := make([][]Distance, nodeCount)

	for i := range result {
		result[i] = make([]Distance, nodeCount)

		for j := range result[i] {
			result[i][j] = -1
		}
	}

	return Matrix(result)
}
