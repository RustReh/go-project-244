package formatter

import (
	"fmt"
	"sort"
	"strings"
)

const indentSize = 4

func FormatStylish(nodes []*DiffNode, depth int) string {
	if len(nodes) == 0 {
		return "{}"
	}

	lines := []string{"{"}

	for _, node := range nodes {
		lines = append(lines, formatNode(node, depth))
	}

	// Closing brace indent: depth * indentSize (0 spaces for root level)
	closingIndent := ""
	if depth > 0 {
		closingIndent = strings.Repeat(" ", depth*indentSize)
	}
	lines = append(lines, fmt.Sprintf("%s}", closingIndent))
	return strings.Join(lines, "\n")
}

func formatNode(node *DiffNode, depth int) string {
	// Base indent for properties at this depth: (depth + 1) * indentSize
	baseIndent := (depth + 1) * indentSize

	switch node.Type {
	case "added":
		markerIndent := strings.Repeat(" ", baseIndent-2)
		return fmt.Sprintf("%s+ %s: %s",
			markerIndent,
			node.Key,
			FormatValue(node.Value, depth+1),
		)
	case "removed":
		markerIndent := strings.Repeat(" ", baseIndent-2)
		return fmt.Sprintf("%s- %s: %s",
			markerIndent,
			node.Key,
			FormatValue(node.Value, depth+1),
		)
	case "unchanged":
		propIndent := strings.Repeat(" ", baseIndent)
		return fmt.Sprintf("%s%s: %s",
			propIndent,
			node.Key,
			FormatValue(node.Value, depth+1),
		)
	case "updated":
		markerIndent := strings.Repeat(" ", baseIndent-2)
		line1 := fmt.Sprintf("%s- %s: %s",
			markerIndent,
			node.Key,
			FormatValue(node.OldVal, depth+1),
		)
		line2 := fmt.Sprintf("%s+ %s: %s",
			markerIndent,
			node.Key,
			FormatValue(node.NewVal, depth+1),
		)
		return line1 + "\n" + line2
	case "nested":
		propIndent := strings.Repeat(" ", baseIndent)
		return fmt.Sprintf("%s%s: %s",
			propIndent,
			node.Key,
			FormatStylish(node.Children, depth+1),
		)
	}

	return ""
}

func FormatValue(value any, depth int) string {
	switch v := value.(type) {
	case map[string]any:
		if len(v) == 0 {
			return "{}"
		}
		return formatMap(v, depth)
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatMap(m map[string]any, depth int) string {
	// Properties inside the map should be at (depth + 2) * indentSize
	propIndent := strings.Repeat(" ", (depth+2)*indentSize)
	// Closing brace for the map should be at (depth + 1) * indentSize
	closingIndent := strings.Repeat(" ", (depth+1)*indentSize)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := []string{"{"}

	for _, k := range keys {
		val := FormatValue(m[k], depth+1)
		lines = append(lines,
			fmt.Sprintf("%s%s: %s", propIndent, k, val),
		)
	}

	lines = append(lines, fmt.Sprintf("%s}", closingIndent))
	return strings.Join(lines, "\n")
}
