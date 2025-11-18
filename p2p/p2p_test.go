package p2p_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/elecbug/netkit/graph/standard_graph"
	"github.com/elecbug/netkit/p2p"
)

func TestGenerateNetwork(t *testing.T) {
	sg := standard_graph.NewStandardGraph()
	g := sg.ErdosRenyiGraph(1000, 50.000/1000, true)
	t.Logf("Generated graph with %d nodes and %d edges\n", len(g.Nodes()), g.EdgeCount())
	src := rand.NewSource(time.Now().UnixNano())

	nodeLatency := func() float64 { return p2p.LogNormalRand(math.Log(100), 0.5, src) }
	edgeLatency := func() float64 { return p2p.LogNormalRand(math.Log(100), 0.3, src) }

	nw, err := p2p.GenerateNetwork(g, nodeLatency, edgeLatency, &p2p.Config{GossipFactor: 0.35})
	if err != nil {
		t.Fatalf("Failed to generate network: %v", err)
	}

	t.Logf("Generated network with %d nodes\n", len(nw.PeerIDs()))

	msg1 := "Hello, P2P World!"
	msg2 := "Goodbye, P2P World!"
	msg3 := "The quick brown fox jumps over the lazy dog."

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nw.RunNetworkSimulation(ctx)

	t.Logf("Publishing message '%s' from node %d\n", msg1, nw.PeerIDs()[0])
	err = nw.Publish(nw.PeerIDs()[0], msg1, p2p.Flooding)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}
	time.Sleep(1000 * time.Millisecond)
	t.Logf("Reachability of message '%s': %f\n", msg1, nw.Reachability(msg1))

	t.Logf("Publishing message '%s' from node %d\n", msg2, nw.PeerIDs()[1])
	err = nw.Publish(nw.PeerIDs()[1], msg2, p2p.Gossiping)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}
	time.Sleep(300 * time.Millisecond)
	cancel()
	time.Sleep(700 * time.Millisecond)
	t.Logf("Reachability of message '%s': %f\n", msg2, nw.Reachability(msg2))

	nw.RunNetworkSimulation(context.Background())
	t.Logf("Publishing message '%s' from node %d\n", msg3, nw.PeerIDs()[2])
	err = nw.Publish(nw.PeerIDs()[2], msg3, p2p.Gossiping)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}
	time.Sleep(1000 * time.Millisecond)
	t.Logf("Reachability of message '%s': %f\n", msg3, nw.Reachability(msg3))

	result := make(map[string]map[string]any)

	for _, nodeID := range nw.PeerIDs() {
		if info, err := nw.MessageInfo(nodeID, msg1); err == nil {
			result[fmt.Sprintf("msg_1-node_%d", nodeID)] = info
		}
		if info, err := nw.MessageInfo(nodeID, msg2); err == nil {
			result[fmt.Sprintf("msg_2-node_%d", nodeID)] = info
		}
		if info, err := nw.MessageInfo(nodeID, msg3); err == nil {
			result[fmt.Sprintf("msg_3-node_%d", nodeID)] = info
		}
	}

	data, _ := json.Marshal(result)

	os.WriteFile("p2p_result.log", data, 0644)
}

func TestMetrics(t *testing.T) {
	sg := standard_graph.NewStandardGraph()
	g := sg.ErdosRenyiGraph(1000, 50.000/1000, true)
	t.Logf("Generated graph with %d nodes and %d edges\n", len(g.Nodes()), g.EdgeCount())
	src := rand.NewSource(time.Now().UnixNano())

	nodeLatency := func() float64 { return p2p.LogNormalRand(math.Log(100), 0.5, src) }
	edgeLatency := func() float64 { return p2p.LogNormalRand(math.Log(100), 0.3, src) }

	nw, err := p2p.GenerateNetwork(g, nodeLatency, edgeLatency, &p2p.Config{GossipFactor: 0.35})
	if err != nil {
		t.Fatalf("Failed to generate network: %v", err)
	}

	t.Logf("Generated network with %d nodes\n", len(nw.PeerIDs()))

	msg1 := "Hello, P2P World!"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nw.RunNetworkSimulation(ctx)

	t.Logf("Publishing message '%s' from node %d\n", msg1, nw.PeerIDs()[0])
	err = nw.Publish(nw.PeerIDs()[0], msg1, p2p.Flooding)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}
	time.Sleep(1000 * time.Millisecond)
	t.Logf("Number of nodes: %d\n", len(nw.PeerIDs()))
	t.Logf("Reachability of message '%s': %f\n", msg1, nw.Reachability(msg1))
	t.Logf("First message reception times of message '%s': %v\n", msg1, nw.FirstMessageReceptionTimes(msg1))
	t.Logf("Number of duplicate messages of message '%s': %d\n", msg1, nw.NumberOfDuplicateMessages(msg1))
}
