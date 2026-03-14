# Docs | HERY

HERY (Hierarchical Entity Relational YAML) is a data model and storage system that extends YAML with entity management. It adds five reserved properties to organize content into versioned, schema-validated, relational entities.

**Standard version:** Draft 3.5

## Reserved Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `_type` | string (URI) | Yes* | Entity type URI with version (e.g., `amadla.org/entity/application@v1.0.0`) |
| `_extends` | string (URI) | No | Inherit values from another entity instance via deep merge (data only, no execution ordering) |
| `_meta` | object | No | Metadata for filtering and search (schema-defined) |
| `_body` | object | No | Entity content data (schema-defined) |
| `_requires` | array of strings | No | Hard dependencies on other entities for execution ordering (DAG + topological sort) |

\* `_type` can be inherited via `_extends` when not explicitly set.

## Entity Identity

Entity identity is derived from the **git path** (directory position in the repository), not declared in the document. There is no `_id` property.

## Storage

- **Global entity cache:** `~/.cache/hery/entity/` (cloned entity types, by git path)
- **Project-level SQLite cache:** `.hery.cache` (gitignored, rebuilt from source `.hery` files)
- **Project-level lock file:** `hery.lock` (JSON, committed to Git, portable snapshot of merged state)

## CLI Commands

- `hery entity get` -- Fetch an entity by URI
- `hery entity list` -- List entities
- `hery entity validate` -- Validate entities against their schemas
- `hery query` -- Two-stage query (selection flags + `--jq` transformation)
- `hery compose` -- Compose and merge entities into the lock file and cache

## Further Reading

- [cache.md](cache.md) -- SQLite cache structure
- [compose.md](compose.md) -- Entity composition and merging
- [query.md](query.md) -- Two-stage query model
- [schema.md](schema.md) -- JSON Schema conventions
- [storage.md](storage.md) -- File system layout
- [validation.md](validation.md) -- Entity validation
- [version.md](version.md) -- Versioning model
