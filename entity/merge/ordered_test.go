package merge

import (
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestOrderedDeepMerge_ParentOnlyKeys(t *testing.T) {
	parent := yaml.MapSlice{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
	}
	result := OrderedDeepMerge(parent, nil)
	assert.Equal(t, parent, result)
}

func TestOrderedDeepMerge_ChildAppendedAtEnd(t *testing.T) {
	parent := yaml.MapSlice{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
	}
	child := yaml.MapSlice{
		{Key: "c", Value: 3},
	}
	result := OrderedDeepMerge(parent, child)
	expected := yaml.MapSlice{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}
	assert.Equal(t, expected, result)
}

func TestOrderedDeepMerge_ChildOverridesKeepsParentPosition(t *testing.T) {
	parent := yaml.MapSlice{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}
	child := yaml.MapSlice{
		{Key: "b", Value: 99},
	}
	result := OrderedDeepMerge(parent, child)
	expected := yaml.MapSlice{
		{Key: "a", Value: 1},
		{Key: "b", Value: 99},
		{Key: "c", Value: 3},
	}
	assert.Equal(t, expected, result)
}

func TestOrderedDeepMerge_NestedMapSliceRecurses(t *testing.T) {
	parent := yaml.MapSlice{
		{Key: "ssh", Value: yaml.MapSlice{
			{Key: "user", Value: "root"},
			{Key: "port", Value: 22},
		}},
	}
	child := yaml.MapSlice{
		{Key: "ssh", Value: yaml.MapSlice{
			{Key: "port", Value: 2222},
		}},
	}
	result := OrderedDeepMerge(parent, child)
	expected := yaml.MapSlice{
		{Key: "ssh", Value: yaml.MapSlice{
			{Key: "user", Value: "root"},
			{Key: "port", Value: 2222},
		}},
	}
	assert.Equal(t, expected, result)
}

func TestOrderedDeepMerge_ChildReplacesArray(t *testing.T) {
	parent := yaml.MapSlice{
		{Key: "arch", Value: []any{"x86_64"}},
	}
	child := yaml.MapSlice{
		{Key: "arch", Value: []any{"arm64", "riscv64"}},
	}
	result := OrderedDeepMerge(parent, child)
	expected := yaml.MapSlice{
		{Key: "arch", Value: []any{"arm64", "riscv64"}},
	}
	assert.Equal(t, expected, result)
}
