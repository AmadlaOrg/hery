package merge

// DeepMerge merges parent and child maps using HERY deep merge semantics:
//   - Objects: merge recursively, child overrides parent for same key
//   - Arrays: child replaces parent entirely
//   - Scalars: child wins
//   - Missing keys in child: parent value preserved
func DeepMerge(parent, child map[string]any) map[string]any {
	result := make(map[string]any, len(parent)+len(child))

	// Start with all parent keys
	for k, v := range parent {
		result[k] = v
	}

	// Apply child overrides
	for k, childVal := range child {
		parentVal, exists := result[k]
		if !exists {
			result[k] = childVal
			continue
		}

		parentMap, parentIsMap := parentVal.(map[string]any)
		childMap, childIsMap := childVal.(map[string]any)

		if parentIsMap && childIsMap {
			// Both are objects: recurse
			result[k] = DeepMerge(parentMap, childMap)
		} else {
			// Arrays, scalars, or type mismatch: child wins
			result[k] = childVal
		}
	}

	return result
}
