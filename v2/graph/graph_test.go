package graph_test

import (
	"fmt"
	"testing"

	"github.com/elecbug/netkit/v2/graph"
)

// TestGraph tests the basic operations of the Graph, including node and edge manipulation,
// string representation, adjacency matrix, serialization/deserialization, and pathfinding.
func TestGraph(t *testing.T) {
	fmt.Println("Test directed, weighted graph")
	dwg := graph.New(true, true)
	testNodeOperations(t, dwg)
	testEdgeOperations(t, dwg, true, true)
	testOtherOperations(t, dwg, true, true)
	fmt.Println("- Verified directed, weighted graph operations")

	fmt.Println("Test directed, unweighted graph")
	dug := graph.New(true, false)
	testNodeOperations(t, dug)
	testEdgeOperations(t, dug, true, false)
	testOtherOperations(t, dug, true, false)
	fmt.Println("- Verified directed, unweighted graph operations")

	fmt.Println("Test undirected, weighted graph")
	uwg := graph.New(false, true)
	testNodeOperations(t, uwg)
	testEdgeOperations(t, uwg, false, true)
	testOtherOperations(t, uwg, false, true)
	fmt.Println("- Verified undirected, weighted graph operations")

	fmt.Println("Test undirected, unweighted graph")
	uug := graph.New(false, false)
	testNodeOperations(t, uug)
	testEdgeOperations(t, uug, false, false)
	testOtherOperations(t, uug, false, false)
	fmt.Println("- Verified undirected, unweighted graph operations")

	testPath(t)
	fmt.Println("- Verified path operations")
}

// testNodeOperations tests node-related operations on the graph.
func testNodeOperations(t *testing.T, g *graph.Graph) {
	var err error

	fmt.Println("- Test node operations")

	// Test AddNode
	if err = g.AddNode("A"); err != nil {
		t.Fatalf("unexpected error adding node A: %v", err)
	}
	if err = g.AddNode("A"); err == nil {
		t.Fatalf("expected error adding duplicate node A, got nil")
	}
	if err = g.AddNode("B"); err != nil {
		t.Fatalf("unexpected error adding node B: %v", err)
	}
	if err = g.AddNode("C"); err != nil {
		t.Fatalf("unexpected error adding node C: %v", err)
	}

	// Test RemoveNode
	if err = g.RemoveNode("D"); err == nil {
		t.Fatalf("expected error removing non-existent node D, got nil")
	}
	if err = g.RemoveNode("C"); err != nil {
		t.Fatalf("unexpected error removing node C: %v", err)
	}

	// Test HasNode
	if !g.HasNode("A") {
		t.Fatalf("expected node A to exist")
	}
	if g.HasNode("C") {
		t.Fatalf("expected node C to not exist")
	}

	// Test Node
	if _, err = g.Node("A"); err != nil {
		t.Fatalf("unexpected error getting node A: %v", err)
	}
	if _, err = g.Node("C"); err == nil {
		t.Fatalf("expected error getting non-existent node C, got nil")
	}

	// Test Nodes
	nodes := g.Nodes()
	expectedNodes := []graph.NodeID{"A", "B"}
	if !equalNodeSlices(nodes, expectedNodes) {
		t.Fatalf("expected nodes %v, got %v", expectedNodes, nodes)
	}

	// Test Size
	if g.Size() != 2 {
		t.Fatalf("expected graph size 2, got %d", g.Size())
	}

	if err = g.AddNode("C"); err != nil {
		t.Fatalf("unexpected error adding node C: %v", err)
	}
	if err = g.AddNode("D"); err != nil {
		t.Fatalf("unexpected error adding node D: %v", err)
	}
	if err = g.AddNode("E"); err != nil {
		t.Fatalf("unexpected error adding node E: %v", err)
	}
}

