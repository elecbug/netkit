package p2p

import (
	"context"
	"sync"
	"time"
)

// ID represents a unique identifier for a node in the P2P network.
type ID uint64

// p2pNode represents a node in the P2P network.
type p2pNode struct {
	id          ID
	nodeLatency float64
	edges       map[ID]p2pEdge

	recvFrom map[string]map[ID]struct{} // content -> set of senders
	sentTo   map[string]map[ID]struct{} // content -> set of targets
	seenAt   map[string]time.Time       // content -> first arrival time

	msgQueue chan Message
	mu       sync.Mutex

	alive bool
}

// p2pEdge represents a connection from one node to another in the P2P network.
type p2pEdge struct {
	TargetID ID
	Latency  float64 // in milliseconds
}

// newNode creates a new Node with the given ID and node latency.
func newNode(id ID, nodeLatency float64) *p2pNode {
	return &p2pNode{
		id:          id,
		nodeLatency: nodeLatency,
		edges:       make(map[ID]p2pEdge),

		recvFrom: make(map[string]map[ID]struct{}),
		sentTo:   make(map[string]map[ID]struct{}),
		seenAt:   make(map[string]time.Time),

		msgQueue: make(chan Message, 1000),
		mu:       sync.Mutex{},
	}
}

// eachRun starts the message handling routine for the node.
func (n *p2pNode) eachRun(network *Network, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	go func(ctx context.Context) {
		n.alive = true

		for msg := range n.msgQueue {
			select {
			case <-ctx.Done():
				n.alive = false
				return
			default:
				first := false
				var excludeSnapshot map[ID]struct{}

				n.mu.Lock()
				if _, ok := n.recvFrom[msg.Content]; !ok {
					n.recvFrom[msg.Content] = make(map[ID]struct{})
				}
				n.recvFrom[msg.Content][msg.From] = struct{}{}

				if _, ok := n.seenAt[msg.Content]; !ok {
					n.seenAt[msg.Content] = time.Now()
					first = true
					excludeSnapshot = copyIDSet(n.recvFrom[msg.Content])
				}
				n.mu.Unlock()

				if first {
					go func(msg Message, exclude map[ID]struct{}) {
						time.Sleep(time.Duration(n.nodeLatency) * time.Millisecond)
						n.publish(network, msg, exclude)
					}(msg, excludeSnapshot)
				}
			}
		}
	}(ctx)
}

// copyIDSet creates a shallow copy of a set of IDs.
func copyIDSet(src map[ID]struct{}) map[ID]struct{} {
	dst := make(map[ID]struct{}, len(src))
	for k := range src {
		dst[k] = struct{}{}
	}
	return dst
}

// publish sends the message to neighbors, excluding 'exclude' and already-sent targets.
func (n *p2pNode) publish(network *Network, msg Message, exclude map[ID]struct{}) {
	content := msg.Content
	protocol := msg.Protocol

	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.sentTo[content]; !ok {
		n.sentTo[content] = make(map[ID]struct{})
	}

	willSendEdges := make([]p2pEdge, 0)

	for _, edge := range n.edges {
		if _, wasSender := exclude[edge.TargetID]; wasSender {
			continue
		}
		if _, already := n.sentTo[content][edge.TargetID]; already {
			continue
		}
		if _, received := n.recvFrom[content][edge.TargetID]; received {
			continue
		}
		n.sentTo[content][edge.TargetID] = struct{}{}

		willSendEdges = append(willSendEdges, edge)
	}

	if protocol == Gossiping && len(willSendEdges) > 0 {
		k := int(float64(len(willSendEdges)) * network.cfg.GossipFactor)
		willSendEdges = willSendEdges[:k]
	}

	for _, edge := range willSendEdges {
		edgeCopy := edge

		go func(e p2pEdge) {
			time.Sleep(time.Duration(e.Latency) * time.Millisecond)
			network.nodes[e.TargetID].msgQueue <- Message{From: n.id, Content: content, Protocol: protocol}
		}(edgeCopy)
	}
}
