package p2p

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/elecbug/netkit/v2/graph"
)

// PeerID is a type alias for graph.NodeID, representing the unique identifier for a peer in the P2P network.
type PeerID graph.NodeID

// P2P represents a peer-to-peer network, containing a map of peers and the configuration for the network.
type P2P struct {
	peers map[PeerID]*peer
	cfg   *Config
}

// Config holds configuration parameters for the P2P network, including functions to generate processing and network latencies.
type Config struct {
	// ProcessingLatencyFunc generates the latency for processing a message at the source peer. It should return the latency in milliseconds.
	ProcessingLatencyFunc func(src PeerID) float64
	// NetworkLatencyFunc generates the latency for a message sent from src to dst. It should return the latency in milliseconds.
	NetworkLatencyFunc func(src PeerID, dst PeerID) float64
}

// New creates a new P2P network from the given graph. It returns an error if the graph is weighted,
// as weighted graphs are not supported for P2P generation.
func New(source *graph.Graph, cfg *Config) (*P2P, error) {
	if source.IsWeighted() {
		return nil, fmt.Errorf("weighted graphs are not supported for P2P generation")
	}

	nodes := make(map[PeerID]*peer)
	maps := make(map[graph.NodeID]PeerID)

	// create nodes
	for _, gn := range source.Nodes() {
		n := newPeer(PeerID(gn), cfg.ProcessingLatencyFunc(PeerID(gn)))
		n.edges = make(map[PeerID]edge)

		nodes[n.id] = n
		maps[gn] = n.id
	}

	for _, gn := range source.Nodes() {
		n := nodes[PeerID(gn)]

		node, err := source.Node(gn)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s from graph: %v", gn, err)
		}

		for _, neighbor := range node.Neighbors() {
			j := maps[neighbor]

			edge := edge{
				targetID:       PeerID(j),
				networkLatency: cfg.NetworkLatencyFunc(PeerID(gn), PeerID(j)),
			}

			n.edges[edge.targetID] = edge
		}
	}

	return &P2P{peers: nodes, cfg: cfg}, nil
}

// Free clears all peers from the P2P network, effectively resetting it to an empty state.
func (p *P2P) Free() {
	for id := range p.peers {
		p.peers[id].eachStop()
		delete(p.peers, id)
	}
}

/* Basic Actions */

// Run starts the message handling routines for all peers in the network.
func (p *P2P) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(len(p.peers))

	for _, peer := range p.peers {
		peer.eachRun(p, wg, ctx)
	}

	wg.Wait()
}

// ExpireSimulation runs the simulation until the reachability of the specified message stabilizes or a timeout occurs.
func (p *P2P) ExpireSimulation(cancel context.CancelFunc, msg string, expirationDuration, timeoutDuration, checkInterval time.Duration) {
	startTime := time.Now()
	lastChangeTime := startTime
	beforeRch := p.Reachability(msg)

	for {
		currentRch := p.Reachability(msg)

		if currentRch > beforeRch {
			beforeRch = currentRch
			lastChangeTime = time.Now()
		}

		if time.Since(lastChangeTime) > expirationDuration {
			break
		}
		if time.Since(startTime) > timeoutDuration {
			break
		}

		time.Sleep(checkInterval)
	}

	for _, peer := range p.peers {
		peer.eachStop()
	}

	cancel()
}

/* Network Information */

// PeerIDs returns a slice of all node IDs in the network.
func (p *P2P) PeerIDs() []PeerID {
	ids := make([]PeerID, 0, len(p.peers))

	for id := range p.peers {
		ids = append(ids, id)
	}

	slices.Sort(ids)

	return ids
}

/* Message Handling */

// Publish sends a message to the specified peer's message queue.
func (p *P2P) Publish(id PeerID, msg string, protocol ProtocolFunc, params map[string]any) error {
	if peer, ok := p.peers[id]; ok {
		if !peer.alive {
			return fmt.Errorf("peer %s is not alive", id)
		}

		peer.msgQueue <- Message{
			Publisher: id,
			From:      id,
			Content:   msg,
			Protocol:  protocol,
			Params:    params,
			HopCount:  0,
		}
		return nil
	}

	return fmt.Errorf("peer %s not found", id)
}

