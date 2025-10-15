package p2p

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
	"github.com/elecbug/netkit/p2p/broadcast"
)

// Config holds configuration parameters for the P2P network.
type Network struct {
	nodes map[ID]*p2pNode
	cfg   *Config
}

// GenerateNetwork creates a P2P network from the given graph.
// nodeLatency and edgeLatency are functions that generate latencies for nodes and edges respectively.
func GenerateNetwork(g *graph.Graph, nodeLatency, edgeLatency func() float64, cfg *Config) (*Network, error) {
	nodes := make(map[ID]*p2pNode)
	maps := make(map[node.ID]ID)

	// create nodes
	for _, gn := range g.Nodes() {
		num, err := strconv.Atoi(gn.String())

		if err != nil {
			return nil, err
		}

		n := newNode(ID(num), nodeLatency())
		n.edges = make(map[ID]p2pEdge)

		nodes[n.id] = n
		maps[gn] = n.id
	}

	for _, gn := range g.Nodes() {
		num, err := strconv.Atoi(gn.String())

		if err != nil {
			return nil, err
		}

		n := nodes[ID(num)]

		for _, neighbor := range g.Neighbors(gn) {
			j := maps[neighbor]

			edge := p2pEdge{
				TargetID: ID(j),
				Latency:  edgeLatency(),
			}

			n.edges[edge.TargetID] = edge
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

// NodeIDs returns a slice of all node IDs in the network.
func (n *Network) NodeIDs() []ID {
	ids := make([]ID, 0, len(n.nodes))

	for id := range n.nodes {
		ids = append(ids, id)
	}

	return ids
}

// Publish sends a message to the specified node's message queue.
func (n *Network) Publish(nodeID ID, msg string, protocol broadcast.Protocol) error {
	if node, ok := n.nodes[nodeID]; ok {
		if !node.alive {
			return fmt.Errorf("node %d is not alive", nodeID)
		}

		node.msgQueue <- Message{From: nodeID, Content: msg, Protocol: protocol}
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
func (n *Network) MessageInfo(nodeID ID, content string) (map[string]any, error) {
	node := n.nodes[nodeID]

	if node == nil {
		return nil, fmt.Errorf("node %d not found", nodeID)
	}

	node.mu.Lock()
	defer node.mu.Unlock()

	info := make(map[string]any)

	info["recv"] = make([]ID, 0)
	for k := range node.recvFrom[content] {
		info["recv"] = append(info["recv"].([]ID), k)
	}

	info["sent"] = make([]ID, 0)
	for k := range node.sentTo[content] {
		info["sent"] = append(info["sent"].([]ID), k)
	}

	info["seen"] = node.seenAt[content].String()

	return info, nil
}
