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
	errContains  string
	format       string
}

type baseDiffTestCase struct {
	name         string
	file1        string
	file2        string
	expectedFile string
}

type fixtureFileSet struct {
	extension       string
	extensionByFile map[string]string
	expectedByCase  map[string]string
}

var commonDiffTestCases = []baseDiffTestCase{
	{"Flat diff", "flat_file1", "flat_file2", "flat_diff.txt"},
	{"Same files", "same1", "same2", "same.txt"},
	{"Empty files", "empty1", "empty2", "empty.txt"},
	{"Empty vs filled", "empty_vs_filled1", "empty_vs_filled2", "empty_vs_filled.txt"},
	{"Only deleted", "only_deleted", "only_added", "only_deleted.txt"},
	{"Only added", "only_added", "only_deleted", "only_added.txt"},
	{"Different types", "different_types", "different_types2", "different_types.txt"},
	{"Completely different", "completely_different1", "completely_different2", "completely_different.txt"},
}

func TestGenDiff_JSON(t *testing.T) {
	runDiffTests(t, buildFixtureDiffCases(fixtureFileSet{extension: "json"}))
}

func TestGenDiff_YAML(t *testing.T) {
	runDiffTests(t, buildFixtureDiffCases(fixtureFileSet{
		extension: "yml",
		extensionByFile: map[string]string{
			"flat_file1": "yaml",
			"flat_file2": "yaml",
		},
		expectedByCase: map[string]string{
			"Different types": "different_types_yaml.txt",
		},
	}))
}

func TestGenDiff_Mixed(t *testing.T) {
	tests := []diffTestCase{
		diffCase("YAML and JSON mixed", "flat_file1.yaml", "flat_file2.json", "flat_diff.txt"),
		diffCase("JSON and YAML mixed", "flat_file1.json", "flat_file2.yaml", "flat_diff.txt"),
	}

	runDiffTests(t, tests)
}

func TestGenDiff_Plain(t *testing.T) {
	tests := []diffTestCase{
		formattedDiffCase("Flat diff plain format", "flat_file1.json", "flat_file2.json", "flat_diff_plain.txt", "plain"),
		formattedDiffCase("Nested diff plain format", "nested1.json", "nested2.json", "nested_plain.txt", "plain"),
	}

	runDiffTests(t, tests)
}

func TestGenDiff_JSONFormat(t *testing.T) {
	tests := []diffTestCase{
		formattedDiffCase("Flat diff json format", "flat_file1.json", "flat_file2.json", "flat_diff_json.txt", "json"),
		formattedDiffCase("Nested diff json format", "nested1.json", "nested2.json", "nested_json.txt", "json"),
	}

	runDiffTests(t, tests)
}

func TestGenDiff_Errors(t *testing.T) {
	tests := []diffTestCase{
		errorDiffCase("Nonexistent file1", "nonexistent.json", "flat_file1.json", "failed to stat file"),
		errorDiffCase("Nonexistent file2", "flat_file1.json", "nonexistent.json", "failed to stat file"),
		errorDiffCase("Both nonexistent files", "nonexistent1.json", "nonexistent2.json", "failed to stat file"),
		errorDiffCase("Unsupported input extension", "unsupported.txt", "flat_file2.json", "unsupported ext: .txt"),
		formattedErrorDiffCase("Unknown output format", "flat_file1.json", "flat_file2.json", "xml", "unsupported format: xml"),
	}

	runDiffTests(t, tests)
}

func TestGenDiff_EmptyPaths(t *testing.T) {
	tests := []struct {
		name  string
		file1 string
		file2 string
	}{
		{name: "Empty file1", file1: "", file2: "file.json"},
		{name: "Empty file2", file1: "file.json", file2: ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenDiff(tt.file1, tt.file2, "")
			require.Error(t, err)
			require.ErrorIs(t, err, ErrEmptyPath)
		})
	}
}

func TestGenDiff_Nested(t *testing.T) {
	tests := []diffTestCase{
		diffCase("Nested JSON structures", "nested1.json", "nested2.json", "nested.txt"),
		diffCase("Nested YAML structures", "nested1.yml", "nested2.yml", "nested.txt"),
	}

	runDiffTests(t, tests)
}

func buildFixtureDiffCases(files fixtureFileSet) []diffTestCase {
	tests := make([]diffTestCase, 0, len(commonDiffTestCases))

	for _, tt := range commonDiffTestCases {
		expectedFile := tt.expectedFile
		if override, ok := files.expectedByCase[tt.name]; ok {
			expectedFile = override
		}

		tests = append(tests, diffCase(
			tt.name,
			fixtureFileName(tt.file1, files),
			fixtureFileName(tt.file2, files),
			expectedFile,
		))
	}

	return tests
}

func fixtureFileName(filename string, files fixtureFileSet) string {
	extension := files.extension
	if override, ok := files.extensionByFile[filename]; ok {
		extension = override
	}

	return filename + "." + extension
}

func diffCase(name, file1, file2, expectedFile string) diffTestCase {
	return diffTestCase{
		name:         name,
		file1:        fixturePath(file1),
		file2:        fixturePath(file2),
		expectedFile: expectedFile,
	}
}

func formattedDiffCase(name, file1, file2, expectedFile, format string) diffTestCase {
	test := diffCase(name, file1, file2, expectedFile)
	test.format = format

	return test
}

func errorDiffCase(name, file1, file2, errContains string) diffTestCase {
	return diffTestCase{
		name:        name,
		file1:       fixturePath(file1),
		file2:       fixturePath(file2),
		expectErr:   true,
		errContains: errContains,
	}
}

func formattedErrorDiffCase(name, file1, file2, format, errContains string) diffTestCase {
	test := errorDiffCase(name, file1, file2, errContains)
	test.format = format

	return test
}

func runDiffTests(t *testing.T, tests []diffTestCase) {
	t.Helper()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenDiff(tt.file1, tt.file2, tt.format)
			if tt.expectErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.ErrorContains(t, err, tt.errContains)
				}
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
