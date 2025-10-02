package p2p

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// ID represents a unique identifier for a node in the P2P network.
type ID uint64

// Edge represents a connection from one node to another in the P2P network.
type Edge struct {
	TargetID ID
	Latency  float64 // in milliseconds
}

// Node represents a node in the P2P network.
type Node struct {
	ID      ID
	Latency float64
	Edges   map[ID]Edge

	RecvFrom map[string]map[ID]struct{} // content -> set of senders
	SentTo   map[string]map[ID]struct{} // content -> set of targets
	SeenAt   map[string]time.Time       // content -> first arrival time

	msgQueue chan Message
	mu       sync.Mutex
}

// Degree returns the number of edges connected to the node.
func (n *Node) Degree() int {
	return len(n.Edges)
}

// Message represents a message sent between nodes in the P2P network.
type Message struct {
	From    ID
	Content string
}

// eachRun starts the message handling routine for the node.
func (n *Node) eachRun(network map[ID]*Node) {
	go func() {
		n.msgQueue = make(chan Message, 1000)
		n.RecvFrom = make(map[string]map[ID]struct{})
		n.SentTo = make(map[string]map[ID]struct{})
		n.SeenAt = make(map[string]time.Time)

		for msg := range n.msgQueue {
			first := false
			var excludeSnapshot map[ID]struct{}

			n.mu.Lock()
			if _, ok := n.RecvFrom[msg.Content]; !ok {
				n.RecvFrom[msg.Content] = make(map[ID]struct{})
			}
			n.RecvFrom[msg.Content][msg.From] = struct{}{}

			if _, ok := n.SeenAt[msg.Content]; !ok {
				n.SeenAt[msg.Content] = time.Now()
				first = true
				excludeSnapshot = copyIDSet(n.RecvFrom[msg.Content])
			}
			n.mu.Unlock()

			if first {
				go func(content string, exclude map[ID]struct{}) {
					time.Sleep(time.Duration(n.Latency) * time.Millisecond)
					n.publish(network, content, exclude)
				}(msg.Content, excludeSnapshot)
			}
		}
	}()
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
func (n *Node) publish(network map[ID]*Node, content string, exclude map[ID]struct{}) {
	n.mu.Lock()
	if _, ok := n.SentTo[content]; !ok {
		n.SentTo[content] = make(map[ID]struct{})
	}

	for _, edge := range n.Edges {
		if _, wasSender := exclude[edge.TargetID]; wasSender {
			continue
		}
		if _, already := n.SentTo[content][edge.TargetID]; already {
			continue
		}
		if _, received := n.RecvFrom[content][edge.TargetID]; received {
			continue
		}
		n.SentTo[content][edge.TargetID] = struct{}{}

		edgeCopy := edge
		go func(e Edge) {
			time.Sleep(time.Duration(e.Latency) * time.Millisecond)
			network[e.TargetID].msgQueue <- Message{From: n.ID, Content: content}
		}(edgeCopy)
	}
	n.mu.Unlock()
}

// LogNormalRand generates a log-normally distributed random number
// with given mu and sigma parameters.
func LogNormalRand(mu, sigma float64, src rand.Source) float64 {
	r := rand.New(src)

	u1 := r.Float64()
	u2 := r.Float64()
	z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	return math.Exp(mu + sigma*z)
}
