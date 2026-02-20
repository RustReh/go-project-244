package formatter

import (
	"reflect"
	"sort"
)

type DiffNode struct {
	Type     string      `json:"type"`
	Key      string      `json:"key"`
	Value    any         `json:"value,omitempty"`
	OldVal   any         `json:"old_value,omitempty"`
	NewVal   any         `json:"new_value,omitempty"`
	Children []*DiffNode `json:"children,omitempty"`
}

func BuildDiff(a, b map[string]any) []*DiffNode {
	keys := collectKeys(a, b)
	sort.Strings(keys)

	var diff []*DiffNode

	for _, key := range keys {
		valA, okA := a[key]
		valB, okB := b[key]

		switch {
		case !okA:
			diff = append(diff, &DiffNode{Type: "added", Key: key, Value: valB})

		case !okB:
			diff = append(diff, &DiffNode{Type: "removed", Key: key, Value: valA})

		case reflect.DeepEqual(valA, valB):
			diff = append(diff, &DiffNode{Type: "unchanged", Key: key, Value: valA})

		case isMap(valA) && isMap(valB):
			children := BuildDiff(valA.(map[string]any), valB.(map[string]any))
			diff = append(diff, &DiffNode{Type: "nested", Key: key, Children: children})

		default:
			diff = append(diff, &DiffNode{
				Type:   "updated",
				Key:    key,
				OldVal: valA,
				NewVal: valB,
			})
		}
	}

	return diff
}

func collectKeys(a, b map[string]any) []string {
	keys := make(map[string]struct{})

	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}

	return result
}

func isMap(val any) bool {
	_, ok := val.(map[string]any)
	return ok
}
