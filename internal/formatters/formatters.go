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

type Formatter interface {
	Format([]differ.DiffNode) (string, error)
}

type StylishFormatter struct {
	Depth int
}

func (f StylishFormatter) Format(nodes []differ.DiffNode) (string, error) {
	depth := f.Depth
	if depth == 0 {
		depth = 1
	}

	return FormatStylish(nodes, depth), nil
}

type PlainFormatter struct{}

func (PlainFormatter) Format(nodes []differ.DiffNode) (string, error) {
	return FormatPlain(nodes, ""), nil
}

type JSONFormatter struct{}

func (JSONFormatter) Format(nodes []differ.DiffNode) (string, error) {
	return FormatJSON(nodes)
}

func NewFormatter(format string) (Formatter, error) {
	return newFormatter(format, 1)
}

func newFormatter(format string, depth int) (Formatter, error) {
	switch format {
	case Stylish:
		return StylishFormatter{Depth: depth}, nil
	case Plain:
		return PlainFormatter{}, nil
	case JSON:
		return JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func FormatNodes(nodes []differ.DiffNode, format string, depth int) (*string, error) {
	formatter, err := newFormatter(format, depth)
	if err != nil {
		return nil, err
	}

	result, err := formatter.Format(nodes)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
