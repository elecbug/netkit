package p2p

import (
	"math/rand/v2"
	"slices"
)

// Message represents a message sent between nodes in the P2P network.
type Message struct {
	Publisher PeerID         // ID of the peer that originally published the message
	From      PeerID         // ID of the peer that sent the message to the current peer
	Content   string         // the actual content of the message
	HopCount  int            // the number of hops the message has taken from the publisher to the current peer
	Protocol  ProtocolFunc   // the protocol function that determines how the message should be processed and forwarded
	Params    map[string]any // additional parameters for the protocol function
}

// ProtocolFunc defines the function signature for custom protocols in the P2P network.
// It takes the current peer's ID, the message being processed, a list of neighbor IDs,
// lists of peers the message has been sent to and received from, and any additional broadcast parameters.
// It returns a pointer to a slice of PeerIDs that the message should be forwarded to.
type ProtocolFunc func(id PeerID, msg Message, neighbors []PeerID, sentPeers []PeerID, receivedPeers []PeerID, broadcastParams map[string]any) *[]PeerID

// Flooding is a simple broadcast protocol where each peer forwards the message to all its neighbors except those it has already sent to or received from.
var Flooding ProtocolFunc = func(id PeerID, msg Message, neighbors []PeerID, sentPeers []PeerID, receivedPeers []PeerID, broadcastParams map[string]any) *[]PeerID {
	targets := make([]PeerID, 0)
	for _, neighbor := range neighbors {
		if slices.Contains(sentPeers, neighbor) {
			continue
		}
		if slices.Contains(receivedPeers, neighbor) {
			continue
		}

		targets = append(targets, neighbor)
	}

	return &targets
}

// Gossip is a broadcast protocol where each peer forwards the message to a random subset of its neighbors, determined by the gossip factor.
var Gossip ProtocolFunc = func(id PeerID, msg Message, neighbors []PeerID, sentPeers []PeerID, receivedPeers []PeerID, broadcastParams map[string]any) *[]PeerID {
	targets := make([]PeerID, 0)
	for _, neighbor := range neighbors {
		if slices.Contains(sentPeers, neighbor) {
			continue
		}
		if slices.Contains(receivedPeers, neighbor) {
			continue
		}

		targets = append(targets, neighbor)
	}

	gossipFactor, ok := broadcastParams["gossip_factor"].(float64)
	if !ok {
		gossipFactor = 0.5 // default gossip factor
	}

	if len(targets) > 0 {
		rand.Shuffle(len(targets), func(i, j int) {
			targets[i], targets[j] = targets[j], targets[i]
		})

		k := int(float64(len(targets)) * gossipFactor)
		targets = targets[:k]
	}

	return &targets
}
