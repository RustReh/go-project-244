package parser


import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)


func ParseFile(path string) (map[string]any, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path '%s': %w", path, err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", absPath, err)
	}

	ext := filepath.Ext(absPath)
	switch ext {
	case ".json":
		return parseJSON(data)
	case ".yaml", ".yml":
		return parseYAML(data)
	case ".toml":
		return parseTOML(data)
	default:
		return nil, fmt.Errorf("unsupported file format: %s (expected .json, .yaml, .yml, .toml)", ext)
	}
}


func parseJSON(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return result, nil
}


func parseYAML(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	return result, nil
}


func parseTOML(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := toml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid TOML: %w", err)
	}
	return result, nil
}
