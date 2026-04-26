package code

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
)

type Format string

const (
	Stylish Format = "stylish"
)

func GenDiff(filepath1, filepath2 string, format Format) (*string, error) {
	data1, err := readFile(filepath1)
	if err != nil {
		return nil, err
	}

	data2, err := readFile(filepath2)
	if err != nil {
		return nil, err
	}

	ext1 := filepath.Ext(filepath1)
	ext2 := filepath.Ext(filepath2)

	if !strings.EqualFold(ext1, ext2) {
		return nil, fmt.Errorf("files have different extensions")
	}

	parsed1, err := parse(data1, ext1)
	if err != nil {
		return nil, err
	}

	parsed2, err := parse(data2, ext2)
	if err != nil {
		return nil, err
	}

	nodes := buildDiffNodes(parsed1, parsed2)
	formattedNodes, err := formatNodes(nodes, format, 1)
	if err != nil {
		return nil, err
	}

	return formattedNodes, nil
}

func readFile(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path %s is a directory", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return data, nil
}

func parse(data []byte, ext string) (map[string]interface{}, error) {
	switch ext {
	case ".json":
		return parseJSON(data)
	default:
		return nil, fmt.Errorf("unsupported ext: %s", ext)
	}
}

func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}

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

func buildDiffNodes(data1, data2 map[string]interface{}) []DiffNode {
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
				node.Children = buildDiffNodes(map1, map2)
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

func formatNodes(nodes []DiffNode, format Format, depth int) (*string, error) {
	switch format {
	case Stylish:
		result := formatStylish(nodes, depth)
		return &result, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func formatStylish(nodes []DiffNode, depth int) string {
	indent := strings.Repeat("    ", depth-1)
	signIndent := indent + "  "
	lines := []string{"{"}

	for _, node := range nodes {
		switch node.Type {
		case NodeAdded:
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", signIndent, node.Key, stringify(node.NewValue, depth)))
		case NodeRemoved:
			lines = append(lines, fmt.Sprintf("%s- %s: %s", signIndent, node.Key, stringify(node.OldValue, depth)))
		case NodeUnchanged:
			lines = append(lines, fmt.Sprintf("%s  %s: %s", signIndent, node.Key, stringify(node.OldValue, depth)))
		case NodeChanged:
			lines = append(lines, fmt.Sprintf("%s- %s: %s", signIndent, node.Key, stringify(node.OldValue, depth)))
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", signIndent, node.Key, stringify(node.NewValue, depth)))
		case NodeNested:
			lines = append(lines, fmt.Sprintf("%s  %s: %s", signIndent, node.Key, formatStylish(node.Children, depth+1)))
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
