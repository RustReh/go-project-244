package formatter

import "encoding/json"

func FormatJSON(nodes []*DiffNode) (string, error) {
	bytes, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
