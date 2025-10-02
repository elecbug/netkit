package p2p_test

import (
	"testing"

	"github.com/elecbug/netkit/network-graph/graph/standard_graph"
	"github.com/elecbug/netkit/p2p"
)

func TestGenerateNetwork(t *testing.T) {
	g := standard_graph.ErdosRenyiGraph(1000, 0.01, true)

	nodeLatency := func() float64 { return p2p.LogNormalRand(3, 0.5, nil) }
	edgeLatency := func() float64 { return p2p.LogNormalRand(2, 0.3, nil) }

	nw := p2p.GenerateNetwork(g, nodeLatency, edgeLatency)
	p2p.RunNetworkSimulation(nw)
	p2p.Publish(nw[0], "Hello, P2P Network!")
}
