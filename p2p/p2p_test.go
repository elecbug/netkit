package p2p_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/elecbug/netkit/network-graph/graph/standard_graph"
	"github.com/elecbug/netkit/p2p"
)

func TestGenerateNetwork(t *testing.T) {
	g := standard_graph.ErdosRenyiGraph(1000, 0.005, true)
	t.Logf("Generated graph with %d nodes and %d edges\n", len(g.Nodes()), g.EdgeCount())
	src := rand.NewSource(time.Now().UnixNano())

	nodeLatency := func() float64 { return p2p.LogNormalRand(5.704, 0.5, src) }
	edgeLatency := func() float64 { return p2p.LogNormalRand(5.704, 0.3, src) }

	nw, _ := p2p.GenerateNetwork(g, nodeLatency, edgeLatency)
	t.Logf("Generated network with %d nodes\n", len(nw))
	for id, node := range nw {
		t.Logf("Node %d: latency=%.2fms, edges=%v\n", id, node.Latency, node.Edges)
	}

	p2p.RunNetworkSimulation(nw)
	p2p.Publish(nw[0], "Hello, P2P Network!")

	time.Sleep(5 * time.Second)

	count := 0
	for id, node := range nw {
		c := len(node.SentTo["Hello, P2P Network!"])
		t.Logf("Node %d sent %d/%d\n", id, c, len(node.Edges))
		t.Logf("Node %d data: recv: %v, sent: %v, seen: %v\n",
			id,
			node.RecvFrom["Hello, P2P Network!"],
			node.SentTo["Hello, P2P Network!"],
			node.SeenAt["Hello, P2P Network!"],
		)
		count += c
	}

	t.Logf("Total received count: %d\n", count)
}
