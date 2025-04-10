package graph

import (
	"fmt"
)

// NodeID represents a unique identifier assigned to nodes in a graph.
// It is defined as an unsigned integer type to ensure non-negative values.
type NodeID uint

// String converts the Identifier to its string representation.
// This is useful for displaying the node's unique identifier in a readable format.
func (id NodeID) String() string {
	// Use fmt.Sprintf to format the Identifier as a decimal string.
	return fmt.Sprintf("%d", id)
}

// const ALL_NODE = NodeID(math.MaxUint)
