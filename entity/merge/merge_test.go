package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepMerge_ScalarOverride(t *testing.T) {
	parent := map[string]any{"port": 80, "protocol": "tcp"}
	child := map[string]any{"port": 443}

	result := DeepMerge(parent, child)

	assert.Equal(t, 443, result["port"])
	assert.Equal(t, "tcp", result["protocol"])
}

func TestDeepMerge_NestedObjects(t *testing.T) {
	parent := map[string]any{
		"database": map[string]any{
			"engine": "postgres",
			"port":   5432,
			"pool":   map[string]any{"min": 5, "max": 20},
		},
	}
	child := map[string]any{
		"database": map[string]any{
			"port": 5433,
			"pool": map[string]any{"max": 50},
		},
	}

	result := DeepMerge(parent, child)

	db := result["database"].(map[string]any)
	assert.Equal(t, "postgres", db["engine"])
	assert.Equal(t, 5433, db["port"])
	pool := db["pool"].(map[string]any)
	assert.Equal(t, 5, pool["min"])
	assert.Equal(t, 50, pool["max"])
}

func TestDeepMerge_ArrayReplace(t *testing.T) {
	parent := map[string]any{
		"tags": []any{"a", "b", "c"},
	}
	child := map[string]any{
		"tags": []any{"x"},
	}

	result := DeepMerge(parent, child)

	assert.Equal(t, []any{"x"}, result["tags"])
}

func TestDeepMerge_ChildOnlyKeys(t *testing.T) {
	parent := map[string]any{"a": 1}
	child := map[string]any{"b": 2}

	result := DeepMerge(parent, child)

	assert.Equal(t, 1, result["a"])
	assert.Equal(t, 2, result["b"])
}

func TestDeepMerge_EmptyParent(t *testing.T) {
	parent := map[string]any{}
	child := map[string]any{"port": 443}

	result := DeepMerge(parent, child)

	assert.Equal(t, 443, result["port"])
}

func TestDeepMerge_EmptyChild(t *testing.T) {
	parent := map[string]any{"port": 80}
	child := map[string]any{}

	result := DeepMerge(parent, child)

	assert.Equal(t, 80, result["port"])
}

func TestDeepMerge_TypeMismatch(t *testing.T) {
	parent := map[string]any{
		"config": map[string]any{"key": "value"},
	}
	child := map[string]any{
		"config": "simple-string",
	}

	result := DeepMerge(parent, child)

	assert.Equal(t, "simple-string", result["config"])
}

func TestDeepMerge_NilValues(t *testing.T) {
	parent := map[string]any{"key": "value"}
	child := map[string]any{"key": nil}

	result := DeepMerge(parent, child)

	assert.Nil(t, result["key"])
}
