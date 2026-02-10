package formatter


import (
	"fmt"
	"sort"
	"strings"
	"reflect"
)



type DiffNode struct {
	Type    string      `json:"type"`
	Key     string      `json:"key"`
	Value   any `json:"value,omitempty"`
	OldVal  any `json:"old_value,omitempty"`
	NewVal  any `json:"new_value,omitempty"`
	Children []*DiffNode `json:"children,omitempty"`
}

func BuildDiff(a, b map[string]any) []*DiffNode {
	keys := collectKeys(a, b)
	sort.Strings(keys)

	var diff []*DiffNode
	for _, key := range keys {
		valA, okA := a[key]
		valB, okB := b[key]

		switch {
		case !okA:
			diff = append(diff, &DiffNode{Type: "added", Key: key, Value: valB})
		case !okB:
			diff = append(diff, &DiffNode{Type: "removed", Key: key, Value: valA})
		case reflect.DeepEqual(valA, valB):
			diff = append(diff, &DiffNode{Type: "unchanged", Key: key, Value: valA})
		case isMap(valA) && isMap(valB):
			nestedA := valA.(map[string]any)
			nestedB := valB.(map[string]any)
			children := BuildDiff(nestedA, nestedB)
			if len(children) > 0 {
				diff = append(diff, &DiffNode{Type: "nested", Key: key, Children: children})
			} else {
				diff = append(diff, &DiffNode{Type: "unchanged", Key: key, Value: valA})
			}
		default:
			diff = append(diff, &DiffNode{Type: "updated", Key: key, OldVal: valA, NewVal: valB})
		}
	}
	return diff
}


func collectKeys(a, b map[string]any) []string {
	keys := make(map[string]bool)
	for k := range a {
		keys[k] = true
	}
	for k := range b {
		keys[k] = true
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	return result
}

func isMap(val any) bool {
	_, ok := val.(map[string]any)
	return ok
}


func FormatStylish(nodes []*DiffNode, indent int) string {
	var lines []string
	prefix := strings.Repeat(" ", indent)

	for _, node := range nodes {
		switch node.Type {
		case "added":
			lines = append(lines, fmt.Sprintf("%s  + %s: %v", prefix, node.Key, FormatValue(node.Value, indent+4)))
		case "removed":
			lines = append(lines, fmt.Sprintf("%s  - %s: %v", prefix, node.Key, FormatValue(node.Value, indent+4)))
		case "updated":
			lines = append(lines, fmt.Sprintf("%s  - %s: %v", prefix, node.Key, FormatValue(node.OldVal, indent+4)))
			lines = append(lines, fmt.Sprintf("%s  + %s: %v", prefix, node.Key, FormatValue(node.NewVal, indent+4)))
		case "unchanged":
			lines = append(lines, fmt.Sprintf("%s    %s: %v", prefix, node.Key, FormatValue(node.Value, indent+4)))
		case "nested":
			lines = append(lines, fmt.Sprintf("%s  %s: {", prefix, node.Key))
			lines = append(lines, FormatStylish(node.Children, indent+4))
			lines = append(lines, fmt.Sprintf("%s  }", prefix))
		}
	}

	result := strings.Join(lines, "\n")

	if indent == 0 {
		if len(lines) == 0 {
			return "{}"
		}
		return "{\n" + result + "\n}"
	}

	return result
}


func FormatValue(val any, indent int) string {
	if m, ok := val.(map[string]any); ok {
		if len(m) == 0 {
			return "{}"
		}
		return "{...}"
	}
	return fmt.Sprintf("%v", val)
}


func FormatPlain(nodes []*DiffNode, path string) string {
	var lines []string

	for _, node := range nodes {
		currentPath := node.Key
		if path != "" {
			currentPath = path + "." + node.Key
		}

		switch node.Type {
		case "added":
			lines = append(lines, fmt.Sprintf("Property '%s' was added with value: %v", currentPath, node.Value))
		case "removed":
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", currentPath))
		case "updated":
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %v to %v", currentPath, node.OldVal, node.NewVal))
		case "nested":
			lines = append(lines, FormatPlain(node.Children, currentPath))
		}
	}

	return strings.Join(lines, "\n")
}

