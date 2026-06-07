package formatters

import (
	"code/internal/differ"
	"encoding/json"
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

	expected := `{
  "key": "",
  "type": "root",
  "children": [
    {
      "key": "common",
      "type": "nested",
      "children": [
        {
          "key": "follow",
          "type": "added",
          "value2": false
        },
        {
          "key": "setting3",
          "type": "changed",
          "value1": true
        }
      ]
    },
    {
      "key": "host",
      "type": "unchanged",
      "value1": "hexlet.io"
    },
    {
      "key": "proxy",
      "type": "deleted"
    },
    {
      "key": "timeout",
      "type": "changed",
      "value1": 50,
      "value2": 20
    },
    {
      "key": "verbose",
      "type": "added",
      "value2": true
    }
  ]
}`

	if got != expected {
		t.Errorf("unexpected json output:\n%s", got)
	}
}

func TestFormatJSON_CanBeDecodedAsObject(t *testing.T) {
	got, err := FormatJSON([]differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(got), &result); err != nil {
		t.Fatalf("expected JSON object output, got error: %v", err)
	}

	if result["type"] != "root" {
		t.Errorf("expected root object type, got %v", result["type"])
	}
}

func TestFormatJSON_EmptyDiff(t *testing.T) {
	got, err := FormatJSON([]differ.DiffNode{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{
  "key": "",
  "type": "root"
}`

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
