package formatter

import (
	"fmt"
	"sort"
	"strings"
)

func FormatStylish(nodes []*DiffNode, depth int) string {
	if len(nodes) == 0 {
		return "{}"
	}

	indent := getIndent(depth)
	lines := []string{"{"}

	for _, node := range nodes {
		lines = append(lines, formatNode(node, depth+1))
	}

	lines = append(lines, fmt.Sprintf("%s}", indent))
	return strings.Join(lines, "\n")
}

func formatNode(node *DiffNode, depth int) string {
	indent := getIndent(depth)

	switch node.Type {
	case "added":
		return fmt.Sprintf("%s+ %s: %s", indent, node.Key, FormatValue(node.Value, depth))
	case "removed":
		return fmt.Sprintf("%s- %s: %s", indent, node.Key, FormatValue(node.Value, depth))
	case "unchanged":
		return fmt.Sprintf("%s  %s: %s", indent, node.Key, FormatValue(node.Value, depth))
	case "updated":
		line1 := fmt.Sprintf("%s- %s: %s", indent, node.Key, FormatValue(node.OldVal, depth))
		line2 := fmt.Sprintf("%s+ %s: %s", indent, node.Key, FormatValue(node.NewVal, depth))
		return line1 + "\n" + line2
	case "nested":
		lines := []string{fmt.Sprintf("%s  %s: {", indent, node.Key)}
		for _, child := range node.Children {
			lines = append(lines, formatNode(child, depth+1))
		}
		lines = append(lines, fmt.Sprintf("%s  }", indent))
		return strings.Join(lines, "\n")
	}
	return ""
}

func getIndent(depth int) string {
	if depth == 0 {
		return ""
	}
	return strings.Repeat(" ", 4+(depth-1)*2)
}

func FormatValue(value any, depth int) string {
	switch v := value.(type) {
	case map[string]any:
		if len(v) == 0 {
			return "{}"
		}
		return formatMapValue(v, depth)
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

func formatMapValue(m map[string]any, depth int) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	indent := getIndent(depth)
	innerIndent := getIndent(depth + 1)

	lines := []string{"{"}
	for _, k := range keys {
		valStr := FormatValue(m[k], depth+1)
		lines = append(lines, fmt.Sprintf("%s    %s: %s", innerIndent, k, valStr))
	}
	lines = append(lines, fmt.Sprintf("%s  }", indent))
	return strings.Join(lines, "\n")
}
