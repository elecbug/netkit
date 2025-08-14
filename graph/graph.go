package graph

import (
	"fmt"

	"github.com/elecbug/go-dspkg/graph/graph_type"
)

// `Graph` type expressed internally through matrix operations.
// This can express weighted/unweighted and directed/undirected, etc.
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

// `NodeID` wrapper type.
type NodeID int

// `Degree` means node's incoming/outgoing edge count.
type Degree struct {
	Incoming int
	Outgoing int
}

// `NewGraph` creates and initializes a new Graph instance.
// Returns a pointer to the newly created Graph.
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

// `AddNode` adds a new node to the graph.
// Returns the newly created node's ID and an error if insertion fails.
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

// `RemoveNode` removes a node from the graph using its ID.
// Returns an error if the node does not exist.
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
			}
		}
	}

	return nil
}

// `FindNode` retrieves a node from the graph by its ID.
// Returns the Node and an error if the node does not exist.
func (g *Graph) FindNode(id NodeID) (bool, error) {
	if id > NodeID(g.recentlyID) {
		return false, fmt.Errorf("entered ID %d is bigger then currently node count %d", id, g.recentlyID)
	}

	return g.state[int(id)], nil
}

// `AddEdge` adds an unweighted edge between two nodes in the graph.
// Returns an error if the edge cannot be added.
func (g *Graph) AddEdge(from, to NodeID) error {
	return g.AddWeightEdge(from, to, 1)
}

// `AddWeightEdge` adds a weighted edge between two nodes in the graph.
// Returns an error if the edge cannot be added.
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

// `RemoveEdge` removes an edge between two nodes in the graph.
// For undirected graphs, the reverse edge is also removed.
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

// `FindEdge` searches for an edge between two nodes in the graph and returns its distance.
// `Distance` of the edge if it exists.
// An error if the edge or either of the nodes does not exist, or if attempting to find a self-loop edge.
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

// `Degree` returns degree of node from the graph by its ID.
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

// `AliveNodes` return list of all activated node
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

// `Matrix` converts the graph to an adjacency matrix representation.
// Returns a Matrix where each element represents the distance between two nodes.
func (g Graph) ToMatrix() Matrix {
	return g.matrix
}

// `String` returns a string representation of the Matrix.
// This method formats the matrix for easy readability:
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

// `NodeCount` returns the number of nodes in the graph.
func (g Graph) NodeCount() int {
	return g.nodeCount
}

// `EdgeCount` returns the number of edges in the graph.
func (g Graph) EdgeCount() int {
	return g.edgeCount
}

// `Type` returns the type of the graph (e.g., directed/undirected, weighted/unweighted).
func (g Graph) Type() graph_type.GraphType {
	return g.graphType
}

// `Version` returns whether the graph has been updated since the last algorithmic computation.
func (g Graph) Version() int {
	return g.updateVersion
}
