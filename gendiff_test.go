package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type diffTestCase struct {
	name         string
	file1        string
	file2        string
	expectedFile string
	expectErr    bool
	format       string
}

func TestGenDiff_JSON(t *testing.T) {
	tests := []diffTestCase{
		{
			name:         "Flat diff",
			file1:        fixturePath("flat_file1.json"),
			file2:        fixturePath("flat_file2.json"),
			expectedFile: "flat_diff.txt",
		},
		{
			name:         "Same files",
			file1:        fixturePath("same1.json"),
			file2:        fixturePath("same2.json"),
			expectedFile: "same.txt",
		},
		{
			name:         "Empty files",
			file1:        fixturePath("empty1.json"),
			file2:        fixturePath("empty2.json"),
			expectedFile: "empty.txt",
		},
		{
			name:         "Empty vs filled",
			file1:        fixturePath("empty_vs_filled1.json"),
			file2:        fixturePath("empty_vs_filled2.json"),
			expectedFile: "empty_vs_filled.txt",
		},
		{
			name:         "Only deleted",
			file1:        fixturePath("only_deleted.json"),
			file2:        fixturePath("only_added.json"),
			expectedFile: "only_deleted.txt",
		},
		{
			name:         "Only added",
			file1:        fixturePath("only_added.json"),
			file2:        fixturePath("only_deleted.json"),
			expectedFile: "only_added.txt",
		},
		{
			name:         "Different types",
			file1:        fixturePath("different_types.json"),
			file2:        fixturePath("different_types2.json"),
			expectedFile: "different_types.txt",
		},
		{
			name:         "Completely different",
			file1:        fixturePath("completely_different1.json"),
			file2:        fixturePath("completely_different2.json"),
			expectedFile: "completely_different.txt",
		},
	}

	runDiffTests(t, tests)
}

func TestGenDiff_YAML(t *testing.T) {
	tests := []diffTestCase{
		{
			name:         "Flat diff",
			file1:        fixturePath("flat_file1.yml"),
			file2:        fixturePath("flat_file2.yml"),
			expectedFile: "flat_diff.txt",
		},
		{
			name:         "Same files",
			file1:        fixturePath("same1.yml"),
			file2:        fixturePath("same2.yml"),
			expectedFile: "same.txt",
		},
		{
			name:         "Empty files",
			file1:        fixturePath("empty1.yml"),
			file2:        fixturePath("empty2.yml"),
			expectedFile: "empty.txt",
		},
		{
			name:         "Empty vs filled",
			file1:        fixturePath("empty_vs_filled1.yml"),
			file2:        fixturePath("empty_vs_filled2.yml"),
			expectedFile: "empty_vs_filled.txt",
		},
		{
			name:         "Only deleted",
			file1:        fixturePath("only_deleted.yml"),
			file2:        fixturePath("only_added.yml"),
			expectedFile: "only_deleted.txt",
		},
		{
			name:         "Only added",
			file1:        fixturePath("only_added.yml"),
			file2:        fixturePath("only_deleted.yml"),
			expectedFile: "only_added.txt",
		},
		{
			name:         "Different types",
			file1:        fixturePath("different_types.yml"),
			file2:        fixturePath("different_types2.yml"),
			expectedFile: "different_types_yaml.txt",
		},
		{
			name:         "Completely different",
			file1:        fixturePath("completely_different1.yml"),
			file2:        fixturePath("completely_different2.yml"),
			expectedFile: "completely_different.txt",
		},
	}

	runDiffTests(t, tests)
}

func TestGenDiff_Mixed(t *testing.T) {
	tests := []diffTestCase{
		{
			name:         "YAML and JSON mixed",
			file1:        fixturePath("file1.yml"),
			file2:        fixturePath("file2.json"),
			expectedFile: "flat_diff.txt",
		},
		{
			name:         "JSON and YAML mixed",
			file1:        fixturePath("file1.json"),
			file2:        fixturePath("file2.yml"),
			expectedFile: "flat_diff.txt",
		},
	}

	runDiffTests(t, tests)
}

func TestGenDiff_Plain(t *testing.T) {
	tests := []diffTestCase{
		{
			name:         "Flat diff plain format",
			file1:        fixturePath("file1.json"),
			file2:        fixturePath("file2.json"),
			expectedFile: "flat_diff_plain.txt",
			format:       "plain",
		},
		{
			name:         "Nested diff plain format",
			file1:        fixturePath("nested1.json"),
			file2:        fixturePath("nested2.json"),
			expectedFile: "nested_plain.txt",
			format:       "plain",
		},
	}

	runDiffTests(t, tests)
}

func TestGenDiff_JSONFormat(t *testing.T) {
	tests := []diffTestCase{
		{
			name:         "Flat diff json format",
			file1:        fixturePath("file1.json"),
			file2:        fixturePath("file2.json"),
			expectedFile: "flat_diff_json.txt",
			format:       "json",
		},
		{
			name:         "Nested diff json format",
			file1:        fixturePath("nested1.json"),
			file2:        fixturePath("nested2.json"),
			expectedFile: "nested_json.txt",
			format:       "json",
		},
	}

	runDiffTests(t, tests)
}

func TestGenDiff_Errors(t *testing.T) {
	tests := []diffTestCase{
		{
			name:      "Nonexistent file1",
			file1:     fixturePath("nonexistent.json"),
			file2:     fixturePath("file1.json"),
			expectErr: true,
		},
		{
			name:      "Nonexistent file2",
			file1:     fixturePath("file1.json"),
			file2:     fixturePath("nonexistent.json"),
			expectErr: true,
		},
		{
			name:      "Both nonexistent files",
			file1:     fixturePath("nonexistent1.json"),
			file2:     fixturePath("nonexistent2.json"),
			expectErr: true,
		},
		{
			name:      "Unsupported format",
			file1:     fixturePath("file1.txt"),
			file2:     fixturePath("file2.json"),
			expectErr: true,
		},
		{
			name:      "Unknown output format",
			file1:     fixturePath("file1.json"),
			file2:     fixturePath("file2.json"),
			format:    "xml",
			expectErr: true,
		},
	}

	runDiffTests(t, tests)
}

func TestGenDiff_EmptyPaths(t *testing.T) {
	_, err := GenDiff("", "file.json", "")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmptyPath)

	_, err = GenDiff("file.json", "", "")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrEmptyPath)
}

func TestGenDiff_Nested(t *testing.T) {
	tests := []diffTestCase{
		{
			name:         "Nested JSON structures",
			file1:        fixturePath("nested1.json"),
			file2:        fixturePath("nested2.json"),
			expectedFile: "nested.txt",
		},
		{
			name:         "Nested YAML structures",
			file1:        fixturePath("nested1.yml"),
			file2:        fixturePath("nested2.yml"),
			expectedFile: "nested.txt",
		},
	}

	runDiffTests(t, tests)
}

func runDiffTests(t *testing.T, tests []diffTestCase) {
	t.Helper()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenDiff(tt.file1, tt.file2, tt.format)
			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, tt.expectedFile)

			expected := readExpected(t, tt.expectedFile)
			require.Equal(t, expected, result)
		})
	}
}

func readExpected(t *testing.T, filename string) string {
	t.Helper()

	path := filepath.Join("testdata", "expected", filename)
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	normalized := strings.ReplaceAll(string(content), "\r\n", "\n")
	return strings.TrimRight(normalized, "\n")
}

func fixturePath(segments ...string) string {
	return filepath.Join(append([]string{"testdata", "fixture"}, segments...)...)
}
