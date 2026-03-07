package query

import (
	"testing"

	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuery_NoFilters(t *testing.T) {
	mockDb := &database.MockCacheDatabase{}
	mockDb.EXPECT().QueryRows("SELECT merged_json FROM entities").Return([]map[string]any{
		{"merged_json": `{"_type":"example.com/App@v1.0.0","_body":{"name":"test"}}`},
	}, nil)

	q := &SQuery{Database: mockDb}
	results, err := q.Query(SelectionOpts{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "test", results[0]["_body"].(map[string]any)["name"])
}

func TestQuery_WithTypeFilter(t *testing.T) {
	mockDb := &database.MockCacheDatabase{}
	mockDb.EXPECT().QueryRows(
		"SELECT merged_json FROM entities WHERE entity_type GLOB ?",
		mock.Anything,
	).Return([]map[string]any{
		{"merged_json": `{"_type":"example.com/App@v1.0.0","_body":{"name":"matched"}}`},
	}, nil)

	q := &SQuery{Database: mockDb}
	results, err := q.Query(SelectionOpts{Type: "example.com/App*"})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestQuery_WithJQ(t *testing.T) {
	mockDb := &database.MockCacheDatabase{}
	mockDb.EXPECT().QueryRows("SELECT merged_json FROM entities").Return([]map[string]any{
		{"merged_json": `{"_type":"example.com/App@v1.0.0","_body":{"name":"test","color":"blue"}}`},
	}, nil)

	q := &SQuery{Database: mockDb}
	results, err := q.Query(SelectionOpts{JQ: "._body"})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "test", results[0]["name"])
	assert.Equal(t, "blue", results[0]["color"])
}

func TestQuery_EmptyResults(t *testing.T) {
	mockDb := &database.MockCacheDatabase{}
	mockDb.EXPECT().QueryRows("SELECT merged_json FROM entities WHERE entity_self = ?", mock.Anything).
		Return([]map[string]any{}, nil)

	q := &SQuery{Database: mockDb}
	results, err := q.Query(SelectionOpts{Self: "nonexistent"})
	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestApplyJQ(t *testing.T) {
	input := map[string]any{
		"name":  "test",
		"count": 42,
	}
	results, err := ApplyJQ(".name", input)
	assert.NoError(t, err)
	assert.Equal(t, []any{"test"}, results)
}

func TestApplyJQ_InvalidExpression(t *testing.T) {
	_, err := ApplyJQ("invalid[[[", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid jq expression")
}