// testEdgeOperations tests edge-related operations on the graph, including adding, removing, and checking edges, as well as verifying the graph's string representation and adjacency matrix.
func testEdgeOperations(t *testing.T, g *graph.Graph, directed, weighted bool) {
	var err error

	fmt.Println("- Test edge operations")

	// Test AddEdge
	if weighted {
		if err = g.AddEdge("A", "B", nil); err == nil {
			t.Fatalf("expected error adding edge to unweighted graph, got nil")
		}
		if err = g.AddEdge("A", "B", graph.NewWeight(1)); err != nil {
			t.Fatalf("unexpected error adding edge from A to B: %v", err)
		}
		if err = g.AddEdge("A", "B", graph.NewWeight(1)); err == nil {
			t.Fatalf("expected error adding duplicate edge from A to B, got nil")
		}
		if err = g.AddEdge("B", "C", graph.NewWeight(2)); err != nil {
			t.Fatalf("unexpected error adding edge from B to C: %v", err)
		}
		if err = g.AddEdge("C", "D", graph.NewWeight(3)); err != nil {
			t.Fatalf("unexpected error adding edge from C to D: %v", err)
		}
		if err = g.AddEdge("D", "E", graph.NewWeight(4)); err != nil {
			t.Fatalf("unexpected error adding edge from D to E: %v", err)
		}
	} else {
		if err = g.AddEdge("A", "B", nil); err != nil {
			t.Fatalf("unexpected error adding edge from A to B: %v", err)
		}
		if err = g.AddEdge("A", "B", graph.NewWeight(1)); err == nil {
			t.Fatalf("expected error adding edge with weight to unweighted graph, got nil")
		}
		if err = g.AddEdge("A", "B", nil); err == nil {
			t.Fatalf("expected error adding duplicate edge from A to B, got nil")
		}
		if err = g.AddEdge("B", "C", nil); err != nil {
			t.Fatalf("unexpected error adding edge from B to C: %v", err)
		}
		if err = g.AddEdge("C", "D", nil); err != nil {
			t.Fatalf("unexpected error adding edge from C to D: %v", err)
		}
		if err = g.AddEdge("D", "E", nil); err != nil {
			t.Fatalf("unexpected error adding edge from D to E: %v", err)
		}
	}

	// Test RemoveEdge
	if err = g.RemoveEdge("A", "C"); err == nil {
		t.Fatalf("expected error removing non-existent edge from A to C, got nil")
	}
	if err = g.RemoveEdge("A", "B"); err != nil {
		t.Fatalf("unexpected error removing edge from A to B: %v", err)
	}
	if err = g.RemoveEdge("A", "B"); err == nil {
		t.Fatalf("expected error removing non-existent edge from A to B, got nil")
	}

	// Test HasEdge
	if g.HasEdge("A", "B") {
		t.Fatalf("expected edge from A to B to not exist")
	}
	if !g.HasEdge("B", "C") {
		t.Fatalf("expected edge from B to C to exist")
	}

	if !directed {
		if g.HasEdge("B", "A") {
			t.Fatalf("expected edge from B to A to not exist in undirected graph")
		}
		if !g.HasEdge("C", "B") {
			t.Fatalf("expected edge from C to B to exist in undirected graph")
		}
	}

	// Test EdgeWeight
	if _, err = g.EdgeWeight("A", "B"); err == nil {
		t.Fatalf("expected error getting weight of non-existent edge from A to B, got nil")
	}
	if weighted {
		if weight, err := g.EdgeWeight("B", "C"); err != nil {
			t.Fatalf("unexpected error getting weight of edge from B to C: %v", err)
		} else if weight != 2 {
			t.Fatalf("expected weight of edge from B to C to be 2, got %v", weight)
		}

		if !directed {
			if weight, err := g.EdgeWeight("C", "B"); err != nil {
				t.Fatalf("unexpected error getting weight of edge from C to B: %v", err)
			} else if weight != 2 {
				t.Fatalf("expected weight of edge from C to B to be 2 in undirected graph, got %v", weight)
			}
		}
	} else {
		if weight, err := g.EdgeWeight("B", "C"); err != nil {
			t.Fatalf("unexpected error getting weight of edge from B to C: %v", err)
		} else if weight != 1 {
			t.Fatalf("expected weight of edge from B to C to be 1 in unweighted graph, got %v", weight)
		}

		if !directed {
			if weight, err := g.EdgeWeight("C", "B"); err != nil {
				t.Fatalf("unexpected error getting weight of edge from C to B: %v", err)
			} else if weight != 1 {
				t.Fatalf("expected weight of edge from C to B to be 1 in unweighted graph, got %v", weight)
			}
		}
	}
}

