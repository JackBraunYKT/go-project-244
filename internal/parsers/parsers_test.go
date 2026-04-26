package parsers

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestParse_JSON(t *testing.T) {
	data := []byte(`{"host": "hexlet.io", "timeout": 50, "verbose": true}`)

	result, err := Parse(data, ".json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["host"] != "hexlet.io" {
		t.Errorf("expected host = hexlet.io, got %v", result["host"])
	}
	if result["timeout"] != float64(50) {
		t.Errorf("expected timeout = 50, got %v", result["timeout"])
	}
	if result["verbose"] != true {
		t.Errorf("expected verbose = true, got %v", result["verbose"])
	}
}

func TestParse_YAML(t *testing.T) {
	data := []byte("host: hexlet.io\ntimeout: 50\nverbose: true\n")

	result, err := Parse(data, ".yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["host"] != "hexlet.io" {
		t.Errorf("expected host = hexlet.io, got %v", result["host"])
	}
	if result["timeout"] != 50 {
		t.Errorf("expected timeout = 50, got %v", result["timeout"])
	}
	if result["verbose"] != true {
		t.Errorf("expected verbose = true, got %v", result["verbose"])
	}
}

func TestParse_YML(t *testing.T) {
	data := []byte("host: hexlet.io\n")

	result, err := Parse(data, ".yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["host"] != "hexlet.io" {
		t.Errorf("expected host = hexlet.io, got %v", result["host"])
	}
}

func TestParse_UnsupportedExt(t *testing.T) {
	_, err := Parse([]byte("data"), ".xml")
	if err == nil {
		t.Fatal("expected error for unsupported ext, got nil")
	}
	if err.Error() != "unsupported ext: .xml" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse([]byte(`{invalid json`), ".json")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
	// проверяем что ошибка обёрнута через %w и содержит нужный текст
	if err.Error()[:21] != "failed to parse JSON:" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	_, err := Parse([]byte("invalid: yaml: :\n"), ".yaml")
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestParse_EmptyJSON(t *testing.T) {
	result, err := Parse([]byte(`{}`), ".json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParse_NestedJSON(t *testing.T) {
	data := []byte(`{"common": {"setting1": "Value 1", "setting2": 200}}`)

	result, err := Parse(data, ".json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	nested, ok := result["common"].(map[string]interface{})
	if !ok {
		t.Fatal("expected common to be a map")
	}
	if nested["setting1"] != "Value 1" {
		t.Errorf("expected setting1 = 'Value 1', got %v", nested["setting1"])
	}
}

func TestParse_JSONErrorWrapped(t *testing.T) {
	_, err := Parse([]byte(`{invalid`), ".json")

	var syntaxErr *json.SyntaxError // проверяем что исходная ошибка доступна через errors.As
	if !errors.As(err, &syntaxErr) {
		t.Errorf("expected wrapped json.SyntaxError, got %T: %v", err, err)
	}
}
