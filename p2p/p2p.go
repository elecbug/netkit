// Package p2p provides types and interfaces for a peer-to-peer networking simulation.
package p2p

// PeerID represents a unique identifier for a node in the P2P network.
type PeerID uint64

// Message represents a message sent between nodes in the P2P network.
type Message struct {
	From           PeerID
	Content        string
	Protocol       BroadcastProtocol
	HopCount       int
	CustomProtocol func(Message, []PeerID, []PeerID, []PeerID, map[string]any) *[]PeerID
}

// Config holds configuration parameters for the P2P network.
type Config struct {
	GossipFactor float64        // fraction of neighbors to gossip to
	CustomParams map[string]any // parameters for custom protocols
}
