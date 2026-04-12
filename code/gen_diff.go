package code

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
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

	result := genDiff(parsed1, parsed2)

	return &result, nil
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

func genDiff(parsedData1, parsedData2 map[string]interface{}) string {
	fmt.Println(parsedData1)
	fmt.Println(parsedData2)

	keys := getMergedSortedKeys(parsedData1, parsedData2)

	var result string

	for k, v := range parsedData1 {

	}

	return result
}

func getMergedSortedKeys(m1, m2 map[string]interface{}) []string {
	keys := slices.Collect(maps.Keys(m1))

	for k := range m2 {
		if _, ok := m1[k]; !ok {
			keys = append(keys, k)
		}
	}

	slices.Sort(keys)

	return keys
}
