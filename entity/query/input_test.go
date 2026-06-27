package query

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDocs(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLen   int
		wantTypes []string
	}{
		{
			name: "yaml multi-doc stream",
			input: "---\n_type: system/cpu@v1\n_body:\n  cores: 4\n" +
				"---\n_type: system/memory@v1\n_body:\n  size: 8\n",
			wantLen:   2,
			wantTypes: []string{"system/cpu@v1", "system/memory@v1"},
		},
		{
			name:      "yaml single doc",
			input:     "_type: system/cpu@v1\n_body:\n  cores: 4\n",
			wantLen:   1,
			wantTypes: []string{"system/cpu@v1"},
		},
		{
			name:      "json array",
			input:     `[{"_type":"a@v1"},{"_type":"b@v1"}]`,
			wantLen:   2,
			wantTypes: []string{"a@v1", "b@v1"},
		},
		{
			name:      "json single object",
			input:     `{"_type":"a@v1"}`,
			wantLen:   1,
			wantTypes: []string{"a@v1"},
		},
		{
			name:      "ndjson stream",
			input:     "{\"_type\":\"a@v1\"}\n{\"_type\":\"b@v1\"}\n",
			wantLen:   2,
			wantTypes: []string{"a@v1", "b@v1"},
		},
		{
			name:    "empty input",
			input:   "   \n  ",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docs, err := LoadDocs(strings.NewReader(tt.input))
			require.NoError(t, err)
			require.Len(t, docs, tt.wantLen)
			for i, want := range tt.wantTypes {
				assert.Equal(t, want, docs[i]["_type"])
			}
		})
	}
}

func TestLoadDocs_NestedMapsAreJQFriendly(t *testing.T) {
	// goccy must decode nested YAML mappings into map[string]any (not
	// map[interface{}]interface{}) so downstream jq/JSON handling works.
	docs, err := LoadDocs(strings.NewReader("_type: a@v1\n_body:\n  nested:\n    key: value\n"))
	require.NoError(t, err)
	require.Len(t, docs, 1)
	body, ok := docs[0]["_body"].(map[string]any)
	require.True(t, ok, "_body should be map[string]any")
	nested, ok := body["nested"].(map[string]any)
	require.True(t, ok, "nested should be map[string]any")
	assert.Equal(t, "value", nested["key"])
}

func TestLoadDocs_InvalidJSON(t *testing.T) {
	_, err := LoadDocs(strings.NewReader(`{"_type": broken`))
	assert.Error(t, err)
}

func TestLoadDocs_JSONArrayNonObjectElement(t *testing.T) {
	_, err := LoadDocs(strings.NewReader(`[{"_type":"a@v1"}, "notanobject"]`))
	assert.Error(t, err)
}
