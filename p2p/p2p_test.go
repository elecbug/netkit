package p2p_test

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/elecbug/netkit/network-graph/algorithm"
	"github.com/elecbug/netkit/network-graph/graph/standard_graph"
	"github.com/elecbug/netkit/p2p"
)

func TestGenerateNetwork(t *testing.T) {
	g := standard_graph.ErdosRenyiGraph(1000, 0.05, true)
	t.Logf("Generated graph with %d nodes and %d edges\n", len(g.Nodes()), g.EdgeCount())
	src := rand.NewSource(time.Now().UnixNano())

	nodeLatency := func() float64 { return p2p.LogNormalRand(5.704, 0.5, src) }
	edgeLatency := func() float64 { return p2p.LogNormalRand(5.704, 0.3, src) }
	queuingLatency := func() float64 { return p2p.LogNormalRand(5.0, 0.2, src) }

	nw, _ := p2p.GenerateNetwork(g, nodeLatency, edgeLatency, queuingLatency)
	t.Logf("Generated network with %d nodes\n", len(nw))
	for id, node := range nw {
		t.Logf("Node %d: validation_latency=%.2fms, queuing_latency=%.2fms, edges=%v\n", id, node.ValidationLatency, node.QueuingLatency, node.Edges)
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

func TestExpCase(t *testing.T) {
	run := false

	if run {
		for i := 4; i < 10; i++ {
			for j := 0; j < 100; j++ {
				filename := fmt.Sprintf("temp/p2p_result-%02d-%03d.log", i, j)

				if _, err := os.Stat(filename); err == nil {
					t.Logf("File %s already exists, skipping...\n", filename)
					continue
				}

				t.Logf("Experiment case: %02d-%03d\n", i, j)
				r := rand.New(rand.NewSource(time.Now().UnixNano()))

				n := r.Int()%(128-119) + 119
				g := standard_graph.ErdosRenyiGraph(n, float64(i)/float64(n), true)

				nodeLatency := func() float64 { return p2p.LogNormalRand(math.Log(100), 0.1, rand.NewSource(time.Now().UnixNano())) }
				edgeLatency := func() float64 { return p2p.LogNormalRand(math.Log(50), 0.1, rand.NewSource(time.Now().UnixNano())) }
				queuingLatency := func() float64 { return p2p.LogNormalRand(math.Log(50), 0.1, rand.NewSource(time.Now().UnixNano())) }

				nw, _ := p2p.GenerateNetwork(g, nodeLatency, edgeLatency, queuingLatency)
				msg := "Hello, P2P World!"

				t.Logf("Generated graph with %d nodes and %d edges\n", len(g.Nodes()), g.EdgeCount())

				p2p.RunNetworkSimulation(nw)
				p2p.Publish(nw[0], msg)
				time.Sleep(4 * time.Second)

				result := make(map[string]map[string]any)

				for id, node := range nw {
					result[fmt.Sprintf("node_%d", id)] = map[string]any{}
					result[fmt.Sprintf("node_%d", id)]["recv"] = node.RecvFrom[msg]
					result[fmt.Sprintf("node_%d", id)]["sent"] = node.SentTo[msg]
					result[fmt.Sprintf("node_%d", id)]["seen"] = node.SeenAt[msg]
				}

				data, _ := json.Marshal(result)

				os.WriteFile(filename, data, 0644)

				algorithm.CacheClear()
			}
		}
	}
}
