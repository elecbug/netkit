package graph

import (
	"testing"

	"github.com/elecbug/go-dspkg/graph/graph_type"
)

func TestGraph(t *testing.T) {
	g := NewGraph(graph_type.UNDIRECTED_UNWEIGHTED, 20)

	for i := 0; i < 10; i++ {
		_, err := g.AddNode()

		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("--Add 10 nodes--")
	t.Logf("\n%s", g.String())
	t.Logf("Node: %d, Edge: %d", g.nodeCount, g.edgeCount)

	g.RemoveNode(NodeID(2))
	g.RemoveNode(NodeID(6))

	t.Log("--Remove 2 nodes--")
	t.Logf("\n%s", g.String())
	t.Logf("Node: %d, Edge: %d", g.nodeCount, g.edgeCount)

	g.AddEdge(0, 9)
	g.AddEdge(1, 5)
	g.AddEdge(4, 7)
	g.AddEdge(5, 9)

	t.Log("--Add 4 Edges--")
	t.Logf("\n%s", g.String())
	t.Logf("Node: %d, Edge: %d", g.nodeCount, g.edgeCount)

	a, _ := g.FindNode(NodeID(1))
	t.Logf("Find Node 1: %v", a)

	b, _ := g.FindNode(NodeID(2))
	t.Logf("Find Node 2: %v", b)

	c, _ := g.FindEdge(0, 8)
	t.Logf("Find Edge 0-->8: %v", c)

	d, _ := g.FindEdge(0, 9)
	t.Logf("Find Edge 0-->9: %v", d)

	_, ee := g.FindNode(14)
	if ee != nil {
		t.Logf("Find Node 14: %v", ee)
	}
	_, ff := g.FindEdge(0, 15)
	if ff != nil {
		t.Logf("Find Edge 0-->15: %v", ff)
	}
}
