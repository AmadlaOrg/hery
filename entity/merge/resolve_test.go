package merge

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveExtendsChain_NoExtends(t *testing.T) {
	body := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_body": map[string]any{"name": "test"},
	}

	result, err := ResolveExtendsChain(body, nil)
	assert.NoError(t, err)
	assert.Equal(t, body, result)
}

func TestResolveExtendsChain_SingleExtends(t *testing.T) {
	baseBody := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_body": map[string]any{
			"name":  "parent-name",
			"color": "blue",
		},
	}

	childBody := map[string]any{
		"_type":   "example.com/Entity@v1.0.0",
		"_extends": "parent-entity",
		"_body": map[string]any{
			"name": "child-name",
		},
	}

	lookup := func(selfURI string) (map[string]any, error) {
		if selfURI == "parent-entity" {
			return baseBody, nil
		}
		return nil, fmt.Errorf("not found: %s", selfURI)
	}

	result, err := ResolveExtendsChain(childBody, lookup)
	assert.NoError(t, err)

	// _extends should be removed
	_, hasExtends := result["_extends"]
	assert.False(t, hasExtends)

	// child name should override base
	body := result["_body"].(map[string]any)
	assert.Equal(t, "child-name", body["name"])
	// base color should be preserved
	assert.Equal(t, "blue", body["color"])
}

func TestResolveExtendsChain_GrandExtendsChain(t *testing.T) {
	grandparent := map[string]any{
		"_body": map[string]any{
			"a": "from-grandparent",
			"b": "from-grandparent",
			"c": "from-grandparent",
		},
	}

	parent := map[string]any{
		"_extends": "grandparent",
		"_body": map[string]any{
			"b": "from-parent",
		},
	}

	child := map[string]any{
		"_extends": "parent",
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

	result, err := ResolveExtendsChain(child, lookup)
	assert.NoError(t, err)

	body := result["_body"].(map[string]any)
	assert.Equal(t, "from-grandparent", body["a"])
	assert.Equal(t, "from-parent", body["b"])
	assert.Equal(t, "from-child", body["c"])
}

func TestResolveExtendsChain_CircularReference(t *testing.T) {
	lookup := func(selfURI string) (map[string]any, error) {
		// Both point to each other
		return map[string]any{
			"_extends": "entity-a",
		}, nil
	}

	child := map[string]any{
		"_extends": "entity-a",
	}

	_, err := ResolveExtendsChain(child, lookup)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular _extends reference")
}

func TestResolveExtendsChain_LookupError(t *testing.T) {
	child := map[string]any{
		"_extends": "missing-entity",
	}

	lookup := func(selfURI string) (map[string]any, error) {
		return nil, fmt.Errorf("entity not found")
	}

	_, err := ResolveExtendsChain(child, lookup)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to resolve _extends")
}

func TestResolveExtendsChain_MetaMerge(t *testing.T) {
	baseBody := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_meta": map[string]any{
			"name":     "parent-name",
			"category": "Application",
			"tags":     []any{"base", "production"},
		},
		"_body": map[string]any{
			"port": 80,
		},
	}

	childBody := map[string]any{
		"_type":   "example.com/Entity@v1.0.0",
		"_extends": "parent-entity",
		"_meta": map[string]any{
			"name": "child-name",
		},
		"_body": map[string]any{
			"port": 443,
		},
	}

	lookup := func(uri string) (map[string]any, error) {
		if uri == "parent-entity" {
			return baseBody, nil
		}
		return nil, fmt.Errorf("not found: %s", uri)
	}

	result, err := ResolveExtendsChain(childBody, lookup)
	assert.NoError(t, err)

	// _meta should be deep merged: child name overrides, base category preserved
	meta := result["_meta"].(map[string]any)
	assert.Equal(t, "child-name", meta["name"])
	assert.Equal(t, "Application", meta["category"])
	// tags is an array — child didn't set it, so base value preserved
	assert.Equal(t, []any{"base", "production"}, meta["tags"])
}

