package formatters

import (
	"code/differ"
	"fmt"
)

type Format string

const (
	Stylish Format = "stylish"
	Plain   Format = "plain"
	JSON    Format = "json"
)

func FormatNodes(nodes []differ.DiffNode, format Format, depth int) (*string, error) {
	var result string

	switch format {
	case Stylish:
		result = FormatStylish(nodes, depth)
	case Plain:
		result = FormatPlain(nodes, "common")
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
