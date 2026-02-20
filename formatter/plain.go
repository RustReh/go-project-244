package formatter

import (
	"fmt"
	"strings"
)

func FormatPlain(nodes []*DiffNode, path string) string {
	var lines []string

	for _, node := range nodes {
		currentPath := node.Key
		if path != "" {
			currentPath = path + "." + node.Key
		}

		switch node.Type {
		case "added":
			lines = append(lines, fmt.Sprintf("Property '%s' was added with value: %s",
				currentPath, formatPlainValue(node.Value)))

		case "removed":
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", currentPath))

		case "updated":
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s",
				currentPath, formatPlainValue(node.OldVal), formatPlainValue(node.NewVal)))

		case "nested":
			nested := FormatPlain(node.Children, currentPath)
			if nested != "" {
				lines = append(lines, nested)
			}
		}
	}

	return strings.Join(lines, "\n")
}

func formatPlainValue(value any) string {
	switch v := value.(type) {
	case map[string]any:
		return "[complex value]"
	case []any:
		return "[complex value]"
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}
