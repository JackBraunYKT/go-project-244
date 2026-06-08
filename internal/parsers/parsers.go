package parsers

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Parser interface {
	Parse([]byte) (map[string]interface{}, error)
}

type JSONParser struct{}

// Parse разбирает JSON-данные в карту конфигурации.
func (JSONParser) Parse(data []byte) (map[string]interface{}, error) {
	return parseJSON(data)
}

type YAMLParser struct{}

// Parse разбирает YAML-данные в карту конфигурации.
func (YAMLParser) Parse(data []byte) (map[string]interface{}, error) {
	return parseYAML(data)
}

// NewParser возвращает парсер для указанного расширения файла.
func NewParser(ext string) (Parser, error) {
	switch ext {
	case ".json":
		return JSONParser{}, nil
	case ".yaml", ".yml":
		return YAMLParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported ext: %s", ext)
	}
}

// Parse разбирает данные в карту конфигурации с помощью парсера, выбранного по расширению файла.
func Parse(data []byte, ext string) (map[string]interface{}, error) {
	parser, err := NewParser(ext)
	if err != nil {
		return nil, err
	}

	return parser.Parse(data)
}

func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}

func parseYAML(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return result, nil
}
