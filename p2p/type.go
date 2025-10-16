package p2p

import "github.com/elecbug/netkit/p2p/broadcast"

// Message represents a message sent between nodes in the P2P network.
type Message struct {
	From     ID
	Content  string
	Protocol broadcast.Protocol
}

// Config holds configuration parameters for the P2P network.
type Config struct {
	GossipFactor float64 // fraction of neighbors to gossip to
}

// ID represents a unique identifier for a node in the P2P network.
type ID uint64