// testOtherOperations tests other graph operations such as String, Matrix, and Serialize/Deserialize.
func testOtherOperations(t *testing.T, g *graph.Graph, directed, weighted bool) {
	fmt.Println("- Test other graph operations")

	// Test String
	var expectedString string
	if weighted && directed {
		expectedString = "Graph:\n  Directed: true\n  Weighted: true\n  Nodes:\n    A: map[]\n    B: map[C:2]\n    C: map[D:3]\n    D: map[E:4]\n    E: map[]\n"
	} else if !weighted && directed {
		expectedString = "Graph:\n  Directed: true\n  Weighted: false\n  Nodes:\n    A: map[]\n    B: map[C:1]\n    C: map[D:1]\n    D: map[E:1]\n    E: map[]\n"
	} else if weighted && !directed {
		expectedString = "Graph:\n  Directed: false\n  Weighted: true\n  Nodes:\n    A: map[]\n    B: map[C:2]\n    C: map[B:2 D:3]\n    D: map[C:3 E:4]\n    E: map[D:4]\n"
	} else if !weighted && !directed {
		expectedString = "Graph:\n  Directed: false\n  Weighted: false\n  Nodes:\n    A: map[]\n    B: map[C:1]\n    C: map[B:1 D:1]\n    D: map[C:1 E:1]\n    E: map[D:1]\n"
	}

	if g.String() != expectedString {
		t.Fatalf("expected graph string:\n%s\nGot:\n%s", expectedString, g.String())
	}

	// Test Matrix
	matrix := g.Matrix()

	var expectedMatrix graph.Matrix
	if weighted && directed {
		expectedMatrix = [][]graph.Weight{
			{0, 0, 0, 0, 0},
			{0, 0, 2, 0, 0},
			{0, 0, 0, 3, 0},
			{0, 0, 0, 0, 4},
			{0, 0, 0, 0, 0},
		}
	} else if !weighted && directed {
		expectedMatrix = [][]graph.Weight{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 1, 0},
			{0, 0, 0, 0, 1},
			{0, 0, 0, 0, 0},
		}
	} else if weighted && !directed {
		expectedMatrix = [][]graph.Weight{
			{0, 0, 0, 0, 0},
			{0, 0, 2, 0, 0},
			{0, 2, 0, 3, 0},
			{0, 0, 3, 0, 4},
			{0, 0, 0, 4, 0},
		}
	} else if !weighted && !directed {
		expectedMatrix = [][]graph.Weight{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 0, 1, 0},
			{0, 0, 1, 0, 1},
			{0, 0, 0, 1, 0},
		}
	}

	if !equalWeightMatrices(*matrix, expectedMatrix) {
		t.Fatalf("expected adjacency matrix:\n%v\nGot:\n%v", expectedMatrix, matrix)
	}

	// Test Serialize/Deserialize
	serialized, err := g.Serialize()
	if err != nil {
		t.Fatalf("unexpected error serializing graph: %v", err)
	}

	deserialized, err := graph.Deserialize(serialized)
	if err != nil {
		t.Fatalf("unexpected error deserializing graph: %v", err)
	}

	if deserialized.String() != g.String() {
		t.Fatalf("expected deserialized graph string:\n%s\nGot:\n%s", g.String(), deserialized.String())
	}

	if deserialized.Hash() != g.Hash() {
		t.Fatalf("expected deserialized graph hash %s, got %s", g.Hash(), deserialized.Hash())
	}
}

// testPath tests the Path method of the graph, which calculates the path and total distance between a sequence of nodes.
func testPath(t *testing.T) {
	fmt.Println("Test path operations")
	g := graph.New(true, true)

	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	g.AddNode("D")
	g.AddNode("E")

	g.AddEdge("A", "B", graph.NewWeight(1))
	g.AddEdge("B", "C", graph.NewWeight(2))
	g.AddEdge("C", "D", graph.NewWeight(3.14))
	g.AddEdge("D", "E", graph.NewWeight(4))

	path, err := g.Path("A", "B", "C", "D", "E")

	if err != nil {
		t.Fatalf("unexpected error finding path from A to E: %v", err)
	}

	// Verify path properties
	if path.IsInfinite() {
		t.Fatalf("expected finite path from A to E, got infinite path")
	}

	if !equalNodeSlices(path.Nodes(), []graph.NodeID{"A", "B", "C", "D", "E"}) {
		t.Fatalf("expected path nodes [A B C D E], got %v", path.Nodes())
	}

	if path.TotalDistance() != 10.14 {
		t.Fatalf("expected total distance 10.14, got %v", path.TotalDistance())
	}
}

// equalNodeSlices checks if two slices of NodeIDs contain the same elements, regardless of order.
func equalNodeSlices(a, b []graph.NodeID) bool {
	if len(a) != len(b) {
		return false
	}

	nodeSet := make(map[graph.NodeID]struct{})
	for _, id := range a {
		nodeSet[id] = struct{}{}
	}

	for _, id := range b {
		if _, exists := nodeSet[id]; !exists {
			return false
		}
	}

	return true
}

// equalWeightMatrices checks if two adjacency matrices are equal.
func equalWeightMatrices(a, b graph.Matrix) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}
