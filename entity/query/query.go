package query

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/itchyny/gojq"
)

// Query interface for the query service.
type Query interface {
	Query(opts SelectionOpts) ([]map[string]any, error)
}

// SelectionOpts holds the CLI flag values for stage-1 selection.
type SelectionOpts struct {
	Type string // --type glob pattern (matched with GLOB operator)
	Meta string // --meta substring match in meta_json
	Tag  string // --tag substring match in meta_json
	JQ   string // --jq expression for stage-2 transformation
}

// queryImpl implements Query with a database backend.
type queryImpl struct {
	Database database.Database
}

// Query performs two-stage query: selection from SQLite, then optional jq transformation.
func (q *queryImpl) Query(opts SelectionOpts) ([]map[string]any, error) {
	// Stage 1: Build selection query
	var conditions []string
	var args []any

	if opts.Type != "" {
		conditions = append(conditions, "entity_type GLOB ?")
		args = append(args, opts.Type)
	}
	if opts.Meta != "" {
		conditions = append(conditions, "meta_json LIKE ?")
		args = append(args, "%"+opts.Meta+"%")
	}
	if opts.Tag != "" {
		conditions = append(conditions, "meta_json LIKE ?")
		args = append(args, "%"+opts.Tag+"%")
	}

	query := "SELECT merged_json FROM entities"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := q.Database.QueryRows(query, args...)
	if err != nil {
		return nil, fmt.Errorf("selection query failed: %w", err)
	}

	// Parse merged_json from each row
	var results []map[string]any
	for _, row := range rows {
		jsonStr, ok := row["merged_json"].(string)
		if !ok || jsonStr == "" {
			continue
		}
		var doc map[string]any
		if unmarshalErr := json.Unmarshal([]byte(jsonStr), &doc); unmarshalErr != nil {
			continue
		}
		results = append(results, doc)
	}

	// Stage 2: Apply jq transformation if provided
	if opts.JQ != "" {
		return ApplyJQToAll(opts.JQ, results)
	}

	return results, nil
}

// ApplyJQ applies a jq expression to a single JSON input and returns the results.
func ApplyJQ(expression string, input any) ([]any, error) {
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("invalid jq expression: %w", err)
	}

	var results []any
	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, fmt.Errorf("jq execution error: %w", err)
		}
		results = append(results, v)
	}

	return results, nil
}

// ApplyJQToAll applies a jq expression to each document and collects all results.
func ApplyJQToAll(expression string, docs []map[string]any) ([]map[string]any, error) {
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("invalid jq expression: %w", err)
	}

	var results []map[string]any
	for _, doc := range docs {
		iter := query.Run(doc)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				return nil, fmt.Errorf("jq execution error: %w", err)
			}
			if m, ok := v.(map[string]any); ok {
				results = append(results, m)
			}
		}
	}

	return results, nil
}

// FormatJSON marshals a value to indented JSON string.
func FormatJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
