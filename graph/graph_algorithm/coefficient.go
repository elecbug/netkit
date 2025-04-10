package graph_algorithm

import (
	"sync"

	"github.com/elecbug/go-dspkg/graph"
)

// `ClusteringCoefficient` computes the local and global clustering coefficients for a graph using a `Unit`.
// Local clustering coefficient measures the degree to which nodes in a graph cluster together.
// Global clustering coefficient is the average of local coefficients across all nodes.
// Return (local, global) clustering coefficient.
func (u *Unit) ClusteringCoefficient() (map[graph.NodeID]float64, float64) {
	g := u.graph
	matrix := g.ToMatrix() // Get adjacency matrix representation of the graph.
	n := len(matrix)       // Number of nodes in the graph.

	// Map to store local clustering coefficients for each node.
	localCoeffs := make(map[graph.NodeID]float64)
	globalSum := 0.0 // Sum of all local clustering coefficients for computing the global coefficient.

	// Iterate over all nodes to calculate local coefficients.
	for v := 0; v < n; v++ {
		neighbors := []int{}

		// Identify neighbors of the current node.
		for i := 0; i < n; i++ {
			if matrix[v][i] != graph.INF_DISTANCE && matrix[v][i] > 0 {
				neighbors = append(neighbors, i)
			}
		}

		k := len(neighbors) // Degree of the node.
		if k < 2 {
			// If a node has fewer than 2 neighbors, its clustering coefficient is 0.
			localCoeffs[graph.NodeID(v)] = 0.0
			continue
		}

		// Count the number of edges between neighbors.
		e := 0
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				if matrix[neighbors[i]][neighbors[j]] != graph.INF_DISTANCE && matrix[neighbors[i]][neighbors[j]] > 0 {
					e++
				}
			}
		}

		// Compute the local clustering coefficient.
		Cv := float64(2*e) / float64(k*(k-1))
		localCoeffs[graph.NodeID(v)] = Cv
		globalSum += Cv
	}

	// Compute the global clustering coefficient (average of local coefficients).
	globalCoeff := globalSum / float64(n)

	return localCoeffs, globalCoeff
}

// `ClusteringCoefficient` computes the local and global clustering coefficients for a graph using a `ParallelUnit`.
// The computation is performed in parallel for better performance on larger graphs.
// Return (local, global) clustering coefficient.
func (pu *ParallelUnit) ClusteringCoefficient() (map[graph.NodeID]float64, float64) {
	g := pu.graph
	matrix := g.ToMatrix() // Get adjacency matrix representation of the graph.
	n := len(matrix)       // Number of nodes in the graph.

	// Map to store local clustering coefficients for each node.
	localCoeffs := make(map[graph.NodeID]float64)
	globalSum := float64(0)

	// Type to store intermediate results from goroutines.
	type result struct {
		node  graph.NodeID
		value float64
	}

	resultChan := make(chan result, n) // Channel for goroutine results.
	var wg sync.WaitGroup              // WaitGroup to synchronize goroutines.

	// Launch a goroutine for each node to compute its local clustering coefficient.
	for v := 0; v < n; v++ {
		wg.Add(1)
		go func(node int) {
			defer wg.Done()

			neighbors := []int{}

			// Identify neighbors of the current node.
			for i := 0; i < n; i++ {
				if matrix[node][i] != graph.INF_DISTANCE && matrix[node][i] > 0 {
					neighbors = append(neighbors, i)
				}
			}

			k := len(neighbors) // Degree of the node.
			if k < 2 {
				// If a node has fewer than 2 neighbors, its clustering coefficient is 0.
				resultChan <- result{node: graph.NodeID(node), value: 0.0}
				return
			}

			// Count the number of edges between neighbors.
			e := 0
			for i := 0; i < k; i++ {
				for j := i + 1; j < k; j++ {
					if matrix[neighbors[i]][neighbors[j]] != graph.INF_DISTANCE && matrix[neighbors[i]][neighbors[j]] > 0 {
						e++
					}
				}
			}

			// Compute the local clustering coefficient.
			Cv := float64(2*e) / float64(k*(k-1))
			resultChan <- result{node: graph.NodeID(node), value: Cv}
		}(v)
	}

	// Close the result channel after all goroutines complete.
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Aggregate results from goroutines.
	for res := range resultChan {
		localCoeffs[res.node] = res.value
		globalSum += res.value
	}

	// Compute the global clustering coefficient (average of local coefficients).
	globalCoeff := globalSum / float64(n)

	return localCoeffs, globalCoeff
}

