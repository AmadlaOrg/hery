package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity/resolve"
	"github.com/goccy/go-yaml"
)

// ProjectCacheFile is the derived SQLite cache at the project root. The .hery
// source files are the source of truth; the cache is rebuilt whenever they
// change (see .docs/cache.md).
const ProjectCacheFile = ".hery.cache"

// ProjectCache is the project-level query cache: a materialized snapshot of a
// resolved .hery directory, plus a manifest of the source files it was built
// from so staleness is detectable with stat calls alone.
type ProjectCache interface {
	// LoadFresh returns the cached entity documents when the cache exists and
	// every manifest entry still matches the filesystem. Any miss, mismatch,
	// or error reads as "not fresh" — the caller falls back to resolving.
	LoadFresh() ([]map[string]any, bool)

	// Rebuild replaces the cache contents with the docs of a resolve.Result
	// and records the source manifest (every layer directory and .hery file
	// with its mtime and size).
	Rebuild(res *resolve.Result) error
}

type projectCache struct {
	dbPath string
}

// NewProject returns the ProjectCache for a project directory.
func NewProject(dir string) ProjectCache {
	return &projectCache{dbPath: filepath.Join(dir, ProjectCacheFile)}
}

func (p *projectCache) LoadFresh() ([]map[string]any, bool) {
	if _, err := os.Stat(p.dbPath); err != nil {
		return nil, false
	}

	db := database.New(p.dbPath)
	if err := db.Initialize(); err != nil {
		return nil, false
	}
	defer db.Close()

	manifest, err := db.QueryRows("SELECT path, is_dir, mtime_ns, size, hery_names FROM cache_manifest")
	if err != nil || len(manifest) == 0 {
		return nil, false
	}
	for _, m := range manifest {
		path, _ := m["path"].(string)
		if toInt64(m["is_dir"]) != 0 {
			// Directory staleness is a listing check: the set of .hery files
			// must be unchanged. Directory mtimes are useless here — writing
			// .hery.cache into the project root would perturb them.
			names, _ := m["hery_names"].(string)
			current, listErr := heryNamesJSON(path)
			if listErr != nil || current != names {
				return nil, false
			}
			continue
		}
		info, statErr := os.Stat(path)
		if statErr != nil || info.IsDir() ||
			info.ModTime().UnixNano() != toInt64(m["mtime_ns"]) ||
			info.Size() != toInt64(m["size"]) {
			return nil, false
		}
	}

	rows, err := db.QueryRows("SELECT merged_json FROM entities")
	if err != nil {
		return nil, false
	}
	docs := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		s, _ := r["merged_json"].(string)
		if s == "" {
			continue
		}
		var doc map[string]any
		if json.Unmarshal([]byte(s), &doc) != nil {
			return nil, false
		}
		docs = append(docs, doc)
	}
	return docs, true
}

func (p *projectCache) Rebuild(res *resolve.Result) error {
	db := database.New(p.dbPath)
	if err := db.Initialize(); err != nil {
		return err
	}
	defer db.Close()

	// Apply merges its queues as CreateTable → Insert → Delete, so clearing
	// the previous contents must commit in its own round before the inserts
	// are queued.
	db.CreateTable()
	db.Delete(database.Table{Name: "entities"}, database.SelectClauses{})
	db.Delete(database.Table{Name: "cache_manifest"}, database.SelectClauses{})
	if err := db.Apply(); err != nil {
		return err
	}

	for _, layer := range res.Layers {
		for _, doc := range layer.Docs {
			row, err := docRow(doc)
			if err != nil {
				return err
			}
			db.Insert(database.Table{Name: "entities", Rows: []database.Row{row}})
		}
	}

	// The manifest stats the sources after resolution, so a file changing
	// mid-build is recorded with its new mtime and old content — the next
	// LoadFresh sees it as fresh one run late, never as silently wrong
	// forever; any later change re-triggers a rebuild.
	manifest, err := manifestRows(res)
	if err != nil {
		return err
	}
	db.Insert(database.Table{Name: "cache_manifest", Rows: manifest})

	return db.Apply()
}

// docRow converts one resolved doc into an entities-table row. Only the
// columns knowable from a local .hery source are set; repo-oriented columns
// (uri, repo_url, version, hash, …) stay NULL.
func docRow(d resolve.Doc) (database.Row, error) {
	raw, err := yaml.Marshal(d.Raw)
	if err != nil {
		return nil, fmt.Errorf("%s: marshal doc: %w", d.Path, err)
	}
	var full map[string]any
	if err := yaml.Unmarshal(raw, &full); err != nil {
		return nil, fmt.Errorf("%s: parse doc: %w", d.Path, err)
	}
	mergedJSON, err := json.Marshal(full)
	if err != nil {
		return nil, fmt.Errorf("%s: encode doc: %w", d.Path, err)
	}

	row := database.Row{
		"entity_type": d.Type,
		"merged_json": string(mergedJSON),
		"source_file": d.Path,
		"abs_path":    d.Path,
	}
	if d.Extends != "" {
		row["entity_extends"] = d.Extends
	}
	if meta, ok := full["_meta"].(map[string]any); ok {
		if b, jsonErr := json.Marshal(meta); jsonErr == nil {
			row["meta_json"] = string(b)
		}
		if name, nok := meta["name"].(string); nok {
			row["name"] = name
		}
	}
	if body, ok := full["_body"]; ok {
		if b, jsonErr := json.Marshal(body); jsonErr == nil {
			row["body_json"] = string(b)
		}
	}
	if len(d.Requires) > 0 {
		if b, jsonErr := json.Marshal(d.Requires); jsonErr == nil {
			row["requires_json"] = string(b)
		}
	}
	return row, nil
}

// manifestRows records every layer directory (as its sorted .hery listing,
// catching added/removed files) and every source file (as mtime and size).
func manifestRows(res *resolve.Result) ([]database.Row, error) {
	var rows []database.Row
	for _, layer := range res.Layers {
		names, err := heryNamesJSON(layer.Dir)
		if err != nil {
			return nil, fmt.Errorf("list cache source %s: %w", layer.Dir, err)
		}
		rows = append(rows, database.Row{
			"path":       layer.Dir,
			"is_dir":     true,
			"mtime_ns":   0,
			"size":       0,
			"hery_names": names,
		})
		for _, doc := range layer.Docs {
			info, statErr := os.Stat(doc.Path)
			if statErr != nil {
				return nil, fmt.Errorf("stat cache source %s: %w", doc.Path, statErr)
			}
			rows = append(rows, database.Row{
				"path":     doc.Path,
				"is_dir":   false,
				"mtime_ns": info.ModTime().UnixNano(),
				"size":     info.Size(),
			})
		}
	}
	return rows, nil
}

// heryNamesJSON returns the directory's .hery filenames, sorted, as a JSON
// array — the canonical form compared between manifest and filesystem.
func heryNamesJSON(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".hery" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	b, err := json.Marshal(names)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case bool:
		if n {
			return 1
		}
	}
	return 0
}
