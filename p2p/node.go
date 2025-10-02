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
	ID           ID
	Latency      float64 // in milliseconds
	Edges        map[ID]Edge
	Received     map[string][]ID
	ReceivedTime map[string]time.Time
	msgQueue     chan Message
	mu           sync.Mutex
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
		n.Received = make(map[string][]ID)
		n.ReceivedTime = make(map[string]time.Time)

		n.mu = sync.Mutex{}

		for {
			for msg := range n.msgQueue {
				if _, ok := n.ReceivedTime[msg.Content]; !ok {
					n.mu.Lock()
					n.ReceivedTime[msg.Content] = time.Now()
					n.Received[msg.Content] = []ID{}
					n.Received[msg.Content] = append(n.Received[msg.Content], msg.From)
					n.mu.Unlock()

					go func(msg Message) {
						time.Sleep(time.Duration(n.Latency) * time.Millisecond)
						n.mu.Lock()
						n.publish(network, msg.Content)
						n.mu.Unlock()
					}(msg)
				} else {
					n.mu.Lock()
					n.Received[msg.Content] = append(n.Received[msg.Content], msg.From)
					n.mu.Unlock()
				}
			}
		}
	}()
}

// publish sends the message to all connected nodes except those in the ids slice.
func (n *Node) publish(network map[ID]*Node, msg string) {
	for _, edge := range n.Edges {
		found := false

		for _, id := range n.Received[msg] {
			if id == edge.TargetID {
				found = true
				break
			}
		}

		if found {
			continue
		}

		go func(edge Edge) {
			time.Sleep(time.Duration(edge.Latency) * time.Millisecond)
			msg := Message{From: n.ID, Content: msg}
			network[edge.TargetID].msgQueue <- msg
		}(edge)
	}
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
