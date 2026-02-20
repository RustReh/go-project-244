package formatter

import (
	"encoding/json"
)

// jsonNode представляет узел диффа в формате JSON
type jsonNode struct {
	Status   string             `json:"status"`
	Value    any                `json:"value,omitempty"`
	OldValue any                `json:"oldValue,omitempty"`
	NewValue any                `json:"newValue,omitempty"`
	Children map[string]*jsonNode `json:"children,omitempty"` // ИСПРАВЛЕНО: был *jsonNode
}

// convertToJSONNode преобразует внутреннее представление в JSON-формат
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
			// Рекурсивно преобразуем детей в карту
			jsonN.Children = convertToJSONNode(node.Children)
		}

		result[node.Key] = jsonN
	}

	return result
}

// FormatJSON форматирует дифф в JSON
func FormatJSON(nodes []*DiffNode) (string, error) {
	// 1. Конвертируем внутреннюю структуру в JSON-представление
	jsonData := convertToJSONNode(nodes)

	// 2. Маршалим с отступами
	bytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
