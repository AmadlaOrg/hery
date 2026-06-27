package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/goccy/go-yaml"
	"github.com/itchyny/gojq"
	"github.com/olekukonko/tablewriter"
)

// HERYResultType is the entity type used when wrapping query results in a HERY
// envelope (--hery).
const HERYResultType = "amadla.org/entity/query/result@v1.0.0"

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

// QueryDocs runs the same two-stage query as Query (selection + optional jq
// transformation) against an in-memory set of entity documents instead of the
// SQLite cache. It is the data source for `hery query --from`.
func QueryDocs(docs []map[string]any, opts SelectionOpts) ([]map[string]any, error) {
	var selected []map[string]any
	for _, doc := range docs {
		if matchDoc(doc, opts) {
			selected = append(selected, doc)
		}
	}

	if opts.JQ != "" {
		return ApplyJQToAll(opts.JQ, selected)
	}

	return selected, nil
}

// matchDoc reports whether a single document satisfies the stage-1 selection
// predicates. The predicates mirror the SQLite path: a type glob, plus
// substring matches against the JSON-encoded _meta block.
func matchDoc(doc map[string]any, opts SelectionOpts) bool {
	if opts.Type != "" {
		typ, _ := doc["_type"].(string)
		if !matchType(opts.Type, typ) {
			return false
		}
	}
	if opts.Meta != "" && !metaContains(doc, opts.Meta) {
		return false
	}
	if opts.Tag != "" && !metaContains(doc, opts.Tag) {
		return false
	}
	return true
}

// matchType matches a glob pattern against an entity's _type. When the pattern
// carries no "@version" suffix, the version is stripped from the type before
// matching so `--type system/cpu` matches `system/cpu@v1.0.0`.
func matchType(pattern, typ string) bool {
	if typ == "" {
		return false
	}
	if !strings.Contains(pattern, "@") {
		if i := strings.Index(typ, "@"); i >= 0 {
			typ = typ[:i]
		}
	}
	return globMatch(pattern, typ)
}

// globMatch performs a case-insensitive, fully-anchored glob match where '*'
// matches any run of characters (including '/') and '?' matches a single
// character — the same wildcard semantics as the SQLite GLOB operator used by
// the cached path, but case-insensitive for ergonomics with type URIs.
func globMatch(pattern, s string) bool {
	var b strings.Builder
	b.WriteString("(?i)^")
	for _, r := range pattern {
		switch r {
		case '*':
			b.WriteString(".*")
		case '?':
			b.WriteString(".")
		default:
			b.WriteString(regexp.QuoteMeta(string(r)))
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile(b.String())
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

// metaContains reports whether the JSON-encoded _meta block of a document
// contains substr, mirroring the SQLite `meta_json LIKE '%substr%'` predicate.
func metaContains(doc map[string]any, substr string) bool {
	meta, ok := doc["_meta"]
	if !ok {
		return false
	}
	b, err := json.Marshal(meta)
	if err != nil {
		return false
	}
	return strings.Contains(string(b), substr)
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

// FormatYAML renders results as a "---"-separated multi-document YAML stream,
// the inverse of the YAML input accepted by LoadDocs.
func FormatYAML(results []map[string]any) (string, error) {
	var buf bytes.Buffer
	for _, doc := range results {
		buf.WriteString("---\n")
		b, err := yaml.Marshal(doc)
		if err != nil {
			return "", err
		}
		buf.Write(b)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// FormatYAMLValue renders a single value (e.g. a HERY envelope) as YAML.
func FormatYAMLValue(v any) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(b), "\n"), nil
}

// FormatTable renders results as a human-readable table with TYPE, NAME and
// META columns. NAME is read from _meta.name when present.
func FormatTable(results []map[string]any) (string, error) {
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.Header("TYPE", "NAME", "META")
	for _, doc := range results {
		typ, _ := doc["_type"].(string)
		name := ""
		meta := ""
		if m, ok := doc["_meta"].(map[string]any); ok {
			if n, ok := m["name"].(string); ok {
				name = n
			}
			if b, err := json.Marshal(m); err == nil {
				meta = string(b)
			}
		}
		if err := table.Append(typ, name, meta); err != nil {
			return "", err
		}
	}
	if err := table.Render(); err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// WrapHERY wraps results in a HERY envelope so the output is itself a valid
// HERY entity ({_type, _body:{items:[...]}}).
func WrapHERY(results []map[string]any) map[string]any {
	items := results
	if items == nil {
		items = []map[string]any{}
	}
	return map[string]any{
		"_type": HERYResultType,
		"_body": map[string]any{
			"items": items,
		},
	}
}
