package formatter

import "fmt"

func Format(diff []*DiffNode, format string) (string, error) {
	switch format {
	case "stylish":
		return FormatStylish(diff, 0), nil
	case "plain":
		return FormatPlain(diff, ""), nil
	case "json":
		return FormatJSON(diff)
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
