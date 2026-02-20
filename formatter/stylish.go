package formatter

import (
    "fmt"
    "sort"
    "strings"
)

const indentSize = 2

func FormatStylish(nodes []*DiffNode, depth int) string {
    if len(nodes) == 0 {
        return "{}"
    }

    lines := []string{"{"}
    for _, node := range nodes {
        lines = append(lines, formatNode(node, depth))
    }

    closingIndent := strings.Repeat(" ", depth*indentSize)
    lines = append(lines, fmt.Sprintf("%s}", closingIndent))
    return strings.Join(lines, "\n")
}

func formatNode(node *DiffNode, depth int) string {
    propIndent := strings.Repeat(" ", (depth+1)*indentSize)
    markerIndent := strings.Repeat(" ", (depth+1)*indentSize-2)

    switch node.Type {
    case "added":
        return fmt.Sprintf("%s+ %s: %s", markerIndent, node.Key, FormatValue(node.Value, depth+1))
    case "removed":
        return fmt.Sprintf("%s- %s: %s", markerIndent, node.Key, FormatValue(node.Value, depth+1))
    case "unchanged":
        return fmt.Sprintf("%s%s: %s", propIndent, node.Key, FormatValue(node.Value, depth+1))
    case "updated":
        line1 := fmt.Sprintf("%s- %s: %s", markerIndent, node.Key, FormatValue(node.OldVal, depth+1))
        line2 := fmt.Sprintf("%s+ %s: %s", markerIndent, node.Key, FormatValue(node.NewVal, depth+1))
        return line1 + "\n" + line2
    case "nested":
        nestedBlock := FormatStylish(node.Children, depth+1)
        return fmt.Sprintf("%s%s: %s", propIndent, node.Key, nestedBlock)
    }
    return ""
}

func FormatValue(value any, depth int) string {
    switch v := value.(type) {
    case map[string]any:
        return formatMap(v, depth)
    case string:
        return v  // Без кавычек для stylish!
    case bool:
        return fmt.Sprintf("%t", v)
    case nil:
        return ""
    default:
        return fmt.Sprintf("%v", v)
    }
}

func formatMap(m map[string]any, depth int) string {
    if len(m) == 0 {
        return "{}"
    }

    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    lines := []string{"{"}
    propIndent := strings.Repeat(" ", (depth+1)*indentSize)

    for _, k := range keys {
        lines = append(lines, fmt.Sprintf("%s%s: %s", propIndent, k, FormatValue(m[k], depth+1)))
    }

    closingIndent := strings.Repeat(" ", depth*indentSize)
    lines = append(lines, fmt.Sprintf("%s}", closingIndent))

    return strings.Join(lines, "\n")
}
