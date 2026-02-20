package formatter

import (
	"encoding/json"
)

type jsonNode struct {
	Status   string      `json:"status"`
	Value    any         `json:"value,omitempty"`
	OldValue any         `json:"oldValue,omitempty"`
	NewValue any         `json:"newValue,omitempty"`
	Children interface{} `json:"children,omitempty"`
}

func FormatJSON(nodes []*DiffNode) (string, error) {
	jsonData := convertToJSONNode(nodes)
	bytes, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func convertToJSONNode(nodes []*DiffNode) map[string]*jsonNode {
	result := make(map[string]*jsonNode)
	for _, node := range nodes {
		jsonN := &jsonNode{}
		switch node.Type {
		case "added":
			jsonN.Status = "added"
			jsonN.Value = node.Value
		case "removed":
			jsonN.Status = "removed"
			jsonN.Value = node.Value
		case "unchanged":
			jsonN.Status = "unchanged"
			jsonN.Value = node.Value
		case "updated":
			jsonN.Status = "updated"
			jsonN.OldValue = node.OldVal
			jsonN.NewValue = node.NewVal
		case "nested":
			jsonN.Status = "nested"
			childrenMap := convertToJSONNode(node.Children)
			if len(childrenMap) > 0 {
				jsonN.Children = childrenMap
			}
		}
		result[node.Key] = jsonN
	}
	return result
}
