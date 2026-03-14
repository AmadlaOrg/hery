# Query | Docs | HERY

HERY uses a **two-stage query model** instead of a custom query language.

## Stage 1: Selection (CLI flags)

Selection flags filter which entities are returned. They operate on SQLite-indexed fields for fast filtering.

| Flag | Description | Example |
|------|-------------|---------|
| `--type`, `-t` | Filter by `_type` (glob) | `--type 'amadla.org/entity/application@v*'` |
| `--meta`, `-m` | Filter by `_meta` field (key=value) | `--meta 'category=Application'` |
| `--tag` | Filter by `_meta.tags` (contains) | `--tag production` |

Multiple flags use AND logic. Repeating the same flag uses OR within that flag.

## Stage 2: Transformation (--jq)

The `--jq` flag applies a jq expression to the selected results, using gojq (compiled into the binary -- no external jq required).

```bash
hery query --type '*/application@*' --jq '.[].\_body.name'
```

When `--jq` is omitted, the full entity JSON is returned.

## Output

- Results are always **JSON arrays** (even for single results)
- With `--jq`: follows jq output conventions (newline-delimited JSON)
- `--raw` / `-r`: raw output (no JSON quotes for strings)
- `--compact` / `-c`: compact JSON (no whitespace)
- `--count`: return only the count of matching entities

## Layer Inspection

By default, `hery query` returns the merged view. Use `--layers` to see individual layers before merge, or `--layer <n>` for a specific layer.

## Examples

```bash
# All entities
hery query

# Filter by type
hery query --type 'amadla.org/entity/application@v*'

# Filter by type, extract a field
hery query --type '*/application@*' --jq '.[].\_body.name'

# Filter by meta tag
hery query --tag production --jq '.[].\_meta.name'

# Pipe to downstream tools (UNIX integration)
hery query --type '*/application@*' | doorman inject | weaver render
```
