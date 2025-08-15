package graph

import (
	"fmt"

	"github.com/elecbug/go-dspkg/graph/graph_type"
)

// Graph is represented internally using an adjacency matrix and supports
// directed/undirected and weighted/unweighted configurations.
type Graph struct {
	matrix        Matrix
	capacity      int
	updateVersion int
	graphType     graph_type.GraphType
	recentlyID    int
	state         map[int]bool
	nodeCount     int
	edgeCount     int
}

// NodeID identifies a node within a Graph.
type NodeID int

// Degree contains the incoming and outgoing edge counts for a node.
type Degree struct {
	Incoming int
	Outgoing int
}

// NewGraph creates and initializes a new Graph with the given type and capacity.
func NewGraph(graphType graph_type.GraphType, capacity int) *Graph {
	return &Graph{
		matrix:        newMatrix(capacity),
		capacity:      capacity,
		graphType:     graphType,
		updateVersion: 0,
		recentlyID:    0,
		state:         make(map[int]bool, 0),
		nodeCount:     0,
		edgeCount:     0,
	}
}

// AddNode adds a new node to the graph.
// It returns the created node's ID or an error if the graph is at capacity.
func (g *Graph) AddNode() (NodeID, error) {
	if g.recentlyID >= g.capacity {
		return -1, fmt.Errorf("graph is fulled, you set capacity to %d at init", g.capacity)
	}

	g.updateVersion++
	g.state[g.recentlyID] = true
	g.recentlyID++
	g.nodeCount++

	return NodeID(g.recentlyID - 1), nil
}

// RemoveNode removes the node with the given ID.
func (g *Graph) RemoveNode(id NodeID) error {
	if id > NodeID(g.recentlyID) {
		return fmt.Errorf("entered ID %d is bigger then currently node count %d", id, g.recentlyID)
	}

	if !g.state[int(id)] {
		return fmt.Errorf("node (ID: %d) is already removed", id)
	}

	g.updateVersion++
	g.state[int(id)] = false
	g.nodeCount--

	for x, v := range g.matrix {
		for y := range v {
			if x == int(id) || y == int(id) {
				g.matrix[x][y] = INF_DISTANCE

				if g.graphType == graph_type.UNDIRECTED_UNWEIGHTED || g.graphType == graph_type.UNDIRECTED_WEIGHTED {
					g.matrix[y][x] = INF_DISTANCE
				}
			}
		}
	}

	return nil
}

// FindNode reports whether the node with the given ID exists and is active.
func (g *Graph) FindNode(id NodeID) (bool, error) {
	if id > NodeID(g.recentlyID) {
		return false, fmt.Errorf("entered ID %d is bigger then currently node count %d", id, g.recentlyID)
	}

	return g.state[int(id)], nil
}

// AddEdge adds an unweighted edge between two nodes.
func (g *Graph) AddEdge(from, to NodeID) error {
	return g.AddWeightEdge(from, to, 1)
}

// AddWeightEdge adds a weighted edge between two nodes.
func (g *Graph) AddWeightEdge(from, to NodeID, distance Distance) error {
	if (g.graphType == graph_type.DIRECTED_UNWEIGHTED || g.graphType == graph_type.UNDIRECTED_UNWEIGHTED) && distance != 1 {
		return fmt.Errorf("unweighted graph must use AddEdge() or distance be 1")
	}

	if from == to {
		return fmt.Errorf("does not add edge to self")
	}
	if from > NodeID(g.recentlyID) {
		return fmt.Errorf("entered ID %d is bigger then currently node count %d", from, g.recentlyID)
	}
	if to > NodeID(g.recentlyID) {
		return fmt.Errorf("entered ID %d is bigger then currently node count %d", to, g.recentlyID)
	}
	if !g.state[int(from)] {
		return fmt.Errorf("node (ID: %d) is down or removed", from)
	}
	if !g.state[int(to)] {
		return fmt.Errorf("node (ID: %d) is down or removed", to)
	}

	if g.matrix[from][to] != INF_DISTANCE {
		return fmt.Errorf("node (ID: %d) and (ID: %d) are already connected", from, to)
	}

	g.matrix[from][to] = distance

	// Add a reverse edge for undirected graphs.
	if g.graphType == graph_type.UNDIRECTED_UNWEIGHTED || g.graphType == graph_type.UNDIRECTED_WEIGHTED {
		if g.matrix[to][from] != INF_DISTANCE {
			return fmt.Errorf("node (ID: %d) and (ID: %d) are already connected", to, from)
		}

		g.matrix[to][from] = distance
	}

	g.updateVersion++
	g.edgeCount++

	return nil
}

