package p2p

import (
	"context"
	"sync"
	"time"
)

// peer represents a node in the P2P network.
type peer struct {
	id                PeerID          // unique identifier for the peer
	processingLatency float64         // latency for processing a message at the source peer, in milliseconds
	edges             map[PeerID]edge // connections to other peers, mapping target peer ID to edge information

	recvFrom  map[string]map[PeerID]struct{} // content -> set of senders
	sentTo    map[string]map[PeerID]struct{} // content -> set of targets
	seenAt    map[string]time.Time           // content -> first arrival time
	firstFrom map[string]PeerID              // content -> first sender

	msgQueue chan Message // channel for incoming messages
	mu       sync.Mutex   // mutex to protect access to the peer's state

	alive bool // indicates whether the peer is active in the network
}

// edge represents a connection from one node to another in the P2P network.
type edge struct {
	targetID       PeerID  // ID of the target peer
	networkLatency float64 // latency for a message sent from this peer to the target peer, in milliseconds
}

// newPeer creates a new Node with the given ID and node latency.
func newPeer(id PeerID, nodeLatency float64) *peer {
	return &peer{
		id:                id,
		processingLatency: nodeLatency,
		edges:             make(map[PeerID]edge),

		recvFrom:  make(map[string]map[PeerID]struct{}),
		sentTo:    make(map[string]map[PeerID]struct{}),
		seenAt:    make(map[string]time.Time),
		firstFrom: make(map[string]PeerID),

		msgQueue: make(chan Message, 1000),
		mu:       sync.Mutex{},
	}
}

// eachRun starts the message handling routine for the peer.
func (p *peer) eachRun(network *P2P, wg *sync.WaitGroup, ctx context.Context) {
	go func(ctx context.Context, wg *sync.WaitGroup) {
		p.alive = true
		wg.Done()

		select {
		case <-ctx.Done():
			p.alive = false
			return
		default:
			for msg := range p.msgQueue {
				first := false

				p.mu.Lock()
				if _, ok := p.recvFrom[msg.Content]; !ok {
					p.recvFrom[msg.Content] = make(map[PeerID]struct{})
				}
				p.recvFrom[msg.Content][msg.From] = struct{}{}

				if _, ok := p.seenAt[msg.Content]; !ok {
					p.seenAt[msg.Content] = time.Now()
					p.firstFrom[msg.Content] = msg.From
					first = true
				}
				p.mu.Unlock()

				if first {
					go func(msg Message) {
						time.Sleep(time.Duration(p.processingLatency) * time.Millisecond)
						p.eachPublish(network, msg)
					}(msg)
				}
			}
		}
	}(ctx, wg)
}

// eachPublish sends the message to neighbors, excluding 'exclude' and already-sent targets.
func (p *peer) eachPublish(network *P2P, msg Message) {
	content := msg.Content
	protocol := msg.Protocol
	hopCount := msg.HopCount

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.sentTo[content]; !ok {
		p.sentTo[content] = make(map[PeerID]struct{})
	}
	if _, ok := p.recvFrom[content]; !ok {
		p.recvFrom[content] = make(map[PeerID]struct{})
	}

	willSendEdges := make([]edge, 0)

	allEdges := make([]PeerID, 0)
	for _, edge := range p.edges {
		allEdges = append(allEdges, edge.targetID)
	}

	sentEdges := make([]PeerID, 0)
	for targetID := range p.sentTo[content] {
		sentEdges = append(sentEdges, targetID)
	}

	receivedEdges := make([]PeerID, 0)
	for senderID := range p.recvFrom[content] {
		receivedEdges = append(receivedEdges, senderID)
	}

	targets := msg.Protocol(p.id, msg, allEdges, sentEdges, receivedEdges, msg.Params)

	for _, targetID := range *targets {
		for _, edge := range p.edges {
			if edge.targetID == targetID {
				willSendEdges = append(willSendEdges, edge)
				break
			}
		}
	}

	for _, e := range willSendEdges {
		edgeCopy := e
		p.sentTo[content][e.targetID] = struct{}{}

		go func(e edge) {
			time.Sleep(time.Duration(e.networkLatency) * time.Millisecond)

			network.peers[e.targetID].msgQueue <- Message{
				Publisher: msg.Publisher,
				From:      p.id,
				Content:   content,
				Protocol:  protocol,
				HopCount:  hopCount + 1,
			}
		}(edgeCopy)
	}
}
