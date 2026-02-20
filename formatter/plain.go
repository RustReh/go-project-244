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
			lines = append(lines,
				fmt.Sprintf("Property '%s' was added with value: %v",
					currentPath,
					node.Value,
				),
			)

		case "removed":
			lines = append(lines,
				fmt.Sprintf("Property '%s' was removed", currentPath),
			)

		case "updated":
			lines = append(lines,
				fmt.Sprintf("Property '%s' was updated. From %v to %v",
					currentPath,
					node.OldVal,
					node.NewVal,
				),
			)

		case "nested":
			nested := FormatPlain(node.Children, currentPath)
			if nested != "" {
				lines = append(lines, nested)
			}
		}
	}

	return strings.Join(lines, "\n")
}
