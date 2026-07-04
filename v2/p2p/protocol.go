package p2p

import (
	"math/rand/v2"
	"slices"
)

// Message represents a message sent between nodes in the P2P network.
type Message struct {
	Publisher     PeerID         // ID of the peer that originally published the message
	From          PeerID         // ID of the peer that sent the message to the current peer
	Content       string         // the actual content of the message
	HopCount      int            // the number of hops the message has taken from the publisher to the current peer
	Protocol      ProtocolFunc   // the protocol function that determines how the message should be processed and forwarded
	StaticParams  map[string]any // additional parameters for the protocol function
	DynamicParams map[string]any // additional parameters that can change during message processing
}

// ProtocolFunc defines the function signature for custom protocols in the P2P network.
// It takes the current peer's ID, the message being processed, a list of neighbor IDs,
// lists of peers the message has been sent to and received from, and any additional static parameters.
// It returns a pointer to a slice of PeerIDs that the message should be forwarded to, along with any dynamic parameters.
type ProtocolFunc func(id PeerID, msg Message, neighbors []PeerID, sentPeers []PeerID, receivedPeers []PeerID, staticParams, dynamicParams map[string]any) (*[]PeerID, map[PeerID]map[string]any)

// Flooding is a simple broadcast protocol where each peer forwards the message to all its neighbors except those it has already sent to or received from.
var Flooding ProtocolFunc = func(id PeerID, msg Message, neighbors []PeerID, sentPeers []PeerID, receivedPeers []PeerID, staticParams, dynamicParams map[string]any) (*[]PeerID, map[PeerID]map[string]any) {
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

	return &targets, nil
}

// Gossip is a broadcast protocol where each peer forwards the message to a random subset of its neighbors, determined by the gossip factor.
var Gossip ProtocolFunc = func(id PeerID, msg Message, neighbors []PeerID, sentPeers []PeerID, receivedPeers []PeerID, staticParams, dynamicParams map[string]any) (*[]PeerID, map[PeerID]map[string]any) {
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

	gossipFactor, ok1 := staticParams["gossip_factor"].(float64)
	gossipNode, ok2 := staticParams["gossip_node"].(int)

	if !ok1 && !ok2 {
		ok1 = true
		gossipFactor = 0.5
	}

	if len(targets) > 0 {
		rand.Shuffle(len(targets), func(i, j int) {
			targets[i], targets[j] = targets[j], targets[i]
		})

		if ok1 {
			k := int(float64(len(targets)) * gossipFactor)
			targets = targets[:k]
		} else if ok2 {
			k := gossipNode
			if k > len(targets) {
				k = len(targets)
			}
			
			targets = targets[:k]
		}
	}

	return &targets, nil
}
