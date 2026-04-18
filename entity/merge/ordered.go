package merge

import "github.com/goccy/go-yaml"

// OrderedDeepMerge merges parent and child MapSlices preserving parent-first key
// order. Child values override parent values at the same key; child-only keys
// are appended at the end in their original order. Nested MapSlices recurse.
// Arrays and scalars: child wins (no concatenation).
func OrderedDeepMerge(parent, child yaml.MapSlice) yaml.MapSlice {
	parentIdx := indexKeys(parent)
	childIdx := indexKeys(child)

	result := make(yaml.MapSlice, 0, len(parent)+len(child))

	for _, item := range parent {
		key := keyString(item.Key)
		if ci, ok := childIdx[key]; ok {
			merged := mergeValues(item.Value, child[ci].Value)
			result = append(result, yaml.MapItem{Key: item.Key, Value: merged})
		} else {
			result = append(result, item)
		}
	}

	for _, item := range child {
		key := keyString(item.Key)
		if _, ok := parentIdx[key]; !ok {
			result = append(result, item)
		}
	}

	return result
}

func mergeValues(parent, child any) any {
	pMap, pOk := parent.(yaml.MapSlice)
	cMap, cOk := child.(yaml.MapSlice)
	if pOk && cOk {
		return OrderedDeepMerge(pMap, cMap)
	}
	return child
}

func indexKeys(ms yaml.MapSlice) map[string]int {
	idx := make(map[string]int, len(ms))
	for i, item := range ms {
		idx[keyString(item.Key)] = i
	}
	return idx
}

func keyString(k any) string {
	if s, ok := k.(string); ok {
		return s
	}
	return ""
}
