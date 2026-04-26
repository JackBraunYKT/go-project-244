package formatters

import (
	"code/differ"
	"fmt"
	"strings"
)

func FormatPlain(diff []differ.DiffNode, path string) string {
	var lines []string
	for _, node := range diff {
		fullPath := path + "." + node.Key
		switch node.Type {
		case differ.NodeAdded:
			lines = append(lines, fmt.Sprintf(
				"Property '%s' was added with value: %s",
				fullPath, formatPlainValue(node.NewValue),
			))
		case differ.NodeRemoved:
			lines = append(lines, fmt.Sprintf(
				"Property '%s' was removed", fullPath,
			))
		case differ.NodeChanged:
			lines = append(lines, fmt.Sprintf(
				"Property '%s' was updated. From %s to %s",
				fullPath, formatPlainValue(node.OldValue), formatPlainValue(node.NewValue),
			))
		case differ.NodeNested:
			lines = append(lines, FormatPlain(node.Children, node.Key))
		}
	}
	return strings.Join(lines, "\n")
}

func formatPlainValue(value any) string {
	switch v := value.(type) {
	case map[string]any:
		return "[complex value]"
	case []any:
		return "[complex value]"
	case string:
		return fmt.Sprintf("'%s'", v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}
