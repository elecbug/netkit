package p2p

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
)

// Message represents a message sent between nodes in the P2P network.
type Message struct {
	From     ID
	Content  string
	Protocol BroadcastProtocol
}

// Config holds configuration parameters for the P2P network.
type Config struct {
	GossipFactor float64 // fraction of neighbors to gossip to
}

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

// GetNode retrieves a node by its ID.
func (n *Network) GetNode(id ID) *p2pNode {
	return n.nodes[id]
}

// Publish sends a message to the specified node's message queue.
func (n *Network) Publish(nodeID ID, msg string, protocol BroadcastProtocol) error {
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
