// Package broadcast defines the protocols used for communication in the P2P network.
package broadcast

// Protocol defines the protocol used for broadcasting messages in the P2P network.
type Protocol int

var (
	Flooding  Protocol = 0
	Gossiping Protocol = 1
)
