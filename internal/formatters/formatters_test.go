package formatters

import (
	"code/internal/differ"
	"strings"
	"testing"
)

func TestFormatNodes_Stylish(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	result, err := FormatNodes(nodes, Stylish, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !strings.Contains(*result, "+ host: hexlet.io") {
		t.Errorf("expected stylish output, got:\n%s", *result)
	}
}

func TestFormatNodes_Plain(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	result, err := FormatNodes(nodes, Plain, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestFormatNodes_JSON(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "host", Type: differ.NodeAdded, NewValue: "hexlet.io"},
	}

	result, err := FormatNodes(nodes, JSON, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !strings.Contains(*result, "hexlet.io") {
		t.Errorf("expected json output, got:\n%s", *result)
	}
}

func TestFormatNodes_UnsupportedFormat(t *testing.T) {
	nodes := []differ.DiffNode{}

	result, err := FormatNodes(nodes, "xml", 1)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if result != nil {
		t.Error("expected nil result on error")
	}
	if err.Error() != "unsupported format: xml" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFormatNodes_ReturnsPointer(t *testing.T) {
	nodes := []differ.DiffNode{}

	result, err := FormatNodes(nodes, Stylish, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil pointer")
	}
}

func TestFormatNodes_SupportedFormatsContainsAll(t *testing.T) {
	expected := []string{Stylish, Plain, JSON}

	for _, format := range expected {
		found := false
		for _, supported := range SupportedFormats {
			if supported == format {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %q in SupportedFormats", format)
		}
	}
}

func TestFormatNodes_AllSupportedFormatsWork(t *testing.T) {
	nodes := []differ.DiffNode{
		{Key: "key", Type: differ.NodeAdded, NewValue: "value"},
	}

	for _, format := range SupportedFormats {
		result, err := FormatNodes(nodes, format, 1)
		if err != nil {
			t.Errorf("format %q returned error: %v", format, err)
		}
		if result == nil {
			t.Errorf("format %q returned nil result", format)
		}
	}
}
