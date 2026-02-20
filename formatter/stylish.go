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
	contentIndent := strings.Repeat(" ", (depth+1)*indentSize)
	bracketIndent := strings.Repeat(" ", depth*indentSize)

	lines := []string{"{"}

	for _, node := range nodes {
		switch node.Type {
		case "added":
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", contentIndent, node.Key, FormatValue(node.Value, depth+1)))
		case "removed":
			lines = append(lines, fmt.Sprintf("%s- %s: %s", contentIndent, node.Key, FormatValue(node.Value, depth+1)))
		case "unchanged":
			lines = append(lines, fmt.Sprintf("%s  %s: %s", contentIndent, node.Key, FormatValue(node.Value, depth+1)))
		case "updated":
			lines = append(lines, fmt.Sprintf("%s- %s: %s", contentIndent, node.Key, FormatValue(node.OldVal, depth+1)))
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", contentIndent, node.Key, FormatValue(node.NewVal, depth+1)))
		case "nested":
			lines = append(lines, fmt.Sprintf("%s%s: {", contentIndent, node.Key))
			lines = append(lines, renderNodesRaw(node.Children, depth+1))
			lines = append(lines, fmt.Sprintf("%s}", contentIndent))
		}
	}

	lines = append(lines, fmt.Sprintf("%s}", bracketIndent))
	return strings.Join(lines, "\n")
}

func renderNodesRaw(nodes []*DiffNode, depth int) string {
	if len(nodes) == 0 {
		return ""
	}

	contentIndent := strings.Repeat(" ", (depth+1)*indentSize)
	var lines []string

	for _, node := range nodes {
		switch node.Type {
		case "added":
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", contentIndent, node.Key, FormatValue(node.Value, depth+1)))
		case "removed":
			lines = append(lines, fmt.Sprintf("%s- %s: %s", contentIndent, node.Key, FormatValue(node.Value, depth+1)))
		case "unchanged":
			lines = append(lines, fmt.Sprintf("%s  %s: %s", contentIndent, node.Key, FormatValue(node.Value, depth+1)))
		case "updated":
			lines = append(lines, fmt.Sprintf("%s- %s: %s", contentIndent, node.Key, FormatValue(node.OldVal, depth+1)))
			lines = append(lines, fmt.Sprintf("%s+ %s: %s", contentIndent, node.Key, FormatValue(node.NewVal, depth+1)))
		case "nested":
			lines = append(lines, fmt.Sprintf("%s%s: {", contentIndent, node.Key))
			lines = append(lines, renderNodesRaw(node.Children, depth+1))
			lines = append(lines, fmt.Sprintf("%s}", contentIndent))
		}
	}
	return strings.Join(lines, "\n")
}

func FormatValue(value any, depth int) string {
	switch v := value.(type) {
	case map[string]any:
		if len(v) == 0 {
			return "{}"
		}

		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		keyIndent := strings.Repeat(" ", depth*indentSize)
		bracketIndent := strings.Repeat(" ", (depth-1)*indentSize)

		lines := []string{"{"}
		for _, k := range keys {
			valStr := FormatValue(v[k], depth+1)
			lines = append(lines, fmt.Sprintf("%s%s: %s", keyIndent, k, valStr))
		}
		lines = append(lines, fmt.Sprintf("%s}", bracketIndent))
		return strings.Join(lines, "\n")

	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "<nil>"
	default:
		return fmt.Sprintf("%v", v)
	}
}
