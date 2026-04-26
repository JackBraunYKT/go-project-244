package formatters

import (
	"code/internal/differ"
	"strings"
	"testing"
)

func TestFormatPlain_Added(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "verbose", Type: differ.NodeAdded, NewValue: true},
	}

	got := FormatPlain(nodes, "")
	want := "Property 'verbose' was added with value: true"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatPlain_Removed(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "proxy", Type: differ.NodeRemoved, OldValue: "123.234.53.22"},
	}

	got := FormatPlain(nodes, "")
	want := "Property 'proxy' was removed"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatPlain_Changed(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "timeout", Type: differ.NodeChanged, OldValue: float64(20), NewValue: float64(50)},
	}

	got := FormatPlain(nodes, "")
	want := "Property 'timeout' was updated. From 20 to 50"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatPlain_AddedStringValue(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "setting4", Type: differ.NodeAdded, NewValue: "blah blah"},
	}

	got := FormatPlain(nodes, "")
	want := "Property 'setting4' was added with value: 'blah blah'"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatPlain_AddedComplexValue(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "setting5", Type: differ.NodeAdded, NewValue: map[string]any{"key": "val"}},
	}

	got := FormatPlain(nodes, "")
	want := "Property 'setting5' was added with value: [complex value]"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatPlain_AddedNullValue(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "setting3", Type: differ.NodeChanged, OldValue: true, NewValue: nil},
	}

	got := FormatPlain(nodes, "")
	want := "Property 'setting3' was updated. From true to null"

	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFormatPlain_UnchangedSkipped(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeUnchanged, OldValue: "hexlet.io"},
	}

	got := FormatPlain(nodes, "")

	if got != "" {
		t.Errorf("expected empty string for unchanged node, got: %s", got)
	}
}

// тест выявляет баг: при вложенности путь должен накапливаться
func TestFormatPlain_Nested(t *testing.T) {
	nodes := []differ.DiffNode{
		{
			Key:  "common",
			Type: differ.NodeNested,
			Children: []differ.DiffNode{
				{Key: "follow", Type: differ.NodeAdded, NewValue: false},
			},
		},
	}

	got := FormatPlain(nodes, "")

	if !strings.Contains(got, "Property 'common.follow' was added") {
		t.Errorf("expected 'common.follow' in output, got:\n%s", got)
	}
}

// тест выявляет баг: путь не должен начинаться с точки
func TestFormatPlain_NoLeadingDot(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	got := FormatPlain(nodes, "")

	if strings.HasPrefix(got, "Property '.") {
		t.Errorf("path must not start with dot, got:\n%s", got)
	}
}

func TestFormatPlain_DeepNested(t *testing.T) {
	nodes := []differ.DiffNode{
		{
			Key:  "common",
			Type: differ.NodeNested,
			Children: []differ.DiffNode{
				{
					Key:  "setting6",
					Type: differ.NodeNested,
					Children: []differ.DiffNode{
						{Key: "ops", Type: differ.NodeAdded, NewValue: "vops"},
					},
				},
			},
		},
	}

	got := FormatPlain(nodes, "")

	if !strings.Contains(got, "Property 'common.setting6.ops' was added") {
		t.Errorf("expected 'common.setting6.ops' in output, got:\n%s", got)
	}
}
