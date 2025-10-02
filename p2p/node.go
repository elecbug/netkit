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
	ID       ID
	Latency  float64 // in milliseconds
	Edges    map[ID]Edge
	msgQueue chan Message
}

// Degree returns the number of edges connected to the node.
func (n *Node) Degree() int {
	return len(n.Edges)
}

// Message represents a message sent between nodes in the P2P network.
type Message string

// eachRun starts the message handling routine for the node.
func (n *Node) eachRun(network map[ID]*Node) {
	n.msgQueue = make(chan Message, 100)
	receivedMap := make(map[Message][]ID)
	mu := sync.Mutex{}

	go func() {
		for {
			for msg := range n.msgQueue {
				if _, ok := receivedMap[msg]; !ok {
					mu.Lock()
					receivedMap[msg] = []ID{}
					mu.Unlock()

					go func() {
						time.Sleep(time.Duration(n.Latency) * time.Millisecond)
						publish(n.Edges, receivedMap[msg], network, msg)
					}()
				}
				mu.Lock()
				receivedMap[msg] = append(receivedMap[msg], n.ID)
				mu.Unlock()
			}
		}
	}()
}

// publish sends the message to all connected nodes except those in the ids slice.
func publish(edges map[ID]Edge, ids []ID, network map[ID]*Node, msg Message) {
	for i, edge := range edges {
		found := false

		for _, id := range ids {
			if id == i {
				found = true
				break
			}
		}

		if found {
			continue
		}

		time.Sleep(time.Duration(edge.Latency) * time.Millisecond)
		network[edge.TargetID].msgQueue <- msg
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
