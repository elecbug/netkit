package p2p

// BroadcastProtocol defines the protocol used for broadcasting messages in the P2P network.
type BroadcastProtocol int

const (
	Flooding BroadcastProtocol = iota
	Gossiping
	Custom
)
