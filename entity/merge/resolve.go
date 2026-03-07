package merge

import "fmt"

// ResolveParentChain resolves _parent references and deep merges entity bodies.
// The lookup function retrieves a parsed entity body by its _self URI.
// Returns the fully merged body after walking the parent chain.
func ResolveParentChain(entityBody map[string]any, lookup func(selfURI string) (map[string]any, error)) (map[string]any, error) {
	parentURI, ok := entityBody["_parent"].(string)
	if !ok || parentURI == "" {
		// No parent — return entity body as-is
		return entityBody, nil
	}

	// Guard against circular references
	seen := map[string]bool{}
	return resolveChain(entityBody, lookup, seen)
}

func resolveChain(child map[string]any, lookup func(string) (map[string]any, error), seen map[string]bool) (map[string]any, error) {
	parentURI, ok := child["_parent"].(string)
	if !ok || parentURI == "" {
		return child, nil
	}

	if seen[parentURI] {
		return nil, fmt.Errorf("circular _parent reference detected: %s", parentURI)
	}
	seen[parentURI] = true

	parentBody, err := lookup(parentURI)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve _parent %s: %w", parentURI, err)
	}

	// Recursively resolve the parent's own _parent first
	resolvedParent, err := resolveChain(parentBody, lookup, seen)
	if err != nil {
		return nil, err
	}

	// Deep merge: parent is base, child overrides
	merged := DeepMerge(resolvedParent, child)

	// Remove _parent from the merged result since it's been resolved
	delete(merged, "_parent")

	return merged, nil
}
