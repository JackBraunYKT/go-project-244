package formatters

import (
	"code/internal/differ"
	"strings"
	"testing"
)

func TestFormatStylish_Added(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "timeout", Type: differ.NodeAdded, NewValue: float64(50)},
	}

	got := FormatStylish(nodes, 1)
	want := "{\n  + timeout: 50\n}"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatStylish_Removed(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "verbose", Type: differ.NodeRemoved, OldValue: true},
	}

	got := FormatStylish(nodes, 1)
	want := "{\n  - verbose: true\n}"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatStylish_Unchanged(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeUnchanged, OldValue: "hexlet.io"},
	}

	got := FormatStylish(nodes, 1)
	want := "{\n    host: hexlet.io\n}"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatStylish_Changed(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "timeout", Type: differ.NodeChanged, OldValue: float64(20), NewValue: float64(50)},
	}

	got := FormatStylish(nodes, 1)

	if !strings.Contains(got, "- timeout: 20") {
		t.Errorf("expected '- timeout: 20' in output:\n%s", got)
	}
	if !strings.Contains(got, "+ timeout: 50") {
		t.Errorf("expected '+ timeout: 50' in output:\n%s", got)
	}
}

func TestFormatStylish_Nested(t *testing.T) {
	nodes := []differ.DiffNode{
		{
			Key:  "group",
			Type: differ.NodeNested,
			Children: []differ.DiffNode{
				{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
			},
		},
	}

	got := FormatStylish(nodes, 1)

	if !strings.Contains(got, "group:") {
		t.Errorf("expected 'group:' in output:\n%s", got)
	}
	if !strings.Contains(got, "+ host: hexlet.io") {
		t.Errorf("expected '+ host: hexlet.io' in output:\n%s", got)
	}
}

func TestFormatStylish_NullValue(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "proxy", Type: differ.NodeAdded, NewValue: nil},
	}

	got := FormatStylish(nodes, 1)

	if !strings.Contains(got, "+ proxy: null") {
		t.Errorf("expected '+ proxy: null' in output:\n%s", got)
	}
}

func TestFormatStylish_MapValue(t *testing.T) {
	nodes := []differ.DiffNode{
		{
			Key:  "meta",
			Type: differ.NodeAdded,
			NewValue: map[string]interface{}{
				"env": "dev",
			},
		},
	}

	got := FormatStylish(nodes, 1)

	if !strings.Contains(got, "+ meta:") {
		t.Errorf("expected '+ meta:' in output:\n%s", got)
	}
	if !strings.Contains(got, "env: dev") {
		t.Errorf("expected 'env: dev' in output:\n%s", got)
	}
}

func TestFormatStylish_MultipleNodes(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeUnchanged, OldValue: "hexlet.io"},
		{Key: "timeout", Type: differ.NodeChanged, OldValue: float64(20), NewValue: float64(50)},
		{Key: "proxy", Type: differ.NodeRemoved, OldValue: "123.234.53.22"},
		{Key: "verbose", Type: differ.NodeAdded, NewValue: true},
	}

	got := FormatStylish(nodes, 1)

	cases := []string{
		"host: hexlet.io",
		"- timeout: 20",
		"+ timeout: 50",
		"- proxy: 123.234.53.22",
		"+ verbose: true",
	}

	for _, expected := range cases {
		if !strings.Contains(got, expected) {
			t.Errorf("expected %q in output:\n%s", expected, got)
		}
	}
}

func TestFormatStylish_SortedMapKeys(t *testing.T) {
	nodes := []differ.DiffNode{
		{
			Key:  "meta",
			Type: differ.NodeAdded,
			NewValue: map[string]interface{}{
				"z": "last",
				"a": "first",
			},
		},
	}

	got := FormatStylish(nodes, 1)
	aIndex := strings.Index(got, "a: first")
	zIndex := strings.Index(got, "z: last")

	if aIndex == -1 || zIndex == -1 {
		t.Fatalf("expected both keys in output:\n%s", got)
	}
	if aIndex > zIndex {
		t.Errorf("expected 'a' before 'z', but got:\n%s", got)
	}
}

func TestFormatStylish_DepthIndent(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	depth1 := FormatStylish(nodes, 1)
	depth2 := FormatStylish(nodes, 2)

	// глубина 2 должна содержать больший отступ чем глубина 1
	indent1 := strings.Index(depth1, "+")
	indent2 := strings.Index(depth2, "+")

	if indent2 <= indent1 {
		t.Errorf("expected deeper indent at depth 2 than depth 1")
	}
}
