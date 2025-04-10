package graph

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGraph(t *testing.T) {
	g := NewGraph(UNDIRECTED_UNWEIGHTED, 10)

	_, err := g.AddNode("first")

	if err != nil {
		t.Fatal(err)
	}

	g.AddNode("second")

	if err != nil {
		t.Fatal(err)
	}

	g.AddNode("third")

	if err != nil {
		t.Fatal(err)
	}

	g.AddNode("fourth")

	if err != nil {
		t.Fatal(err)
	}

	g.AddNode("fourth")

	if err != nil {
		t.Fatal(err)
	}

	g.AddNode("fifth")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", spew.Sdump(g))

	nodes, err := g.FindNodesByName("fourth")

	if err != nil {
		t.Fatal(err)
	}

	if len(nodes) != 2 {
		t.Fatal("invalid node count")
	}

	g.RemoveNode(nodes[0].ID())

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", spew.Sdump(g))

	firs, err := g.FindNodesByName("first")

	if err != nil {
		t.Fatal(err)
	}

	secs, err := g.FindNodesByName("second")

	if err != nil {
		t.Fatal(err)
	}

	fifs, err := g.FindNodesByName("fifth")

	if err != nil {
		t.Fatal(err)
	}

	if len(firs) != 1 || len(secs) != 1 || len(fifs) != 1 {
		t.Fatal("invalid node count")
	}

	err = g.AddEdge(firs[0].ID(), secs[0].ID())

	if err != nil {
		t.Fatal(err)
	}

	err = g.AddEdge(firs[0].ID(), fifs[0].ID())

	if err != nil {
		t.Fatal(err)
	}

	err = g.AddEdge(fifs[0].ID(), secs[0].ID())

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", spew.Sdump(g))

	t.Logf("%s\n", spew.Sdump(g.Matrix()))
}
