package p2p

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// p2pNode represents a node in the P2P network.
type p2pNode struct {
	id          PeerID
	nodeLatency float64
	edges       map[PeerID]p2pEdge

	recvFrom map[string]map[PeerID]struct{} // content -> set of senders
	sentTo   map[string]map[PeerID]struct{} // content -> set of targets
	seenAt   map[string]time.Time           // content -> first arrival time

	msgQueue chan Message
	mu       sync.Mutex

	alive bool
}

// p2pEdge represents a connection from one node to another in the P2P network.
type p2pEdge struct {
	targetID    PeerID
	edgeLatency float64 // in milliseconds
}

// newNode creates a new Node with the given ID and node latency.
func newNode(id PeerID, nodeLatency float64) *p2pNode {
	return &p2pNode{
		id:          id,
		nodeLatency: nodeLatency,
		edges:       make(map[PeerID]p2pEdge),

		recvFrom: make(map[string]map[PeerID]struct{}),
		sentTo:   make(map[string]map[PeerID]struct{}),
		seenAt:   make(map[string]time.Time),

		msgQueue: make(chan Message, 1000),
		mu:       sync.Mutex{},
	}
}

// eachRun starts the message handling routine for the node.
func (n *p2pNode) eachRun(network *Network, wg *sync.WaitGroup, ctx context.Context) {
	go func(ctx context.Context, wg *sync.WaitGroup) {
		n.alive = true
		wg.Done()

		for msg := range n.msgQueue {
			select {
			case <-ctx.Done():
				n.alive = false
				return
			default:
				first := false

				n.mu.Lock()
				if _, ok := n.recvFrom[msg.Content]; !ok {
					n.recvFrom[msg.Content] = make(map[PeerID]struct{})
				}
				n.recvFrom[msg.Content][msg.From] = struct{}{}

				if _, ok := n.seenAt[msg.Content]; !ok {
					n.seenAt[msg.Content] = time.Now()
					first = true
				}
				n.mu.Unlock()

				if first {
					go func(msg Message) {
						time.Sleep(time.Duration(n.nodeLatency) * time.Millisecond)
						n.publish(network, msg)
					}(msg)
				}
			}
		}
	}(ctx, wg)
}

// // copyIDSet creates a shallow copy of a set of IDs.
// func copyIDSet(src map[PeerID]struct{}) map[PeerID]struct{} {
// 	dst := make(map[PeerID]struct{}, len(src))
// 	for k := range src {
// 		dst[k] = struct{}{}
// 	}
// 	return dst
// }

// publish sends the message to neighbors, excluding 'exclude' and already-sent targets.
func (n *p2pNode) publish(network *Network, msg Message) {
	content := msg.Content
	protocol := msg.Protocol
	hopCount := msg.HopCount

	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.sentTo[content]; !ok {
		n.sentTo[content] = make(map[PeerID]struct{})
	}
	if _, ok := n.recvFrom[content]; !ok {
		n.recvFrom[content] = make(map[PeerID]struct{})
	}

	willSendEdges := make([]p2pEdge, 0)

	if protocol == Flooding || protocol == Gossiping {
		for _, edge := range n.edges {
			if _, already := n.sentTo[content][edge.targetID]; already {
				continue
			}
			if _, received := n.recvFrom[content][edge.targetID]; received {
				continue
			}
			n.sentTo[content][edge.targetID] = struct{}{}

			willSendEdges = append(willSendEdges, edge)
		}

		if protocol == Gossiping && len(willSendEdges) > 0 {
			rand.Shuffle(len(willSendEdges), func(i, j int) {
				willSendEdges[i], willSendEdges[j] = willSendEdges[j], willSendEdges[i]
			})

			k := int(float64(len(willSendEdges)) * network.cfg.GossipFactor)
			willSendEdges = willSendEdges[:k]
		}
	} else if protocol == Custom {

	} else {
		return
	}

	for _, edge := range willSendEdges {
		edgeCopy := edge

		go func(e p2pEdge) {
			time.Sleep(time.Duration(e.edgeLatency) * time.Millisecond)

			network.nodes[e.targetID].msgQueue <- Message{
				From:     n.id,
				Content:  content,
				Protocol: protocol,
				HopCount: hopCount + 1,
			}
		}(edgeCopy)
	}
}
