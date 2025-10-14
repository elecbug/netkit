package p2p

import (
	"strconv"
	"sync"

	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
)

// GenerateNetwork creates a P2P network from the given graph.
// nodeLatency and edgeLatency are functions that generate latencies for nodes and edges respectively.
func GenerateNetwork(g *graph.Graph, nodeLatency, edgeLatency, queuingLatency func() float64) (map[ID]*Node, error) {
	nodes := make(map[ID]*Node)
	maps := make(map[node.ID]ID)

	// create nodes
	for _, gn := range g.Nodes() {
		num, err := strconv.Atoi(gn.String())

		if err != nil {
			return nil, err
		}

		n := &Node{
			ID:                ID(num),
			ValidationLatency: nodeLatency(),
			Edges:             make(map[ID]Edge),
		}

		nodes[n.ID] = n
		maps[gn] = n.ID
	}

	for _, gn := range g.Nodes() {
		num, err := strconv.Atoi(gn.String())

		if err != nil {
			return nil, err
		}

		n := nodes[ID(num)]

		for _, neighbor := range g.Neighbors(gn) {
			j := maps[neighbor]

			edge := Edge{
				TargetID: ID(j),
				Latency:  edgeLatency(),
			}

			n.Edges[edge.TargetID] = edge
		}
	}

	return nodes, nil
}

// RunNetworkSimulation starts the message handling routines for all nodes in the network.
func RunNetworkSimulation(nodes map[ID]*Node) {
	wg := &sync.WaitGroup{}
	wg.Add(len(nodes))

	for _, n := range nodes {
		n.eachRun(nodes, wg)
	}

	wg.Wait()
}

// Publish sends a message to the specified node's message queue.
func Publish(node *Node, msg string) {
	node.msgQueue <- Message{From: node.ID, Content: msg}
}
