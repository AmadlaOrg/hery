# Cache | Docs | HERY

Entity data is cached in SQLite for fast querying.

- **Cache file:** `.hery.cache` at the project root (gitignored, derived data)
- The `.hery` source files are the source of truth
- Cache is rebuilt from source files as needed

> [!NOTE]
> hery uses embedded SQLite (go-sqlite3). No external database is needed.

## Primary Table: `entities`

The `entities` table stores all entity data in a single flat structure:

| Column | Type | Description |
|--------|------|-------------|
| `entity_type` | TEXT NOT NULL | `_type` URI with version |
| `entity_extends` | TEXT | `_extends` URI |
| `uri` | TEXT UNIQUE | Full entity URI |
| `name` | TEXT | Simple name of the entity |
| `repo_url` | TEXT | Full URL to the repository |
| `origin` | TEXT | Partial path of the entity |
| `version` | TEXT | Entity version |
| `is_latest_version` | BOOLEAN | Whether this is the latest version |
| `is_pseudo_version` | BOOLEAN | Whether this is a pseudo-version |
| `abs_path` | TEXT UNIQUE | Full system path to entity files |
| `have` | BOOLEAN | Whether entity content is on local machine |
| `hash` | TEXT UNIQUE | Hash of entity content for validation |
| `exist` | BOOLEAN | Whether the repository was found |
| `schema_json` | TEXT | Full JSON Schema content |
| `meta_json` | TEXT | `_meta` as JSON string |
| `body_json` | TEXT | `_body` as JSON string |
| `requires_json` | TEXT | `_requires` as JSON string |
| `merged_json` | TEXT | Full merged entity as JSON |
| `source_file` | TEXT | Origin `.hery` file path |
| `insert_date_time` | DATETIME | Row creation timestamp |
| `update_date_time` | DATETIME | Row update timestamp |

## Indexes

Indexed columns for fast selection queries:

- `entity_type` -- for `--type` flag filtering
- `name` -- for name lookups
- `uri` -- for URI lookups
- `version` -- for version filtering
- `is_latest_version` -- for latest-version queries

## Querying

`hery query` loads `merged_json` from the cache and runs the same in-memory
selection predicates (case-insensitive type glob with version stripping,
`_meta` substring match) as the `--from` and `--dir` sources, so results are
identical no matter where the data came from. SQL-side filtering on the
indexed columns is a possible future optimization.

## Staleness: `cache_manifest`

The cache records what it was built from in the `cache_manifest` table:

| Column | Type | Description |
|--------|------|-------------|
| `path` | TEXT PRIMARY KEY | Absolute path of a source `.hery` file or directory |
| `is_dir` | BOOLEAN | TRUE for directory rows |
| `mtime_ns` | INTEGER | File modification time (ns); 0 for directories |
| `size` | INTEGER | File size in bytes; 0 for directories |
| `hery_names` | TEXT | Directory rows: sorted JSON array of `.hery` filenames |

A file row goes stale when its mtime or size changes; a directory row goes
stale when the set of `.hery` filenames changes (catching added or removed
files). Directory mtimes are deliberately not used — writing `.hery.cache`
into the project root would perturb them. Any mismatch, missing table, or
read error makes `hery query` re-resolve the sources and rebuild the cache
transparently.

## Additional Tables

The `body_data_*` tables (TEXT, NUMERIC, REAL, BOOLEAN, DATE, DATETIME) store typed property values for potential future indexed property queries. The primary query path uses the JSON columns in the `entities` table.

See `resources/sql/hery-tables.sql` for the full schema.
