// Package resolve walks a directory of .hery files, resolves _extends chains
// (deep-merging parent _body into child _body while preserving child _type,
// _meta, _extends declaration, and _requires), and tracks cross-directory
// references as discrete merge layers.
//
// Layer model: the starting directory is layer 1. Any _extends or _requires
// reference that resolves to a different directory becomes a deeper layer,
// discovered in BFS order.
//
// Intra-layer ordering: docs with no incoming local _requires are local roots,
// emitted first; their _requires-transitive docs follow in DFS order,
// alphabetical tie-break on file path.
package resolve

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AmadlaOrg/hery/entity/merge"
	"github.com/goccy/go-yaml"
)

// Doc holds a single parsed .hery document along with its reserved fields
// lifted out for convenience. Raw is the order-preserving MapSlice used for
// output.
type Doc struct {
	Path     string
	Raw      yaml.MapSlice
	Type     string
	Extends  string
	Requires []string
}

// Layer is one merge depth — a directory's docs after _extends resolution,
// in emission order.
type Layer struct {
	Dir  string
	Docs []Doc
}

// Result is the full resolution output. Layers[0] is layer 1 (the starting
// directory). Warnings holds non-fatal resolution notices (e.g. type-URI or
// external references that cannot be resolved locally and were skipped).
type Result struct {
	Layers   []Layer
	Warnings []string
}

// Resolver loads a directory tree into layered, resolved Docs.
type Resolver interface {
	Resolve(dir string) (*Result, error)
}

type resolver struct{}

// New returns a Resolver.
func New() Resolver {
	return &resolver{}
}

// Resolve walks dir, resolves _extends, and discovers cross-directory
// references as new layers via BFS.
func (r *resolver) Resolve(dir string) (*Result, error) {
	absStart, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve abs path: %w", err)
	}

	visited := map[string]int{}
	var layers []Layer
	var warnings []string
	queue := []string{absStart}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if _, seen := visited[cur]; seen {
			continue
		}

		docs, err := readDir(cur)
		if err != nil {
			return nil, err
		}
		docs = orderDocs(cur, docs)

		for i := range docs {
			if docs[i].Extends != "" {
				switch {
				case !isLocalRef(docs[i].Extends):
					warnings = append(warnings, fmt.Sprintf(
						"%s: _extends: skipping %q — only ./ and ../ resolve locally; emitting entity un-merged",
						docs[i].Path, docs[i].Extends))
				default:
					extendsAbs, err := resolveRef(cur, docs[i].Extends)
					if err != nil {
						return nil, fmt.Errorf("%s: _extends: %w", docs[i].Path, err)
					}
					if err := applyExtends(&docs[i], extendsAbs); err != nil {
						return nil, err
					}
					extendsDir := filepath.Dir(extendsAbs)
					if extendsDir != cur {
						queue = append(queue, extendsDir)
					}
				}
			}

			for _, req := range docs[i].Requires {
				if !isLocalRef(req) {
					warnings = append(warnings, fmt.Sprintf(
						"%s: _requires: %q is not a local ./ or ../ path — declaration kept, not walked",
						docs[i].Path, req))
					continue
				}
				refAbs, err := resolveRef(cur, req)
				if err != nil {
					return nil, fmt.Errorf("%s: _requires %q: %w", docs[i].Path, req, err)
				}
				reqDir := refAbs
				if info, err := os.Stat(refAbs); err == nil && !info.IsDir() {
					reqDir = filepath.Dir(refAbs)
				}
				if reqDir != cur {
					queue = append(queue, reqDir)
				}
			}
		}

		visited[cur] = len(layers)
		layers = append(layers, Layer{Dir: cur, Docs: docs})
	}

	return &Result{Layers: layers, Warnings: warnings}, nil
}

// readDir loads every .hery file in dir (non-recursive) into Docs, ordered
// by file name.
func readDir(dir string) ([]Doc, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}

	var docs []Doc
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".hery" {
			continue
		}
		full := filepath.Join(dir, e.Name())
		doc, err := readDoc(full)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	sort.Slice(docs, func(i, j int) bool { return docs[i].Path < docs[j].Path })
	return docs, nil
}

// readDoc parses one .hery file into a Doc, preserving key order via
// yaml.MapSlice.
func readDoc(path string) (Doc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Doc{}, fmt.Errorf("read %s: %w", path, err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(data), yaml.UseOrderedMap())
	var raw any
	if err := dec.Decode(&raw); err != nil {
		return Doc{}, fmt.Errorf("decode %s: %w", path, err)
	}
	ms, ok := raw.(yaml.MapSlice)
	if !ok {
		return Doc{}, fmt.Errorf("decode %s: top-level must be a mapping", path)
	}

	doc := Doc{Path: path, Raw: ms}
	for _, item := range ms {
		key, _ := item.Key.(string)
		switch key {
		case "_type":
			if s, ok := item.Value.(string); ok {
				doc.Type = s
			}
		case "_extends":
			if s, ok := item.Value.(string); ok {
				doc.Extends = s
			}
		case "_requires":
			if arr, ok := item.Value.([]any); ok {
				for _, v := range arr {
					if s, ok := v.(string); ok {
						doc.Requires = append(doc.Requires, s)
					}
				}
			}
		}
	}
	return doc, nil
}

