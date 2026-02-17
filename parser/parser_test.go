package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHost = "hexlet.io"

func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644), "failed to create temp file")
	return path
}

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    map[string]any
		expectError bool
	}{
		{
			name: "valid JSON",
			input: `{
				"host": "` + testHost + `",
				"timeout": 50,
				"debug": false
			}`,
			expected: map[string]any{
				"host":    testHost,
				"timeout": 50.0,
				"debug":   false,
			},
			expectError: false,
		},
		{
			name:        "invalid JSON syntax",
			input:       `{ "host": "` + testHost + `", }`,
			expectError: true,
		},
		{
			name:        "empty JSON",
			input:       ``,
			expectError: true,
		},
		{
			name: "nested JSON",
			input: `{
				"common": {
					"setting1": "value1",
					"setting2": 200
				}
			}`,
			expected: map[string]any{
				"common": map[string]any{
					"setting1": "value1",
					"setting2": 200.0,
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseJSON([]byte(tt.input))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseYAML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    map[string]any
		expectError bool
	}{
		{
			name: "valid YAML",
			input: `host: ` + testHost + `
timeout: 50
debug: false
`,
			expected: map[string]any{
				"host":    testHost,
				"timeout": 50,
				"debug":   false,
			},
			expectError: false,
		},
		{
			name: "YAML with nested structure",
			input: `common:
  setting1: value1
  setting2: 200
`,
			expected: map[string]any{
				"common": map[string]any{
					"setting1": "value1",
					"setting2": 200,
				},
			},
			expectError: false,
		},
		{
			name:        "invalid YAML syntax",
			input:       "host: " + testHost + "\n  timeout: 50",
			expectError: true,
		},
		{
			name:        "empty YAML",
			input:       "{}", // пустой валидный YAML
			expected:    map[string]any{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseYAML([]byte(tt.input))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseTOML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    map[string]any
		expectError bool
	}{
		{
			name: "valid TOML",
			input: `host = "` + testHost + `"
timeout = 50
debug = false
`,
			expected: map[string]any{
				"host":    testHost,
				"timeout": int64(50),
				"debug":   false,
			},
			expectError: false,
		},
		{
			name: "TOML with nested table",
			input: `[common]
setting1 = "value1"
setting2 = 200
`,
			expected: map[string]any{
				"common": map[string]any{
					"setting1": "value1",
					"setting2": int64(200),
				},
			},
			expectError: false,
		},
		{
			name:        "invalid TOML syntax",
			input:       `host = ` + testHost,
			expectError: true,
		},
		{
			name:        "empty TOML",
			input:       "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTOML([]byte(tt.input))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				if tt.input == "" {
					assert.Equal(t, 0, len(result), "Пустой TOML должен вернуть пустую мапу или nil")
				} else {
					assert.Equal(t, tt.expected, result)
				}
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		filename    string
		content     string
		expected    map[string]any
		expectError bool
		errorMsg    string
	}{
		{
			name:     "JSON file with relative path",
			filename: "config1.json",
			content: `{
				"host": "` + testHost + `",
				"timeout": 50
			}`,
			expected: map[string]any{
				"host":    testHost,
				"timeout": 50.0,
			},
			expectError: false,
		},
		{
			name:     "YAML file (.yaml extension)",
			filename: "config2.yaml",
			content: `host: ` + testHost + `
timeout: 20
`,
			expected: map[string]any{
				"host":    testHost,
				"timeout": 20,
			},
			expectError: false,
		},
		{
			name:     "YAML file (.yml extension)",
			filename: "config3.yml",
			content: `host: ` + testHost + `
verbose: true
`,
			expected: map[string]any{
				"host":    testHost,
				"verbose": true,
			},
			expectError: false,
		},
		{
			name:     "TOML file",
			filename: "config4.toml",
			content: `host = "` + testHost + `"
debug = false
`,
			expected: map[string]any{
				"host":  testHost,
				"debug": false,
			},
			expectError: false,
		},
		{
			name:        "unsupported extension",
			filename:    "config5.txt",
			content:     "some content",
			expectError: true,
			errorMsg:    "unsupported file format",
		},
		{
			name:        "non-existent file",
			filename:    "nonexistent.json",
			expectError: true,
			errorMsg:    "no such file or directory",
		},
		{
			name:     "absolute path",
			filename: "abs.json",
			content:  `{"key": "value"}`,
			expected: map[string]any{"key": "value"},
		},
		{
			name:        "invalid JSON content",
			filename:    "invalid.json",
			content:     `{ "broken": `,
			expectError: true,
			errorMsg:    "invalid JSON",
		},
		{
			name:        "invalid YAML content",
			filename:    "invalid.yaml",
			content:     "key: value: invalid",
			expectError: true,
			errorMsg:    "invalid YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			if tt.filename != "nonexistent.json" {
				filePath = createTempFile(t, tmpDir, tt.filename, tt.content)
			} else {
				filePath = filepath.Join(tmpDir, tt.filename)
			}

			if tt.name == "absolute path" {
				var err error
				filePath, err = filepath.Abs(filePath)
				require.NoError(t, err)
			}

			result, err := ParseFile(filePath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseFile_RelativePathFromDifferentDir(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := createTempFile(t, tmpDir, "config.json", `{"key": "value"}`)

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(originalDir))
	}()

	tempWorkDir := t.TempDir()
	require.NoError(t, os.Chdir(tempWorkDir))

	relPath, err := filepath.Rel(tempWorkDir, configPath)
	require.NoError(t, err)

	result, err := ParseFile(relPath)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"key": "value"}, result)
}

func TestParseFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()

	emptyJSON := createTempFile(t, tmpDir, "empty.json", "")
	_, err := ParseFile(emptyJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")

	emptyYAML := createTempFile(t, tmpDir, "empty.yaml", "")
	resultYAML, err := ParseFile(emptyYAML)
	require.NoError(t, err)
	assert.Equal(t, 0, len(resultYAML), "Пустой YAML должен вернуть пустую мапу или nil")

	emptyTOML := createTempFile(t, tmpDir, "empty.toml", "")
	resultTOML, err := ParseFile(emptyTOML)
	require.NoError(t, err)
	assert.Equal(t, 0, len(resultTOML), "Пустой TOML должен вернуть пустую мапу или nil")
}

func TestParseFile_PathNormalization(t *testing.T) {
	tmpDir := t.TempDir()
	_ = createTempFile(t, tmpDir, "config.json", `{"key": "value"}`)

	dirtyPath := filepath.Join(tmpDir, ".", "config.json")
	result, err := ParseFile(dirtyPath)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"key": "value"}, result)
}
