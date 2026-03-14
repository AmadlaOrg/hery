package merge

import "fmt"

// ResolveExtendsChain resolves _extends references and deep merges entity bodies.
// The lookup function retrieves a parsed entity body by its URI.
// Returns the fully merged body after walking the extends chain.
func ResolveExtendsChain(entityBody map[string]any, lookup func(extendsURI string) (map[string]any, error)) (map[string]any, error) {
	extendsURI, ok := entityBody["_extends"].(string)
	if !ok || extendsURI == "" {
		// No extends — return entity body as-is
		return entityBody, nil
	}

	// Guard against circular references
	seen := map[string]bool{}
	return resolveChain(entityBody, lookup, seen)
}

func resolveChain(child map[string]any, lookup func(string) (map[string]any, error), seen map[string]bool) (map[string]any, error) {
	extendsURI, ok := child["_extends"].(string)
	if !ok || extendsURI == "" {
		return child, nil
	}

	if seen[extendsURI] {
		return nil, fmt.Errorf("circular _extends reference detected: %s", extendsURI)
	}
	seen[extendsURI] = true

	extendsBody, err := lookup(extendsURI)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve _extends %s: %w", extendsURI, err)
	}

	// Recursively resolve the extended entity's own _extends first
	resolvedExtends, err := resolveChain(extendsBody, lookup, seen)
	if err != nil {
		return nil, err
	}

	// Deep merge: extended entity is base, child overrides
	merged := DeepMerge(resolvedExtends, child)

	// Remove _extends from the merged result since it's been resolved
	delete(merged, "_extends")

	return merged, nil
}
