package code

import (
	"fmt"

	parser "code/parser"
	formatter "code/formatter"
)


func GenDiff(path1, path2, format string) (string, error) {
    data1, err := parser.ParseFile(path1)
    if err != nil {
        return "", err
    }

    data2, err := parser.ParseFile(path2)
    if err != nil {
        return "", err
    }

    diff := formatter.BuildDiff(data1, data2)

    switch format {
    case "stylish":
        return formatter.FormatStylish(diff, 0), nil
    case "plain":
        return formatter.FormatPlain(diff, ""), nil
    case "json":
        return formatter.FormatJSON(diff)
    default:
        return "", fmt.Errorf("unsupported format: %s", format)
    }
}

