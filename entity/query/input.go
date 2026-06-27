package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

// LoadDocs reads entity documents from r, auto-detecting the input format by
// sniffing the first non-whitespace byte:
//
//   - '{' or '['  -> JSON (single object, array of objects, or NDJSON stream)
//   - anything else -> YAML (single document or a "---"-separated multi-doc stream)
//
// Each returned map is one entity document. The entire input is read into
// memory; graphs produced by `hery compose --dir` are small.
func LoadDocs(r io.Reader) ([]map[string]any, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, nil
	}

	switch trimmed[0] {
	case '{', '[':
		return loadJSON(trimmed)
	default:
		return loadYAML(data)
	}
}

// loadJSON decodes a JSON array of objects, a single object, or an NDJSON
// stream of objects. The json.Decoder loop transparently handles all three:
// an array decodes as one []any value, while concatenated objects decode one
// at a time.
func loadJSON(data []byte) ([]map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	var docs []map[string]any
	for {
		var v any
		if err := dec.Decode(&v); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to parse JSON input: %w", err)
		}
		switch t := v.(type) {
		case map[string]any:
			docs = append(docs, t)
		case []any:
			for _, item := range t {
				m, ok := item.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("JSON array element is not an object: %T", item)
				}
				docs = append(docs, m)
			}
		default:
			return nil, fmt.Errorf("unsupported JSON top-level value: %T", v)
		}
	}
	return docs, nil
}

// loadYAML decodes a YAML stream, returning one map per document. goccy's
// decoder splits on "---" separators and unmarshals nested mappings into
// map[string]any, keeping the result JSON- and jq-friendly.
func loadYAML(data []byte) ([]map[string]any, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var docs []map[string]any
	for {
		var m map[string]any
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to parse YAML input: %w", err)
		}
		if m == nil {
			continue
		}
		docs = append(docs, m)
	}
	return docs, nil
}