func TestResolveExtendsChain_RequiresMerge(t *testing.T) {
	baseBody := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_requires": []any{
			"github.com/AmadlaOrg/Entities/Application/DB/RDBMS@v1.0.0",
			"github.com/AmadlaOrg/Entities/Network@v1.0.0",
		},
		"_body": map[string]any{
			"name": "parent",
		},
	}

	childBody := map[string]any{
		"_type":   "example.com/Entity@v1.0.0",
		"_extends": "parent-entity",
		"_requires": []any{
			"github.com/AmadlaOrg/Entities/Application/DB/RDBMS@v2.0.0",
		},
		"_body": map[string]any{
			"name": "child",
		},
	}

	lookup := func(uri string) (map[string]any, error) {
		if uri == "parent-entity" {
			return baseBody, nil
		}
		return nil, fmt.Errorf("not found: %s", uri)
	}

	result, err := ResolveExtendsChain(childBody, lookup)
	assert.NoError(t, err)

	// _requires is an array — child replaces base entirely (arrays replace, not merge)
	requires := result["_requires"].([]any)
	assert.Len(t, requires, 1)
	assert.Equal(t, "github.com/AmadlaOrg/Entities/Application/DB/RDBMS@v2.0.0", requires[0])
}

func TestResolveExtendsChain_RequiresInheritedFromExtends(t *testing.T) {
	baseBody := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_requires": []any{
			"github.com/AmadlaOrg/Entities/Network@v1.0.0",
		},
		"_body": map[string]any{
			"name": "parent",
		},
	}

	childBody := map[string]any{
		"_type":   "example.com/Entity@v1.0.0",
		"_extends": "parent-entity",
		"_body": map[string]any{
			"name": "child",
		},
	}

	lookup := func(uri string) (map[string]any, error) {
		if uri == "parent-entity" {
			return baseBody, nil
		}
		return nil, fmt.Errorf("not found: %s", uri)
	}

	result, err := ResolveExtendsChain(childBody, lookup)
	assert.NoError(t, err)

	// Child didn't set _requires, so base's _requires is inherited
	requires := result["_requires"].([]any)
	assert.Len(t, requires, 1)
	assert.Equal(t, "github.com/AmadlaOrg/Entities/Network@v1.0.0", requires[0])
}

func TestResolveExtendsChain_MetaAndRequiresFullMerge(t *testing.T) {
	baseBody := map[string]any{
		"_type": "example.com/Entity@v1.0.0",
		"_meta": map[string]any{
			"name":        "base-webserver",
			"description": "Base web server config",
			"tags":        []any{"web", "base"},
		},
		"_requires": []any{
			"github.com/AmadlaOrg/Entities/Network@v1.0.0",
		},
		"_body": map[string]any{
			"port":     80,
			"protocol": "tcp",
		},
	}

	childBody := map[string]any{
		"_type":   "example.com/Entity@v1.0.0",
		"_extends": "parent-entity",
		"_meta": map[string]any{
			"name": "my-webserver",
			"tags": []any{"production"},
		},
		"_requires": []any{
			"github.com/AmadlaOrg/Entities/Application/DB/RDBMS@v1.0.0",
		},
		"_body": map[string]any{
			"port": 443,
		},
	}

	lookup := func(uri string) (map[string]any, error) {
		if uri == "parent-entity" {
			return baseBody, nil
		}
		return nil, fmt.Errorf("not found: %s", uri)
	}

	result, err := ResolveExtendsChain(childBody, lookup)
	assert.NoError(t, err)

	// _meta: deep merged (objects recursive, arrays replace)
	meta := result["_meta"].(map[string]any)
	assert.Equal(t, "my-webserver", meta["name"])
	assert.Equal(t, "Base web server config", meta["description"]) // inherited from base
	assert.Equal(t, []any{"production"}, meta["tags"])             // array replaced by child

	// _requires: array replaced by child
	requires := result["_requires"].([]any)
	assert.Len(t, requires, 1)
	assert.Equal(t, "github.com/AmadlaOrg/Entities/Application/DB/RDBMS@v1.0.0", requires[0])

	// _body: deep merged
	body := result["_body"].(map[string]any)
	assert.Equal(t, 443, body["port"])
	assert.Equal(t, "tcp", body["protocol"]) // inherited from base
}
