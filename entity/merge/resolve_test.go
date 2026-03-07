package merge

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveParentChain_NoParent(t *testing.T) {
	body := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_body": map[string]any{"name": "test"},
	}

	result, err := ResolveParentChain(body, nil)
	assert.NoError(t, err)
	assert.Equal(t, body, result)
}

func TestResolveParentChain_SingleParent(t *testing.T) {
	parentBody := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_body": map[string]any{
			"name":  "parent-name",
			"color": "blue",
		},
	}

	childBody := map[string]any{
		"_type":   "example.com/Entity@v1.0.0",
		"_parent": "parent-entity",
		"_body": map[string]any{
			"name": "child-name",
		},
	}

	lookup := func(selfURI string) (map[string]any, error) {
		if selfURI == "parent-entity" {
			return parentBody, nil
		}
		return nil, fmt.Errorf("not found: %s", selfURI)
	}

	result, err := ResolveParentChain(childBody, lookup)
	assert.NoError(t, err)

	// _parent should be removed
	_, hasParent := result["_parent"]
	assert.False(t, hasParent)

	// child name should override parent
	body := result["_body"].(map[string]any)
	assert.Equal(t, "child-name", body["name"])
	// parent color should be preserved
	assert.Equal(t, "blue", body["color"])
}

func TestResolveParentChain_GrandparentChain(t *testing.T) {
	grandparent := map[string]any{
		"_body": map[string]any{
			"a": "from-grandparent",
			"b": "from-grandparent",
			"c": "from-grandparent",
		},
	}

	parent := map[string]any{
		"_parent": "grandparent",
		"_body": map[string]any{
			"b": "from-parent",
		},
	}

	child := map[string]any{
		"_parent": "parent",
		"_body": map[string]any{
			"c": "from-child",
		},
	}

	lookup := func(selfURI string) (map[string]any, error) {
		switch selfURI {
		case "grandparent":
			return grandparent, nil
		case "parent":
			return parent, nil
		}
		return nil, fmt.Errorf("not found: %s", selfURI)
	}

	result, err := ResolveParentChain(child, lookup)
	assert.NoError(t, err)

	body := result["_body"].(map[string]any)
	assert.Equal(t, "from-grandparent", body["a"])
	assert.Equal(t, "from-parent", body["b"])
	assert.Equal(t, "from-child", body["c"])
}

func TestResolveParentChain_CircularReference(t *testing.T) {
	lookup := func(selfURI string) (map[string]any, error) {
		// Both point to each other
		return map[string]any{
			"_parent": "entity-a",
		}, nil
	}

	child := map[string]any{
		"_parent": "entity-a",
	}

	_, err := ResolveParentChain(child, lookup)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular _parent reference")
}

func TestResolveParentChain_LookupError(t *testing.T) {
	child := map[string]any{
		"_parent": "missing-parent",
	}

	lookup := func(selfURI string) (map[string]any, error) {
		return nil, fmt.Errorf("entity not found")
	}

	_, err := ResolveParentChain(child, lookup)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to resolve _parent")
}
