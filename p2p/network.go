package p2p

import (
	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
)

// GenerateNetwork creates a P2P network from the given graph.
// nodeLatency and edgeLatency are functions that generate latencies for nodes and edges respectively.
func GenerateNetwork(g *graph.Graph, nodeLatency, edgeLatency func() float64) map[ID]*Node {
	nodes := make(map[ID]*Node)
	maps := make(map[node.ID]ID)

	// create nodes
	for i, gn := range g.Nodes() {
		n := &Node{
			ID:      ID(i),
			Latency: nodeLatency(),
			Edges:   make(map[ID]Edge),
		}

		nodes[n.ID] = n
		maps[gn] = n.ID
	}

	for i, gn := range g.Nodes() {
		n := nodes[ID(i)]

		for _, neighbor := range g.Neighbors(gn) {
			j := maps[neighbor]

			edge := Edge{
				TargetID: ID(j),
				Latency:  edgeLatency(),
			}

			n.Edges[edge.TargetID] = edge
		}
	}

	return nodes
}

// RunNetworkSimulation starts the message handling routines for all nodes in the network.
func RunNetworkSimulation(nodes map[ID]*Node) {
	for _, n := range nodes {
		n.eachRun(nodes)
	}
}

// Publish sends a message to the specified node's message queue.
func Publish(node *Node, msg Message) {
	node.msgQueue <- msg
}
