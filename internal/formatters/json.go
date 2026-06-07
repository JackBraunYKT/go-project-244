package formatters

import (
	"code/internal/differ"
	"encoding/json"
	"fmt"
)

type jsonDiffNode struct {
	Key      string           `json:"key"`
	Type     string           `json:"type"`
	Value    *json.RawMessage `json:"value,omitempty"`
	OldValue *json.RawMessage `json:"oldValue,omitempty"`
	NewValue *json.RawMessage `json:"newValue,omitempty"`
	Children *[]jsonDiffNode  `json:"children,omitempty"`
}

func FormatJSON(diff []differ.DiffNode) (string, error) {
	nodes, err := buildJSONDiffNodes(diff)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func buildJSONDiffNodes(diff []differ.DiffNode) ([]jsonDiffNode, error) {
	nodes := make([]jsonDiffNode, 0, len(diff))

	for _, node := range diff {
		jsonNode := jsonDiffNode{
			Key:  node.Key,
			Type: node.Type,
		}

		switch node.Type {
		case differ.NodeAdded:
			value, err := marshalJSONValue(node.NewValue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal added value for key %q: %w", node.Key, err)
			}
			jsonNode.Value = value
		case differ.NodeRemoved, differ.NodeUnchanged:
			value, err := marshalJSONValue(node.OldValue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal value for key %q: %w", node.Key, err)
			}
			jsonNode.Value = value
		case differ.NodeChanged:
			oldValue, err := marshalJSONValue(node.OldValue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal old value for key %q: %w", node.Key, err)
			}
			newValue, err := marshalJSONValue(node.NewValue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal new value for key %q: %w", node.Key, err)
			}
			jsonNode.OldValue = oldValue
			jsonNode.NewValue = newValue
		case differ.NodeNested:
			children, err := buildJSONDiffNodes(node.Children)
			if err != nil {
				return nil, err
			}
			jsonNode.Children = &children
		default:
			return nil, fmt.Errorf("unknown diff node type %q for key %q", node.Type, node.Key)
		}

		nodes = append(nodes, jsonNode)
	}

	return nodes, nil
}

func marshalJSONValue(value any) (*json.RawMessage, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	raw := json.RawMessage(data)
	return &raw, nil
}
