package formatters

import (
	"code/internal/differ"
	"fmt"
)

const (
	Stylish = "stylish"
	Plain   = "plain"
	JSON    = "json"
)

var SupportedFormats = []string{
	Stylish,
	Plain,
	JSON,
}

func FormatNodes(nodes []differ.DiffNode, format string, depth int) (*string, error) {
	var result string

	switch format {
	case Stylish:
		result = FormatStylish(nodes, depth)
	case Plain:
		result = FormatPlain(nodes, "")
	case JSON:
		formattedNodes, err := FormatJSON(nodes)
		if err != nil {
			return nil, err
		}
		result = formattedNodes
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return &result, nil
}
