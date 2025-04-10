package graph_err

import (
	"fmt"
)

func InvalidEdge(graphKey, edgeKey string) error {
	return fmt.Errorf("invalid edge for graph type: [(%s) does not fit %s]", edgeKey, graphKey)
}

func SelfEdge(key string) error {
	return fmt.Errorf("node can not connect to self: [%s]", key)
}

func AlreadyEdge(fromKey, toKey string) error {
	return fmt.Errorf("edge is already existed: [%s ---> %s]", fromKey, toKey)
}

func NotExistEdge(fromKey, toKey string) error {
	return fmt.Errorf("edge not exist: [%s ---> %s]", fromKey, toKey)
}

func AlreadyNode(key string) error {
	return fmt.Errorf("node is already existed: [%s]", key)
}

func NotExistNode(key string) error {
	return fmt.Errorf("node not exist: [%s]", key)
}
