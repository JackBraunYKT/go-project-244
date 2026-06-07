package formatters

import (
	"code/internal/differ"
	"encoding/json"
	"fmt"
)

type jsonDiffNode struct {
	Key      string         `json:"key"`
	Type     string         `json:"type"`
	Value1   interface{}    `json:"value1,omitempty"`
	Value2   interface{}    `json:"value2,omitempty"`
	Children []jsonDiffNode `json:"children,omitempty"`
}

func FormatJSON(diff []differ.DiffNode) (string, error) {
	root := jsonDiffNode{
		Key:      "",
		Type:     "root",
		Children: buildJSONDiffNodes(diff),
	}

	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON diff: %w", err)
	}

	return string(data), nil
}

func buildJSONDiffNodes(diff []differ.DiffNode) []jsonDiffNode {
	nodes := make([]jsonDiffNode, 0, len(diff))

	for _, node := range diff {
		jsonNode := jsonDiffNode{
			Key:  node.Key,
			Type: node.Type,
		}

		switch node.Type {
		case differ.NodeAdded:
			jsonNode.Value2 = node.NewValue
		case differ.NodeRemoved, differ.NodeUnchanged:
			if node.Type == differ.NodeRemoved {
				jsonNode.Type = "deleted"
			}
			jsonNode.Value1 = node.OldValue
		case differ.NodeChanged:
			jsonNode.Value1 = node.OldValue
			jsonNode.Value2 = node.NewValue
		case differ.NodeNested:
			jsonNode.Children = buildJSONDiffNodes(node.Children)
		}

		nodes = append(nodes, jsonNode)
	}

	return nodes
}
