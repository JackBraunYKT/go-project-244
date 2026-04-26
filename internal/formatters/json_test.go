package formatters

import (
	"code/internal/differ"
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatJSON_ValidJSON(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// проверяем что вывод является валидным JSON
	var result []interface{}
	if err := json.Unmarshal([]byte(got), &result); err != nil {
		t.Errorf("output is not valid JSON: %v\ngot:\n%s", err, got)
	}
}

func TestFormatJSON_ContainsKey(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, `"Key"`) && !strings.Contains(got, `"key"`) {
		t.Errorf("expected key field in output:\n%s", got)
	}
	if !strings.Contains(got, "hexlet.io") {
		t.Errorf("expected value 'hexlet.io' in output:\n%s", got)
	}
}

func TestFormatJSON_EmptyDiff(t *testing.T) {
	got, err := FormatJSON([]differ.DiffNode{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// пустой слайс должен сериализоваться в []
	if strings.TrimSpace(got) != "[]" {
		t.Errorf("expected '[]' for empty diff, got: %s", got)
	}
}

func TestFormatJSON_NodeTypes(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "added", Type: differ.NodeAdded, NewValue: true},
		{Key: "removed", Type: differ.NodeRemoved, OldValue: "val"},
		{Key: "unchanged", Type: differ.NodeUnchanged, OldValue: 42},
		{Key: "changed", Type: differ.NodeChanged, OldValue: "old", NewValue: "new"},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, nodeType := range []string{
		differ.NodeAdded,
		differ.NodeRemoved,
		differ.NodeUnchanged,
		differ.NodeChanged,
	} {
		if !strings.Contains(got, nodeType) {
			t.Errorf("expected node type %q in output:\n%s", nodeType, got)
		}
	}
}

func TestFormatJSON_Indented(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// MarshalIndent должен добавлять отступы
	if !strings.Contains(got, "\n") || !strings.Contains(got, "  ") {
		t.Errorf("expected indented JSON output:\n%s", got)
	}
}

func TestFormatJSON_NestedChildren(t *testing.T) {
	nodes := []differ.DiffNode{
		{
			Key:  "common",
			Type: differ.NodeNested,
			Children: []differ.DiffNode{
				{Key: "follow", Type: differ.NodeAdded, NewValue: false},
			},
		},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// проверяем что вложенность сериализовалась
	if !strings.Contains(got, "common") {
		t.Errorf("expected 'common' in output:\n%s", got)
	}
	if !strings.Contains(got, "follow") {
		t.Errorf("expected nested key 'follow' in output:\n%s", got)
	}

	// весь вывод должен быть валидным JSON
	var result []interface{}
	if err := json.Unmarshal([]byte(got), &result); err != nil {
		t.Errorf("output is not valid JSON: %v\ngot:\n%s", err, got)
	}
}

func TestFormatJSON_NullValue(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "proxy", Type: differ.NodeAdded, NewValue: nil},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "null") {
		t.Errorf("expected 'null' in output:\n%s", got)
	}
}
