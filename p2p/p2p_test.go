package p2p_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
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

	msg := "Hello, P2P World!"

	p2p.RunNetworkSimulation(nw)
	p2p.Publish(nw[0], msg)

	time.Sleep(5 * time.Second)

	count := 0
	result := make(map[string]map[string]any)

	for id, node := range nw {
		c := len(node.SentTo[msg])
		t.Logf("Node %d sent %d/%d\n", id, c, len(node.Edges))

		result[fmt.Sprintf("node_%d", id)] = map[string]any{}
		result[fmt.Sprintf("node_%d", id)]["recv"] = node.RecvFrom[msg]
		result[fmt.Sprintf("node_%d", id)]["sent"] = node.SentTo[msg]
		result[fmt.Sprintf("node_%d", id)]["seen"] = node.SeenAt[msg]

		count += c
	}

	t.Logf("Total received count: %d\n", count)

	data, _ := json.Marshal(result)

	os.WriteFile("p2p_result.log", data, 0644)
}
