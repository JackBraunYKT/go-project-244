package code

import (
	"code/internal/differ"
	"code/internal/formatters"
	"code/internal/parsers"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Format string

const (
	Stylish Format = "stylish"
)

var (
	ErrEmptyPath = errors.New("file path cannot be empty")
)

func GenDiff(filepath1, filepath2 string, format string) (string, error) {
	if format == "" {
		format = formatters.Stylish
	}

	if filepath1 == "" {
		return "", fmt.Errorf("first file: %w", ErrEmptyPath)
	}
	if filepath2 == "" {
		return "", fmt.Errorf("second file: %w", ErrEmptyPath)
	}

	data1, err := readFile(filepath1)
	if err != nil {
		return "", err
	}

	data2, err := readFile(filepath2)
	if err != nil {
		return "", err
	}

	ext1 := filepath.Ext(filepath1)
	ext2 := filepath.Ext(filepath2)

	parser1, err := parsers.NewParser(ext1)
	if err != nil {
		return "", err
	}

	parsed1, err := parser1.Parse(data1)
	if err != nil {
		return "", err
	}

	parser2, err := parsers.NewParser(ext2)
	if err != nil {
		return "", err
	}

	parsed2, err := parser2.Parse(data2)
	if err != nil {
		return "", err
	}

	nodes := differ.BuildDiffNodes(parsed1, parsed2)
	formatter, err := formatters.NewFormatter(format)
	if err != nil {
		return "", err
	}

	result, err := formatter.Format(nodes)
	if err != nil {
		return "", err
	}

	return result, nil
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
