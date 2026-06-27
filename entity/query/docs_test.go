package query

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func graph() []map[string]any {
	return []map[string]any{
		{"_type": "system@v1.0.0", "_meta": map[string]any{"name": "host"}},
		{"_type": "system/cpu@v1.0.0", "_meta": map[string]any{"name": "cpu", "tier": "fast"}},
		{"_type": "system/memory@v1.0.0", "_meta": map[string]any{"name": "mem"}},
		{"_type": "infrastructure/vm@v1.0.0", "_body": map[string]any{"vcpus": 4}},
	}
}

func typesOf(docs []map[string]any) []string {
	var out []string
	for _, d := range docs {
		out = append(out, d["_type"].(string))
	}
	return out
}

func TestQueryDocs_TypeGlob(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{"exact strips version", "system/cpu", []string{"system/cpu@v1.0.0"}},
		{"star matches across slash", "system/*", []string{"system/cpu@v1.0.0", "system/memory@v1.0.0"}},
		{"star matches parent and children", "system*", []string{"system@v1.0.0", "system/cpu@v1.0.0", "system/memory@v1.0.0"}},
		{"case insensitive", "System/CPU", []string{"system/cpu@v1.0.0"}},
		{"no match", "doesnotexist", nil},
		{"version in pattern", "system/cpu@v1.0.0", []string{"system/cpu@v1.0.0"}},
		{"version mismatch", "system/cpu@v2.0.0", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryDocs(graph(), SelectionOpts{Type: tt.pattern})
			require.NoError(t, err)
			assert.Equal(t, tt.want, typesOf(got))
		})
	}
}

func TestQueryDocs_NoFilterReturnsAll(t *testing.T) {
	got, err := QueryDocs(graph(), SelectionOpts{})
	require.NoError(t, err)
	assert.Len(t, got, 4)
}

func TestQueryDocs_MetaSubstring(t *testing.T) {
	got, err := QueryDocs(graph(), SelectionOpts{Meta: "fast"})
	require.NoError(t, err)
	assert.Equal(t, []string{"system/cpu@v1.0.0"}, typesOf(got))
}

func TestQueryDocs_TagSubstring(t *testing.T) {
	got, err := QueryDocs(graph(), SelectionOpts{Tag: "host"})
	require.NoError(t, err)
	assert.Equal(t, []string{"system@v1.0.0"}, typesOf(got))
}

func TestQueryDocs_TypeAndMeta(t *testing.T) {
	got, err := QueryDocs(graph(), SelectionOpts{Type: "system/*", Meta: "fast"})
	require.NoError(t, err)
	assert.Equal(t, []string{"system/cpu@v1.0.0"}, typesOf(got))
}

func TestQueryDocs_WithJQ(t *testing.T) {
	got, err := QueryDocs(graph(), SelectionOpts{Type: "infrastructure/vm", JQ: "._body"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.EqualValues(t, 4, got[0]["vcpus"])
}

func TestFormatYAML_RoundTrips(t *testing.T) {
	out, err := FormatYAML(graph())
	require.NoError(t, err)
	assert.Contains(t, out, "---")
	docs, err := LoadDocs(strings.NewReader(out))
	require.NoError(t, err)
	assert.Len(t, docs, 4)
}

func TestFormatTable(t *testing.T) {
	out, err := FormatTable(graph())
	require.NoError(t, err)
	assert.Contains(t, out, "TYPE")
	assert.Contains(t, out, "system/cpu@v1.0.0")
	assert.Contains(t, out, "cpu")
}

func TestWrapHERY(t *testing.T) {
	env := WrapHERY(graph())
	assert.Equal(t, HERYResultType, env["_type"])
	body := env["_body"].(map[string]any)
	items := body["items"].([]map[string]any)
	assert.Len(t, items, 4)
}

func TestWrapHERY_EmptyIsNonNilSlice(t *testing.T) {
	env := WrapHERY(nil)
	items := env["_body"].(map[string]any)["items"].([]map[string]any)
	assert.NotNil(t, items)
	assert.Len(t, items, 0)
}
