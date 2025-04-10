package graph

import (
	"github.com/elecbug/go-dspkg/graph/graph/internal/graph_err" // Custom error package
)

// `graphNodes` represents a collection of nodes in a graph.
// It maintains two mappings:
//  1. `nodes`: Maps a node's unique identifier to its corresponding Node object.
//  2. `nameMap`: Maps a node's name to a list of identifiers for nodes with that name.
type graphNodes struct {
	nodes   map[NodeID]*Node    // Maps node identifiers to Node instances.
	nameMap map[string][]NodeID // Maps node names to lists of identifiers for nodes with the same name.
}

// `newNodes` creates and initializes a new graphNodes instance.
//
// Parameters:
//   - cap: Initial capacity for the internal maps.
//
// Returns a pointer to the newly created graphNodes instance.
func newNodes(cap int) *graphNodes {
	return &graphNodes{
		nodes:   make(map[NodeID]*Node, cap),
		nameMap: make(map[string][]NodeID, cap),
	}
}

// `insert` adds a new Node to the graphNodes collection.
//
// Parameters:
//   - node: The Node to be inserted.
//
// Returns an error if a node with the same identifier already exists.
func (ns *graphNodes) insert(node *Node) error {
	if _, exists := ns.nodes[node.ID()]; exists {
		// Return an error if the node identifier already exists in the collection.
		return graph_err.AlreadyNode(node.ID().String())
	} else {
		// Add the node to the nodes map.
		ns.nodes[node.ID()] = node

		// Initialize the nameMap entry if it doesn't exist.
		if ns.nameMap[node.Name] == nil {
			ns.nameMap[node.Name] = make([]NodeID, 0)
		}

		// Add the node's identifier to the nameMap.
		ns.nameMap[node.Name] = append(ns.nameMap[node.Name], node.ID())

		return nil
	}
}

// `remove` deletes a Node from the graphNodes collection using its identifier.
//
// Parameters:
//   - identifier: The unique identifier of the Node to remove.
//
// Returns an error if the Node does not exist in the collection.
func (ns *graphNodes) remove(identifier NodeID) error {
	if _, exists := ns.nodes[identifier]; exists {
		// Retrieve the node's name for nameMap cleanup.
		name := ns.nodes[identifier].Name

		// Remove the node from the nodes map.
		delete(ns.nodes, identifier)

		// Remove the node's identifier from the nameMap.
		for i := 0; i < len(ns.nameMap[name]); i++ {
			if ns.nameMap[name][i] == identifier {
				ns.nameMap[name] = append(ns.nameMap[name][:i], ns.nameMap[name][i+1:]...)
				break
			}
		}

		return nil
	} else {
		// Return an error if the node identifier does not exist.
		return graph_err.NotExistNode(identifier.String())
	}
}

// `find` retrieves a Node by its identifier.
//
// Parameters:
//   - identifier: The unique identifier of the Node.
//
// Returns the Node instance if found, or nil if the Node does not exist.
func (ns *graphNodes) find(identifier NodeID) *Node {
	return ns.nodes[identifier]
}

// `findAll` retrieves all Nodes with a given name.
//
// Parameters:
//   - name: The name of the Nodes to find.
//
// Returns a slice of Node pointers matching the given name.
func (ns *graphNodes) findAll(name string) []*Node {
	ids := ns.nameMap[name] // Get all identifiers for the given name.
	var result = make([]*Node, len(ids))

	// Populate the result slice with Node instances.
	for i, id := range ids {
		result[i] = ns.nodes[id]
	}

	return result
}
