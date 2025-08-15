package graph_algorithm

import (
	"math"
	"sync"

	"github.com/elecbug/go-dspkg/graph"
	"github.com/elecbug/go-dspkg/graph/graph_type"
)

// BetweennessCentrality computes the betweenness centrality of each node for a Unit.
// Betweenness centrality measures how often a node appears on the shortest paths between pairs of other nodes.
func (u *Unit) BetweennessCentrality() map[graph.NodeID]float64 {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	centrality := make(map[graph.NodeID]float64)

	// Initialize centrality scores for all nodes to 0.
	for i := 0; i < g.NodeCount(); i++ {
		centrality[graph.NodeID(i)] = 0
	}

	// Count how many times each node appears on the shortest paths.
	for _, path := range u.shortestPaths {
		nodes := path.Nodes()

		for _, n := range nodes {
			// Exclude the source and target nodes of the path.
			if n != nodes[0] && n != nodes[len(nodes)-1] {
				centrality[n]++
			}
		}
	}

	// Normalize the centrality scores.
	n := g.NodeCount()
	if n > 2 {
		for node := range centrality {
			centrality[node] /= float64((n - 1) * (n - 2))
		}
	}

	return centrality
}

// BetweennessCentrality computes the betweenness centrality of each node for a ParallelUnit.
// The computation is performed in parallel for better performance on larger graphs.
func (pu *ParallelUnit) BetweennessCentrality() map[graph.NodeID]float64 {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	centrality := make(map[graph.NodeID]float64)

	// Initialize centrality scores for all nodes to 0.
	for i := 0; i < g.NodeCount(); i++ {
		centrality[graph.NodeID(i)] = 0
	}

	// Define a result type to collect intermediate centrality counts.
	type result struct {
		node  graph.NodeID
		count float64
	}

	resultChan := make(chan result, g.NodeCount())
	var wg sync.WaitGroup

	// Compute centrality scores in parallel.
	for _, path := range pu.shortestPaths {
		wg.Add(1)

		go func(path Path) {
			defer wg.Done()
			nodes := path.Nodes()

			for _, n := range nodes {
				// Exclude the source and target nodes of the path.
				if n != nodes[0] && n != nodes[len(nodes)-1] {
					resultChan <- result{node: n, count: 1}
				}
			}
		}(path)
	}

	// Close the result channel after all goroutines complete.
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Aggregate results from the result channel.
	for res := range resultChan {
		centrality[res.node] += res.count
	}

	// Normalize the centrality scores.
	n := g.NodeCount()
	if n > 2 {
		for node := range centrality {
			centrality[node] /= float64((n - 1) * (n - 2))
		}
	}

	return centrality
}

// DegreeCentrality computes the degree centrality of each node for a Unit.
// Degree centrality is the number of direct connections a node has to other nodes.
func (u *Unit) DegreeCentrality() map[graph.NodeID]float64 {
	g := u.graph
	centrality := make(map[graph.NodeID]float64)
	matrix := g.ToMatrix()
	graphType := g.Type()

	// Initialize centrality scores for all nodes to 0.
	for i := 0; i < g.NodeCount(); i++ {
		centrality[graph.NodeID(i)] = 0
	}

	// Calculate degree based on graph type.
	for i, row := range matrix {
		for _, value := range row {
			if graphType == graph_type.DIRECTED_UNWEIGHTED || graphType == graph_type.UNDIRECTED_UNWEIGHTED {
				if value == 1 {
					centrality[graph.NodeID(i)]++
				}
			} else {
				if value != graph.INF_DISTANCE {
					centrality[graph.NodeID(i)]++
				}
			}
		}
	}

	n := g.NodeCount()
	if n > 1 {
		for node := range centrality {
			centrality[node] /= float64(n - 1)
		}
	}

	return centrality
}