// `RichClubCoefficient` computes the rich club coefficient for a given threshold degree `k`.
// This coefficient measures how well nodes with `k <= degree` are connected to each other.
func (u *Unit) RichClubCoefficient(k int) float64 {
	g := u.graph
	matrix := g.ToMatrix() // Get adjacency matrix representation of the graph.
	n := len(matrix)       // Number of nodes in the graph.

	// Identify nodes with degree >= k
	nodes := []int{}
	for v := 0; v < n; v++ {
		degree := 0
		for i := 0; i < n; i++ {
			if matrix[v][i] != graph.INF_DISTANCE && matrix[v][i] > 0 {
				degree++
			}
		}
		if degree >= k {
			nodes = append(nodes, v)
		}
	}

	Nk := len(nodes) // Number of nodes with degree >= k
	if Nk < 2 {
		// If there are fewer than 2 nodes, the rich club coefficient is undefined (0).
		return 0.0
	}

	// Count the number of edges between these nodes
	Ek := 0
	for i := 0; i < Nk; i++ {
		for j := i + 1; j < Nk; j++ {
			if matrix[nodes[i]][nodes[j]] != graph.INF_DISTANCE && matrix[nodes[i]][nodes[j]] > 0 {
				Ek++
			}
		}
	}

	// Compute the rich club coefficient
	return float64(2*Ek) / float64(Nk*(Nk-1))
}

// `RichClubCoefficient` computes the rich club coefficient for a given threshold degree k using a `ParallelUnit`.
// The computation is performed in parallel for better performance.
func (pu *ParallelUnit) RichClubCoefficient(k int) float64 {
	g := pu.graph
	matrix := g.ToMatrix() // Get adjacency matrix representation of the graph.
	n := len(matrix)       // Number of nodes in the graph.

	// Identify nodes with degree >= k in parallel
	nodesChan := make(chan int, n)
	var wg sync.WaitGroup

	for v := 0; v < n; v++ {
		wg.Add(1)
		go func(node int) {
			defer wg.Done()
			degree := 0
			for i := 0; i < n; i++ {
				if matrix[node][i] != graph.INF_DISTANCE && matrix[node][i] > 0 {
					degree++
				}
			}
			if degree >= k {
				nodesChan <- node
			}
		}(v)
	}

	// Close channel after goroutines finish
	go func() {
		wg.Wait()
		close(nodesChan)
	}()

	// Collect nodes with degree >= k
	nodes := []int{}
	for node := range nodesChan {
		nodes = append(nodes, node)
	}

	Nk := len(nodes) // Number of nodes with degree >= k
	if Nk < 2 {
		// If there are fewer than 2 nodes, the rich club coefficient is undefined (0).
		return 0.0
	}

	// Count the number of edges between these nodes in parallel
	EkChan := make(chan int, Nk*Nk)
	for i := 0; i < Nk; i++ {
		for j := i + 1; j < Nk; j++ {
			wg.Add(1)
			go func(node1, node2 int) {
				defer wg.Done()
				if matrix[node1][node2] != graph.INF_DISTANCE && matrix[node1][node2] > 0 {
					EkChan <- 1
				}
			}(nodes[i], nodes[j])
		}
	}

	// Close edge channel after goroutines finish
	go func() {
		wg.Wait()
		close(EkChan)
	}()

	// Sum up the edges
	Ek := 0
	for edge := range EkChan {
		Ek += edge
	}

	// Compute the rich club coefficient
	return float64(2*Ek) / float64(Nk*(Nk-1))
}