// orderDocs arranges docs by local-root DFS following _requires within the
// same directory. Docs not required by any other local doc are roots, emitted
// in alphabetical path order; their local _requires chain follows.
func orderDocs(dir string, docs []Doc) []Doc {
	if len(docs) <= 1 {
		return docs
	}

	byPath := make(map[string]int, len(docs))
	for i, d := range docs {
		byPath[d.Path] = i
	}

	// Build local adjacency: docPath -> list of local _requires targets.
	adj := make(map[string][]string, len(docs))
	incoming := make(map[string]int, len(docs))
	for _, d := range docs {
		for _, req := range d.Requires {
			target, err := resolveRef(dir, req)
			if err != nil {
				continue
			}
			if _, ok := byPath[target]; ok {
				adj[d.Path] = append(adj[d.Path], target)
				incoming[target]++
			}
		}
	}

	var roots []string
	for _, d := range docs {
		if incoming[d.Path] == 0 {
			roots = append(roots, d.Path)
		}
	}
	sort.Strings(roots)

	for _, children := range adj {
		sort.Strings(children)
	}

	var ordered []Doc
	visited := make(map[string]bool, len(docs))
	var dfs func(path string)
	dfs = func(path string) {
		if visited[path] {
			return
		}
		visited[path] = true
		ordered = append(ordered, docs[byPath[path]])
		for _, child := range adj[path] {
			dfs(child)
		}
	}
	for _, root := range roots {
		dfs(root)
	}
	// Catch any unreachable docs (cycles or orphans): append in path order.
	for _, d := range docs {
		if !visited[d.Path] {
			ordered = append(ordered, d)
		}
	}
	return ordered
}

// applyExtends reads the extended file and deep-merges its _body into
// child._body, keeping child's _type, _meta, _extends (declaration), _requires.
func applyExtends(child *Doc, extendsPath string) error {
	parent, err := readDoc(extendsPath)
	if err != nil {
		return fmt.Errorf("_extends target %s: %w", extendsPath, err)
	}
	parentBody := extractBody(parent.Raw)
	childBody := extractBody(child.Raw)
	merged := merge.OrderedDeepMerge(parentBody, childBody)
	child.Raw = replaceBody(child.Raw, merged)
	return nil
}

func extractBody(ms yaml.MapSlice) yaml.MapSlice {
	for _, item := range ms {
		if k, _ := item.Key.(string); k == "_body" {
			if body, ok := item.Value.(yaml.MapSlice); ok {
				return body
			}
		}
	}
	return nil
}

func replaceBody(ms yaml.MapSlice, body yaml.MapSlice) yaml.MapSlice {
	out := make(yaml.MapSlice, 0, len(ms)+1)
	replaced := false
	for _, item := range ms {
		if k, _ := item.Key.(string); k == "_body" {
			out = append(out, yaml.MapItem{Key: item.Key, Value: body})
			replaced = true
			continue
		}
		out = append(out, item)
	}
	if !replaced {
		out = append(out, yaml.MapItem{Key: "_body", Value: body})
	}
	return out
}

// isLocalRef reports whether ref is a local filesystem reference (./ or ../).
// Anything else (type URIs, external github.com/... refs) cannot be resolved
// by the local directory walker.
func isLocalRef(ref string) bool {
	return strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../")
}

// resolveRef resolves a _extends or _requires reference relative to fromDir.
// Trailing slash is preserved for directory references.
func resolveRef(fromDir, ref string) (string, error) {
	if isLocalRef(ref) {
		joined := filepath.Join(fromDir, ref)
		return filepath.Clean(joined), nil
	}
	return "", fmt.Errorf("unsupported reference %q (only ./ and ../ are supported in MVP)", ref)
}

// Marshal renders a slice of Docs as a multi-document YAML stream, prefixing
// each document with a "---" separator so the output is valid multi-doc YAML
// that round-trips through a standard decoder (e.g. `hery query --from -`).
func Marshal(docs []Doc) ([]byte, error) {
	var buf bytes.Buffer
	for _, d := range docs {
		buf.WriteString("---\n")
		enc := yaml.NewEncoder(&buf, yaml.Indent(2), yaml.IndentSequence(true))
		if err := enc.Encode(d.Raw); err != nil {
			return nil, err
		}
		if err := enc.Close(); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// MarshalAll renders every layer's docs in order.
func MarshalAll(layers []Layer) ([]byte, error) {
	var buf bytes.Buffer
	for _, layer := range layers {
		b, err := Marshal(layer.Docs)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}
