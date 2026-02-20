package formatter

import (
    "encoding/json"
)

type jsonNode struct {
    Key      string `json:"key"`
    Status   string `json:"status"`
    Value    any    `json:"value,omitempty"`
    OldValue any    `json:"oldValue,omitempty"`
    NewValue any    `json:"newValue,omitempty"`
    Children any    `json:"children,omitempty"`
}

func FormatJSON(nodes []*DiffNode) (string, error) {
    j := convertToJSONNodes(nodes)
    payload := map[string]any{
        "diff": j,
    }

    bytes, err := json.MarshalIndent(payload, "", "    ")
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func convertToJSONNodes(nodes []*DiffNode) []*jsonNode {
    result := make([]*jsonNode, 0, len(nodes))
    for _, node := range nodes {
        jsonN := &jsonNode{
            Key:    node.Key,
            Status: node.Type,
        }

        switch node.Type {
        case "added":
            jsonN.NewValue = convertValue(node.Value)
        case "removed":
            jsonN.OldValue = convertValue(node.Value)
        case "unchanged":
            jsonN.OldValue = convertValue(node.Value)
        case "updated":
            jsonN.OldValue = convertValue(node.OldVal)
            jsonN.NewValue = convertValue(node.NewVal)
        case "nested":
            childrenSlice := convertToJSONNodes(node.Children)
            if len(childrenSlice) > 0 {
                jsonN.Children = childrenSlice
            }
        }
        result = append(result, jsonN)
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
