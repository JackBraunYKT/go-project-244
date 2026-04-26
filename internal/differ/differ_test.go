package differ

import (
	"reflect"
	"testing"
)

func TestBuildDiffNodes_Added(t *testing.T) {
	data1 := map[string]interface{}{}
	data2 := map[string]interface{}{"host": "hexlet.io"}

	nodes := BuildDiffNodes(data1, data2)

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if nodes[0].Type != NodeAdded {
		t.Errorf("expected type %q, got %q", NodeAdded, nodes[0].Type)
	}
	if nodes[0].NewValue != "hexlet.io" {
		t.Errorf("expected NewValue = 'hexlet.io', got %v", nodes[0].NewValue)
	}
	if nodes[0].OldValue != nil {
		t.Errorf("expected OldValue = nil, got %v", nodes[0].OldValue)
	}
}

func TestBuildDiffNodes_Removed(t *testing.T) {
	data1 := map[string]interface{}{"host": "hexlet.io"}
	data2 := map[string]interface{}{}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeRemoved {
		t.Errorf("expected type %q, got %q", NodeRemoved, nodes[0].Type)
	}
	if nodes[0].OldValue != "hexlet.io" {
		t.Errorf("expected OldValue = 'hexlet.io', got %v", nodes[0].OldValue)
	}
	if nodes[0].NewValue != nil {
		t.Errorf("expected NewValue = nil, got %v", nodes[0].NewValue)
	}
}

func TestBuildDiffNodes_Unchanged(t *testing.T) {
	data1 := map[string]interface{}{"host": "hexlet.io"}
	data2 := map[string]interface{}{"host": "hexlet.io"}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeUnchanged {
		t.Errorf("expected type %q, got %q", NodeUnchanged, nodes[0].Type)
	}
	if nodes[0].OldValue != "hexlet.io" {
		t.Errorf("expected OldValue = 'hexlet.io', got %v", nodes[0].OldValue)
	}
}

func TestBuildDiffNodes_Changed(t *testing.T) {
	data1 := map[string]interface{}{"timeout": float64(20)}
	data2 := map[string]interface{}{"timeout": float64(50)}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeChanged {
		t.Errorf("expected type %q, got %q", NodeChanged, nodes[0].Type)
	}
	if nodes[0].OldValue != float64(20) {
		t.Errorf("expected OldValue = 20, got %v", nodes[0].OldValue)
	}
	if nodes[0].NewValue != float64(50) {
		t.Errorf("expected NewValue = 50, got %v", nodes[0].NewValue)
	}
}

func TestBuildDiffNodes_Nested(t *testing.T) {
	data1 := map[string]interface{}{
		"common": map[string]interface{}{"setting1": "Value 1"},
	}
	data2 := map[string]interface{}{
		"common": map[string]interface{}{"setting1": "Value 2"},
	}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeNested {
		t.Fatalf("expected type %q, got %q", NodeNested, nodes[0].Type)
	}
	if len(nodes[0].Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(nodes[0].Children))
	}
	if nodes[0].Children[0].Type != NodeChanged {
		t.Errorf("expected child type %q, got %q", NodeChanged, nodes[0].Children[0].Type)
	}
}

func TestBuildDiffNodes_SortedKeys(t *testing.T) {
	data1 := map[string]interface{}{
		"z": "last",
		"a": "first",
		"m": "middle",
	}
	data2 := map[string]interface{}{
		"z": "last",
		"a": "first",
		"m": "middle",
	}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Key != "a" || nodes[1].Key != "m" || nodes[2].Key != "z" {
		t.Errorf("expected sorted keys [a, m, z], got [%s, %s, %s]",
			nodes[0].Key, nodes[1].Key, nodes[2].Key)
	}
}

func TestBuildDiffNodes_MergedKeys(t *testing.T) {
	data1 := map[string]interface{}{"a": 1}
	data2 := map[string]interface{}{"b": 2}

	nodes := BuildDiffNodes(data1, data2)

	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Key != "a" || nodes[1].Key != "b" {
		t.Errorf("expected keys [a, b], got [%s, %s]", nodes[0].Key, nodes[1].Key)
	}
}

func TestBuildDiffNodes_BothEmpty(t *testing.T) {
	nodes := BuildDiffNodes(map[string]interface{}{}, map[string]interface{}{})

	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestBuildDiffNodes_MapVsNonMap(t *testing.T) {
	// если один map а другой нет — должен быть NodeChanged, не NodeNested
	data1 := map[string]interface{}{
		"key": map[string]interface{}{"nested": "val"},
	}
	data2 := map[string]interface{}{
		"key": "string",
	}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeChanged {
		t.Errorf("expected type %q when map vs non-map, got %q", NodeChanged, nodes[0].Type)
	}
}

func TestBuildDiffNodes_DeepNested(t *testing.T) {
	data1 := map[string]interface{}{
		"common": map[string]interface{}{
			"setting6": map[string]interface{}{
				"key": "value",
			},
		},
	}
	data2 := map[string]interface{}{
		"common": map[string]interface{}{
			"setting6": map[string]interface{}{
				"key": "new value",
			},
		},
	}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeNested {
		t.Fatalf("expected nested at level 1")
	}
	if nodes[0].Children[0].Type != NodeNested {
		t.Fatalf("expected nested at level 2")
	}
	if nodes[0].Children[0].Children[0].Type != NodeChanged {
		t.Errorf("expected changed at level 3, got %q", nodes[0].Children[0].Children[0].Type)
	}
}

func TestBuildDiffNodes_NilValue(t *testing.T) {
	data1 := map[string]interface{}{"key": nil}
	data2 := map[string]interface{}{"key": nil}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeUnchanged {
		t.Errorf("expected unchanged for nil == nil, got %q", nodes[0].Type)
	}
}

func TestBuildDiffNodes_KeyPreserved(t *testing.T) {
	data1 := map[string]interface{}{"myKey": "val"}
	data2 := map[string]interface{}{}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Key != "myKey" {
		t.Errorf("expected Key = 'myKey', got %q", nodes[0].Key)
	}
}

func TestBuildDiffNodes_DeepEqualSlice(t *testing.T) {
	val := []interface{}{"a", "b"}
	data1 := map[string]interface{}{"list": val}
	data2 := map[string]interface{}{"list": val}

	nodes := BuildDiffNodes(data1, data2)

	if nodes[0].Type != NodeUnchanged {
		t.Errorf("expected unchanged for equal slices, got %q", nodes[0].Type)
	}
	if !reflect.DeepEqual(nodes[0].OldValue, val) {
		t.Errorf("expected OldValue = %v, got %v", val, nodes[0].OldValue)
	}
}
