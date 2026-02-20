package formatter

import (
	"encoding/json"
)

type jsonNode struct {
	Status   string      `json:"status"`
	Value    any `json:"value,omitempty"`
	OldValue any `json:"oldValue,omitempty"`
	NewValue any `json:"newValue,omitempty"`
	Children any `json:"children,omitempty"`
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
			jsonN.Value = convertValue(node.Value)
		case "removed":
			jsonN.Status = "removed"
			jsonN.Value = convertValue(node.Value)
		case "unchanged":
			jsonN.Status = "unchanged"
			jsonN.Value = convertValue(node.Value)
		case "updated":
			jsonN.Status = "updated"
			jsonN.OldValue = convertValue(node.OldVal)
			jsonN.NewValue = convertValue(node.NewVal)
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

func convertValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return convertMap(v)
	case nil:
		return nil
	default:
		return v
	}
}

func convertMap(m map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		switch val := v.(type) {
		case map[string]any:
			result[k] = convertMap(val)
		default:
			result[k] = val
		}
	}
	return result
}
