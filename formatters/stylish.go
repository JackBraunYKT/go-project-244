package formatters

import (
	"code/differ"
	"fmt"
	"sort"
	"strings"
)

func FormatStylish(nodes []differ.DiffNode, depth int) string {
	indent := strings.Repeat("    ", depth-1)
	signIndent := indent + "  "
	lines := []string{"{"}

	for _, node := range nodes {
		switch node.Type {
		case differ.NodeAdded:
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", signIndent, node.Key, stringify(node.NewValue, depth)))
		case differ.NodeRemoved:
			lines = append(lines, fmt.Sprintf("%s- %s: %s", signIndent, node.Key, stringify(node.OldValue, depth)))
		case differ.NodeUnchanged:
			lines = append(lines, fmt.Sprintf("%s  %s: %s", signIndent, node.Key, stringify(node.OldValue, depth)))
		case differ.NodeChanged:
			lines = append(lines, fmt.Sprintf("%s- %s: %s", signIndent, node.Key, stringify(node.OldValue, depth)))
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", signIndent, node.Key, stringify(node.NewValue, depth)))
		case differ.NodeNested:
			lines = append(lines, fmt.Sprintf("%s  %s: %s", signIndent, node.Key, FormatStylish(node.Children, depth+1)))
		}
	}

	lines = append(lines, indent+"}")
	return strings.Join(lines, "\n")
}

func stringify(value interface{}, depth int) string {
	switch v := value.(type) {
	case map[string]interface{}:
		keys := getSortedKeys(v)
		indent := strings.Repeat("    ", depth)
		lines := []string{"{"}

		for _, key := range keys {
			lines = append(lines, fmt.Sprintf("%s    %s: %s", indent, key, stringify(v[key], depth+1)))
		}

		lines = append(lines, indent+"}")
		return strings.Join(lines, "\n")
	case nil:
		return "null"
	case bool:
		if v {
			return "true"
		}
		return "false"
	case string:
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

func getSortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
