package p2p

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/elecbug/netkit/graph"
)

// Config holds configuration parameters for the P2P network.
type Network struct {
	nodes map[PeerID]*p2pNode
	cfg   *Config
}

// GenerateNetwork creates a P2P network from the given graph.
// nodeLatency and edgeLatency are functions that generate latencies for nodes and edges respectively.
func GenerateNetwork(g *graph.Graph, nodeLatency, edgeLatency func() float64, cfg *Config) (*Network, error) {
	nodes := make(map[PeerID]*p2pNode)
	maps := make(map[graph.NodeID]PeerID)

	// create nodes
	for _, gn := range g.Nodes() {
		num, err := strconv.Atoi(gn.String())

		if err != nil {
			return nil, err
		}

		n := newNode(PeerID(num), nodeLatency())
		n.edges = make(map[PeerID]p2pEdge)

		nodes[n.id] = n
		maps[gn] = n.id
	}

	for _, gn := range g.Nodes() {
		num, err := strconv.Atoi(gn.String())

		if err != nil {
			return nil, err
		}

		n := nodes[PeerID(num)]

		for _, neighbor := range g.Neighbors(gn) {
			j := maps[neighbor]

			edge := p2pEdge{
				targetID:    PeerID(j),
				edgeLatency: edgeLatency(),
			}

			n.edges[edge.targetID] = edge
		}
	}

	return &Network{nodes: nodes, cfg: cfg}, nil
}

// RunNetworkSimulation starts the message handling routines for all nodes in the network.
func (n *Network) RunNetworkSimulation(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(len(n.nodes))

	for _, node := range n.nodes {
		node.eachRun(n, wg, ctx)
	}

	wg.Wait()
}

// PeerIDs returns a slice of all node IDs in the network.
func (n *Network) PeerIDs() []PeerID {
	ids := make([]PeerID, 0, len(n.nodes))

	for id := range n.nodes {
		ids = append(ids, id)
	}

	return ids
}

// Publish sends a message to the specified node's message queue.
func (n *Network) Publish(nodeID PeerID, msg string, protocol BroadcastProtocol, customProtocol func(msg Message, known []PeerID, sent []PeerID, received []PeerID, params map[string]any) []PeerID) error {
	if node, ok := n.nodes[nodeID]; ok {
		if !node.alive {
			return fmt.Errorf("node %d is not alive", nodeID)
		}

		node.msgQueue <- Message{From: nodeID, Content: msg, Protocol: protocol, HopCount: 0, CustomProtocol: customProtocol}
		return nil
	}

	return fmt.Errorf("node %d not found", nodeID)
}

// Reachability calculates the fraction of nodes that have received the specified message.
func (n *Network) Reachability(msg string) float64 {
	total := 0
	reached := 0

	for _, node := range n.nodes {
		total++
		node.mu.Lock()
		if _, ok := node.seenAt[msg]; ok {
			reached++
		}
		node.mu.Unlock()
	}

	return float64(reached) / float64(total)
}

// FirstMessageReceptionTimes returns the first reception times of the specified message across all nodes.
func (n *Network) FirstMessageReceptionTimes(msg string) []time.Time {
	firstTimes := make([]time.Time, 0)

	for _, node := range n.nodes {
		node.mu.Lock()
		if t, ok := node.seenAt[msg]; ok {
			firstTimes = append(firstTimes, t)
		}

		node.mu.Unlock()
	}

	return firstTimes
}

// NumberOfDuplicateMessages counts how many duplicate messages were received across all nodes.
func (n *Network) NumberOfDuplicateMessages(msg string) int {
	dupCount := 0

	for _, node := range n.nodes {
		node.mu.Lock()
		if count, ok := node.recvFrom[msg]; ok {
			dupCount += len(count) - 1
		}
		node.mu.Unlock()
	}

	return dupCount
}

// MessageInfo returns a snapshot of the node's message-related information.
func (n *Network) MessageInfo(nodeID PeerID, content string) (map[string]any, error) {
	node := n.nodes[nodeID]

	if node == nil {
		return nil, fmt.Errorf("node %d not found", nodeID)
	}

	node.mu.Lock()
	defer node.mu.Unlock()

	info := make(map[string]any)

	info["recv"] = make([]PeerID, 0)
	for k := range node.recvFrom[content] {
		info["recv"] = append(info["recv"].([]PeerID), k)
	}

	info["sent"] = make([]PeerID, 0)
	for k := range node.sentTo[content] {
		info["sent"] = append(info["sent"].([]PeerID), k)
	}

	info["seen"] = node.seenAt[content].String()

	return info, nil
}