// Reachability calculates the fraction of peers that have received the specified message.
func (p *P2P) Reachability(msg string) float64 {
	total := 0
	reached := 0

	for _, peer := range p.peers {
		total++
		peer.mu.Lock()
		if _, ok := peer.seenAt[msg]; ok {
			reached++
		}
		peer.mu.Unlock()
	}

	return float64(reached) / float64(total)
}

// FirstMessageReceptionTimes returns the first reception times of the specified message across all peers.
func (p *P2P) FirstMessageReceptionTimes(msg string) []time.Time {
	firstTimes := make([]time.Time, 0)

	for _, peer := range p.peers {
		peer.mu.Lock()
		if t, ok := peer.seenAt[msg]; ok {
			firstTimes = append(firstTimes, t)
		}

		peer.mu.Unlock()
	}

	return firstTimes
}

// FirstMessageReceptions returns the first reception details of the specified message across all peers, including the peer ID, the sender's peer ID, and the timestamp.
func (p *P2P) FirstMessageReceptions(msg string) []struct {
	PeerID    PeerID    `json:"peer_id"`
	From      PeerID    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
} {
	receptions := make([]struct {
		PeerID    PeerID    `json:"peer_id"`
		From      PeerID    `json:"from"`
		Timestamp time.Time `json:"timestamp"`
	}, 0)

	for _, peer := range p.peers {
		peer.mu.Lock()
		if t, ok := peer.seenAt[msg]; ok {
			from := peer.firstFrom[msg]

			receptions = append(receptions, struct {
				PeerID    PeerID    `json:"peer_id"`
				From      PeerID    `json:"from"`
				Timestamp time.Time `json:"timestamp"`
			}{
				PeerID:    peer.id,
				From:      from,
				Timestamp: t,
			})
		}
		peer.mu.Unlock()
	}

	return receptions
}

// DuplicateMessageCount counts how many duplicate messages were received across all peers.
func (p *P2P) DuplicateMessageCount(msg string) int {
	dupCount := 0

	for _, peer := range p.peers {
		peer.mu.Lock()
		if count, ok := peer.recvFrom[msg]; ok {
			dupCount += len(count) - 1
		}
		peer.mu.Unlock()
	}

	return dupCount
}

// MessageInfo returns a snapshot of the peer's message-related information.
func (p *P2P) MessageInfo(peerID PeerID, content string) (map[string]any, error) {
	peer := p.peers[peerID]

	if peer == nil {
		return nil, fmt.Errorf("peer %s not found", peerID)
	}

	peer.mu.Lock()
	defer peer.mu.Unlock()

	info := make(map[string]any)

	info["recv"] = make([]PeerID, 0)
	for k := range peer.recvFrom[content] {
		info["recv"] = append(info["recv"].([]PeerID), k)
	}

	info["sent"] = make([]PeerID, 0)
	for k := range peer.sentTo[content] {
		info["sent"] = append(info["sent"].([]PeerID), k)
	}

	info["seen"] = peer.seenAt[content].String()
	info["first_from"] = peer.firstFrom[content]

	return info, nil
}

// PeerLog returns a copy of the log entries for the specified peer, allowing for inspection of message flow and events.
func (p *P2P) PeerLog(peerID PeerID, content string) (map[string][]logEntry, error) {
	peer := p.peers[peerID]

	if peer == nil {
		return nil, fmt.Errorf("peer %s not found", peerID)
	}

	peer.mu.Lock()
	defer peer.mu.Unlock()

	logCopy := make(map[string][]logEntry)
	for k, v := range peer.log {
		logCopy[k] = append([]logEntry(nil), v...)
	}

	return logCopy, nil
}
