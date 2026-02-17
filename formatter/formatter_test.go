package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHost = "hexlet.io"

func TestBuildDiff(t *testing.T) {
	tests := []struct {
		name     string
		a        map[string]any
		b        map[string]any
		expected []*DiffNode
	}{
		{
			name: "added key",
			a:    map[string]any{"host": testHost},
			b:    map[string]any{"host": testHost, "timeout": 50},
			expected: []*DiffNode{
				{Type: "unchanged", Key: "host", Value: testHost},
				{Type: "added", Key: "timeout", Value: 50},
			},
		},
		{
			name: "removed key",
			a:    map[string]any{"host": testHost, "timeout": 50},
			b:    map[string]any{"host": testHost},
			expected: []*DiffNode{
				{Type: "unchanged", Key: "host", Value: testHost},
				{Type: "removed", Key: "timeout", Value: 50},
			},
		},
		{
			name: "updated value",
			a:    map[string]any{"timeout": 50},
			b:    map[string]any{"timeout": 20},
			expected: []*DiffNode{
				{Type: "updated", Key: "timeout", OldVal: 50, NewVal: 20},
			},
		},
		{
			name: "nested with changes",
			a: map[string]any{
				"common": map[string]any{"setting1": "value1", "setting2": 200},
			},
			b: map[string]any{
				"common": map[string]any{"setting1": "value1", "setting2": 300},
			},
			expected: []*DiffNode{
				{
					Type: "nested",
					Key:  "common",
					Children: []*DiffNode{
						{Type: "unchanged", Key: "setting1", Value: "value1"},
						{Type: "updated", Key: "setting2", OldVal: 200.0, NewVal: 300.0},
					},
				},
			},
		},
		{
			name:     "empty maps",
			a:        map[string]any{},
			b:        map[string]any{},
			expected: []*DiffNode{},
		},
		{
			name: "complex nested structure",
			a: map[string]any{
				"host":    testHost,
				"timeout": 50,
				"debug":   false,
				"common": map[string]any{
					"follow":   false,
					"setting1": "value1",
				},
			},
			b: map[string]any{
				"host":    testHost,
				"timeout": 20,
				"verbose": true,
				"common": map[string]any{
					"follow":   true,
					"setting1": "value1",
				},
			},
			expected: []*DiffNode{
				{
					Type: "nested",
					Key:  "common",
					Children: []*DiffNode{
						{Type: "updated", Key: "follow", OldVal: false, NewVal: true},
						{Type: "unchanged", Key: "setting1", Value: "value1"},
					},
				},
				{Type: "removed", Key: "debug", Value: false},
				{Type: "unchanged", Key: "host", Value: testHost},
				{Type: "updated", Key: "timeout", OldVal: 50, NewVal: 20},
				{Type: "added", Key: "verbose", Value: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildDiff(tt.a, tt.b)

			assert.Len(t, result, len(tt.expected))

			for i, expectedNode := range tt.expected {
				require.Less(t, i, len(result), "Результат содержит меньше узлов, чем ожидалось")
				actualNode := result[i]

				assert.Equal(t, expectedNode.Type, actualNode.Type, "Тип узла не совпадает для ключа %s", expectedNode.Key)
				assert.Equal(t, expectedNode.Key, actualNode.Key, "Ключ узла не совпадает")

				switch expectedNode.Type {
				case "added":
					assert.Equal(t, expectedNode.Value, actualNode.Value, "Значение добавленного узла не совпадает для ключа %s", expectedNode.Key)
				case "removed":
					assert.Equal(t, expectedNode.Value, actualNode.Value, "Значение удалённого узла не совпадает для ключа %s", expectedNode.Key)
				case "updated":
					assert.Equal(t, expectedNode.OldVal, actualNode.OldVal, "Старое значение не совпадает для ключа %s", expectedNode.Key)
					assert.Equal(t, expectedNode.NewVal, actualNode.NewVal, "Новое значение не совпадает для ключа %s", expectedNode.Key)
				case "unchanged":
					assert.Equal(t, expectedNode.Value, actualNode.Value, "Значение неизменённого узла не совпадает для ключа %s", expectedNode.Key)
				case "nested":
					assert.NotNil(t, actualNode.Children, "Дети вложенного узла должны быть не nil для ключа %s", expectedNode.Key)
					assert.Len(t, actualNode.Children, len(expectedNode.Children), "Количество детей не совпадает для ключа %s", expectedNode.Key)
				}
			}
		})
	}
}

func TestCollectKeys(t *testing.T) {
	tests := []struct {
		name     string
		a        map[string]any
		b        map[string]any
		expected []string
	}{
		{
			name:     "both maps empty",
			a:        map[string]any{},
			b:        map[string]any{},
			expected: []string{},
		},
		{
			name:     "only in a",
			a:        map[string]any{"key1": 1, "key2": 2},
			b:        map[string]any{},
			expected: []string{"key1", "key2"},
		},
		{
			name:     "only in b",
			a:        map[string]any{},
			b:        map[string]any{"key3": 3},
			expected: []string{"key3"},
		},
		{
			name:     "overlap",
			a:        map[string]any{"key1": 1, "key2": 2},
			b:        map[string]any{"key2": 2, "key3": 3},
			expected: []string{"key1", "key2", "key3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collectKeys(tt.a, tt.b)
			assert.ElementsMatch(t, tt.expected, result, "Ключи должны совпадать (порядок не важен)")
		})
	}
}

func TestIsMap(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{
			name:     "map[string]any",
			value:    map[string]any{"key": "value"},
			expected: true,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
		{
			name:     "string",
			value:    "not a map",
			expected: false,
		},
		{
			name:     "int",
			value:    42,
			expected: false,
		},
		{
			name:     "slice",
			value:    []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "empty map",
			value:    map[string]any{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMap(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatStylish(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []*DiffNode
		indent   int
		expected string
	}{
		{
			name: "single added node",
			nodes: []*DiffNode{
				{Type: "added", Key: "timeout", Value: 50},
			},
			indent:   0,
			expected: "{\n  + timeout: 50\n}",
		},
		{
			name: "added and removed",
			nodes: []*DiffNode{
				{Type: "removed", Key: "debug", Value: false},
				{Type: "added", Key: "verbose", Value: true},
			},
			indent:   0,
			expected: "{\n  - debug: false\n  + verbose: true\n}",
		},
		{
			name: "nested structure",
			nodes: []*DiffNode{
				{
					Type: "nested",
					Key:  "common",
					Children: []*DiffNode{
						{Type: "added", Key: "follow", Value: true},
					},
				},
			},
			indent:   0,
			expected: "{\n  common: {\n      + follow: true\n  }\n}",
		},
		{
			name: "with indent",
			nodes: []*DiffNode{
				{Type: "added", Key: "key", Value: "value"},
			},
			indent:   4,
			expected: "      + key: value",
		},
		{
			name:     "empty diff",
			nodes:    []*DiffNode{},
			indent:   0,
			expected: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatStylish(tt.nodes, tt.indent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatPlain(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []*DiffNode
		path     string
		expected string
	}{
		{
			name: "added property",
			nodes: []*DiffNode{
				{Type: "added", Key: "timeout", Value: 50},
			},
			path:     "",
			expected: "Property 'timeout' was added with value: 50",
		},
		{
			name: "removed property",
			nodes: []*DiffNode{
				{Type: "removed", Key: "debug", Value: false},
			},
			path:     "",
			expected: "Property 'debug' was removed",
		},
		{
			name: "updated property",
			nodes: []*DiffNode{
				{Type: "updated", Key: "timeout", OldVal: 50, NewVal: 20},
			},
			path:     "",
			expected: "Property 'timeout' was updated. From 50 to 20",
		},
		{
			name: "nested property",
			nodes: []*DiffNode{
				{
					Type: "nested",
					Key:  "common",
					Children: []*DiffNode{
						{Type: "added", Key: "follow", Value: true},
					},
				},
			},
			path:     "",
			expected: "Property 'common.follow' was added with value: true",
		},
		{
			name: "multiple changes",
			nodes: []*DiffNode{
				{Type: "removed", Key: "debug", Value: false},
				{Type: "added", Key: "verbose", Value: true},
			},
			path:     "",
			expected: "Property 'debug' was removed\nProperty 'verbose' was added with value: true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPlain(tt.nodes, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		indent   int
		expected string
	}{
		{
			name:     "string",
			value:    "hello",
			indent:   0,
			expected: "hello",
		},
		{
			name:     "number",
			value:    42,
			indent:   0,
			expected: "42",
		},
		{
			name:     "boolean",
			value:    true,
			indent:   0,
			expected: "true",
		},
		{
			name:     "empty map",
			value:    map[string]any{},
			indent:   0,
			expected: "{}",
		},
		{
			name:     "non-empty map",
			value:    map[string]any{"key": "value"},
			indent:   0,
			expected: "{...}",
		},
		{
			name:     "nil",
			value:    nil,
			indent:   0,
			expected: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValue(tt.value, tt.indent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntegrationFullCycle(t *testing.T) {
	data1 := map[string]any{
		"host":    testHost,
		"timeout": 50,
		"debug":   false,
	}

	data2 := map[string]any{
		"host":    testHost,
		"timeout": 20,
		"verbose": true,
	}

	diff := BuildDiff(data1, data2)
	assert.Len(t, diff, 4)

	stylish := FormatStylish(diff, 0)
	assert.Contains(t, stylish, "- debug: false")
	assert.Contains(t, stylish, "host: "+testHost)
	assert.Contains(t, stylish, "- timeout: 50")
	assert.Contains(t, stylish, "+ timeout: 20")
	assert.Contains(t, stylish, "+ verbose: true")

	plain := FormatPlain(diff, "")
	assert.Contains(t, plain, "Property 'debug' was removed")
	assert.Contains(t, plain, "Property 'timeout' was updated")
	assert.Contains(t, plain, "Property 'verbose' was added")
}
