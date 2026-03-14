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

The `meta_json` and `body_json` columns are queried using SQLite `json_extract()`:

```sql
-- --meta 'category=Application'
SELECT merged_json FROM entities
WHERE json_extract(meta_json, '$.category') = 'Application';

-- --tag production
SELECT merged_json FROM entities
WHERE EXISTS (
    SELECT 1 FROM json_each(json_extract(meta_json, '$.tags'))
    WHERE value = 'production'
);
```

## Additional Tables

The `body_data_*` tables (TEXT, NUMERIC, REAL, BOOLEAN, DATE, DATETIME) store typed property values for potential future indexed property queries. The primary query path uses the JSON columns in the `entities` table.

See `resources/sql/hery-tables.sql` for the full schema.