// DegreeCentrality computes the degree centrality of each node for a ParallelUnit.
// The computation is performed in parallel for better performance on larger graphs.
func (pu *ParallelUnit) DegreeCentrality() map[graph.NodeID]float64 {
	g := pu.graph
	centrality := make(map[graph.NodeID]float64)
	matrix := g.ToMatrix()
	graphType := g.Type()

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		node  graph.NodeID
		count float64
	}, g.NodeCount())

	// Worker goroutines to compute degree centrality.
	for i := 0; i < len(matrix); i++ {
		wg.Add(1)
		go func(nodeIndex int) {
			defer wg.Done()
			count := 0.0
			row := matrix[nodeIndex]
			for _, value := range row {
				if graphType == graph_type.DIRECTED_UNWEIGHTED || graphType == graph_type.UNDIRECTED_UNWEIGHTED {
					if value == 1 {
						count++
					}
				} else {
					if value != graph.INF_DISTANCE {
						count++
					}
				}
			}
			resultChan <- struct {
				node  graph.NodeID
				count float64
			}{node: graph.NodeID(nodeIndex), count: count}
		}(i)
	}

	// Close the result channel after all workers finish.
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Aggregate results.
	for res := range resultChan {
		centrality[res.node] = res.count
	}

	n := g.NodeCount()
	if n > 1 {
		for node := range centrality {
			centrality[node] /= float64(n - 1)
		}
	}

	return centrality
}

// EigenvectorCentrality computes the eigenvector centrality of each node for a Unit.
// Eigenvector centrality assigns scores to nodes based on the importance of their neighbors.
func (u *Unit) EigenvectorCentrality(maxIter int, tol float64) map[graph.NodeID]float64 {
	g := u.graph
	matrix := g.ToMatrix()
	n := len(matrix)

	// Initialize centrality scores with 1/n
	centrality := make([]float64, n)
	for i := 0; i < n; i++ {
		centrality[i] = 1.0 / float64(n)
	}

	for iter := 0; iter < maxIter; iter++ {
		newCentrality := make([]float64, n)

		// Update centrality scores
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if matrix[i][j] != graph.INF_DISTANCE {
					newCentrality[i] += float64(matrix[i][j]) * centrality[j]
				}
			}
		}

		// Normalize the new centrality scores
		norm := 0.0
		for _, value := range newCentrality {
			norm += value * value
		}
		norm = math.Sqrt(norm)

		for i := 0; i < n; i++ {
			newCentrality[i] /= norm
		}

		// Check for convergence
		diff := 0.0
		for i := 0; i < n; i++ {
			diff += math.Abs(newCentrality[i] - centrality[i])
		}

		if diff < tol {
			break
		}

		centrality = newCentrality
	}

	// Convert to map for output
	result := make(map[graph.NodeID]float64)
	for i := 0; i < n; i++ {
		result[graph.NodeID(i)] = centrality[i]
	}

	return result
}

// EigenvectorCentrality computes the eigenvector centrality of each node for a ParallelUnit.
// The computation is performed in parallel for better performance on larger graphs.
func (pu *ParallelUnit) EigenvectorCentrality(maxIter int, tol float64) map[graph.NodeID]float64 {
	g := pu.graph
	matrix := g.ToMatrix()
	n := len(matrix)

	// Initialize centrality scores with 1/n
	centrality := make([]float64, n)
	for i := 0; i < n; i++ {
		centrality[i] = 1.0 / float64(n)
	}

	for iter := 0; iter < maxIter; iter++ {
		newCentrality := make([]float64, n)

		var wg sync.WaitGroup

		// Update centrality scores in parallel
		for i := 0; i < n; i++ {
			wg.Add(1)

			go func(node int) {
				defer wg.Done()
				for j := 0; j < n; j++ {
					if matrix[node][j] != graph.INF_DISTANCE {
						newCentrality[node] += float64(matrix[node][j]) * centrality[j]
					}
				}
			}(i)
		}

		wg.Wait()

		// Normalize the new centrality scores
		norm := 0.0
		for _, value := range newCentrality {
			norm += value * value
		}
		norm = math.Sqrt(norm)

		for i := 0; i < n; i++ {
			newCentrality[i] /= norm
		}

		// Check for convergence
		diff := 0.0
		for i := 0; i < n; i++ {
			diff += math.Abs(newCentrality[i] - centrality[i])
		}

		if diff < tol {
			break
		}

		centrality = newCentrality
	}

	// Convert to map for output
	result := make(map[graph.NodeID]float64)
	for i := 0; i < n; i++ {
		result[graph.NodeID(i)] = centrality[i]
	}

	return result
}
