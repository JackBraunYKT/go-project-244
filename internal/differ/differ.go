package differ

import (
	"maps"
	"reflect"
	"slices"
)

type DiffNode struct {
	Key      string
	Type     string
	OldValue interface{}
	NewValue interface{}
	Children []DiffNode
}

const (
	NodeAdded     = "added"
	NodeRemoved   = "removed"
	NodeUnchanged = "unchanged"
	NodeChanged   = "changed"
	NodeNested    = "nested"
)

func BuildDiffNodes(data1, data2 map[string]interface{}) []DiffNode {
	keys := getMergedSortedKeys(data1, data2)
	nodes := make([]DiffNode, 0, len(keys))

	for _, key := range keys {
		value1, ok1 := data1[key]
		value2, ok2 := data2[key]

		node := DiffNode{
			Key: key,
		}

		switch {
		case !ok1:
			node.Type = NodeAdded
			node.NewValue = value2
		case !ok2:
			node.Type = NodeRemoved
			node.OldValue = value1
		default:
			map1, isMap1 := value1.(map[string]interface{})
			map2, isMap2 := value2.(map[string]interface{})

			if isMap1 && isMap2 {
				node.Type = NodeNested
				node.Children = BuildDiffNodes(map1, map2)
			} else if reflect.DeepEqual(value1, value2) {
				node.Type = NodeUnchanged
				node.OldValue = value1
			} else {
				node.Type = NodeChanged
				node.OldValue = value1
				node.NewValue = value2
			}
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func getMergedSortedKeys(m1, m2 map[string]interface{}) []string {
	keys := slices.Collect(maps.Keys(m1))

	for key := range m2 {
		if _, ok := m1[key]; !ok {
			keys = append(keys, key)
		}
	}

	slices.Sort(keys)

	return keys
}