// RemoveEdge removes an edge between two nodes. For undirected graphs, the
// reverse edge is also removed.
func (g *Graph) RemoveEdge(from, to NodeID) error {
	if from == to {
		return fmt.Errorf("does not add edge to self")
	}
	if from > NodeID(g.recentlyID) {
		return fmt.Errorf("entered ID %d is bigger then currently node count %d", from, g.recentlyID)
	}
	if to > NodeID(g.recentlyID) {
		return fmt.Errorf("entered ID %d is bigger then currently node count %d", to, g.recentlyID)
	}
	if !g.state[int(from)] {
		return fmt.Errorf("node (ID: %d) is down or removed", from)
	}
	if !g.state[int(to)] {
		return fmt.Errorf("node (ID: %d) is down or removed", to)
	}

	if g.matrix[from][to] == INF_DISTANCE {
		return fmt.Errorf("node (ID: %d) and (ID: %d) are already disconnected", from, to)
	}

	g.matrix[from][to] = INF_DISTANCE

	// Add a reverse edge for undirected graphs.
	if g.graphType == graph_type.UNDIRECTED_UNWEIGHTED || g.graphType == graph_type.UNDIRECTED_WEIGHTED {
		if g.matrix[to][from] == INF_DISTANCE {
			return fmt.Errorf("node (ID: %d) and (ID: %d) are already disconnected", to, from)
		}

		g.matrix[to][from] = INF_DISTANCE
	}

	g.updateVersion++
	g.edgeCount--

	return nil
}

// FindEdge returns the distance of the edge from -> to.
func (g *Graph) FindEdge(from, to NodeID) (Distance, error) {
	if from == to {
		return INF_DISTANCE, fmt.Errorf("does not add edge to self")
	}
	if from > NodeID(g.recentlyID) {
		return INF_DISTANCE, fmt.Errorf("entered ID %d is bigger then currently node count %d", from, g.recentlyID)
	}
	if to > NodeID(g.recentlyID) {
		return INF_DISTANCE, fmt.Errorf("entered ID %d is bigger then currently node count %d", to, g.recentlyID)
	}
	if !g.state[int(from)] {
		return INF_DISTANCE, fmt.Errorf("node (ID: %d) is down or removed", from)
	}
	if !g.state[int(to)] {
		return INF_DISTANCE, fmt.Errorf("node (ID: %d) is down or removed", to)
	}

	return g.matrix[from][to], nil
}

// Degree returns the degree counts for the node with the given ID.
func (g *Graph) Degree(id NodeID) (*Degree, error) {
	if id > NodeID(g.recentlyID) {
		return nil, fmt.Errorf("entered ID %d is bigger then currently node count %d", id, g.recentlyID)
	}

	if !g.state[int(id)] {
		return nil, fmt.Errorf("node (ID: %d) is already removed", id)
	}

	degree := Degree{
		Incoming: 0,
		Outgoing: 0,
	}

	for i := 0; i < g.recentlyID; i++ {
		if g.matrix[id][i] != INF_DISTANCE {
			degree.Incoming++
		}
		if g.matrix[i][id] != INF_DISTANCE {
			degree.Outgoing++
		}
	}

	return &degree, nil
}

// AliveNodes returns the list of active node IDs.
func (g *Graph) AliveNodes() []NodeID {
	result := make([]NodeID, g.nodeCount)

	idx := 0
	for i, v := range g.state {
		if v {
			result[idx] = NodeID(i)
			idx++
		}
	}

	return result
}

// ToMatrix returns the adjacency matrix representation of the graph.
func (g Graph) ToMatrix() Matrix {
	return g.matrix
}

// String returns a human-readable representation of the adjacency matrix.
func (g Graph) String() string {
	result := ""
	matrix := g.matrix

	// Iterate over each row of the matrix.
	for _, arr := range [][]Distance(matrix) {
		// Iterate over each element in the row.
		for _, a := range arr {
			if a != INF_DISTANCE {
				// Print the distance if it is not `INF`.
				result += fmt.Sprintf("%3d ", a)
			} else {
				// Use "INF" to represent unreachable nodes.
				result += "INF "
			}
		}

		// Add a newline at the end of each row.
		result += "\n"
	}

	return result
}

// NodeCount returns the number of nodes in the graph.
func (g Graph) NodeCount() int {
	return g.nodeCount
}

// EdgeCount returns the number of edges in the graph.
func (g Graph) EdgeCount() int {
	return g.edgeCount
}

// Type returns the graph type (directed/undirected, weighted/unweighted).
func (g Graph) Type() graph_type.GraphType {
	return g.graphType
}

// Version returns the update version used to invalidate cached computations.
func (g Graph) Version() int {
	return g.updateVersion
}
