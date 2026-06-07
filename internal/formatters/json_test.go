package formatters

import (
	"code/internal/differ"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	nodes := []differ.DiffNode{
		{
			Key:  "common",
			Type: differ.NodeNested,
			Children: []differ.DiffNode{
				{Key: "follow", Type: differ.NodeAdded, NewValue: false},
				{Key: "setting3", Type: differ.NodeChanged, OldValue: true, NewValue: nil},
			},
		},
		{Key: "host", Type: differ.NodeUnchanged, OldValue: "hexlet.io"},
		{Key: "proxy", Type: differ.NodeRemoved, OldValue: nil},
		{Key: "timeout", Type: differ.NodeChanged, OldValue: 50, NewValue: 20},
		{Key: "verbose", Type: differ.NodeAdded, NewValue: true},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `[
  {
    "key": "common",
    "type": "nested",
    "children": [
      {
        "key": "follow",
        "type": "added",
        "value": false
      },
      {
        "key": "setting3",
        "type": "changed",
        "oldValue": true,
        "newValue": null
      }
    ]
  },
  {
    "key": "host",
    "type": "unchanged",
    "value": "hexlet.io"
  },
  {
    "key": "proxy",
    "type": "removed",
    "value": null
  },
  {
    "key": "timeout",
    "type": "changed",
    "oldValue": 50,
    "newValue": 20
  },
  {
    "key": "verbose",
    "type": "added",
    "value": true
  }
]`

	if got != expected {
		t.Errorf("unexpected json output:\n%s", got)
	}
}

func TestFormatJSON_EmptyDiff(t *testing.T) {
	got, err := FormatJSON([]differ.DiffNode{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "[]" {
		t.Errorf("expected empty json array, got: %s", got)
	}
}

func TestFormatJSON_EmptyNestedChildren(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "settings", Type: differ.NodeNested, Children: []differ.DiffNode{}},
	}

	got, err := FormatJSON(nodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `[
  {
    "key": "settings",
    "type": "nested",
    "children": []
  }
]`

	if got != expected {
		t.Errorf("unexpected json output:\n%s", got)
	}
}

func TestFormatJSON_ReturnsMarshalError(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "callback", Type: differ.NodeAdded, NewValue: func() {}},
	}

	_, err := FormatJSON(nodes)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
