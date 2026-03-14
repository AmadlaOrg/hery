# Validation | Docs | HERY

Entities are validated against their JSON Schema to ensure predictable structure and content types.

## How It Works

1. hery parses the YAML entity content
2. Extracts `_type` to determine the entity type
3. Resolves the `_type` URI to fetch the entity type (Git clone/pull to `~/.cache/hery/entity/`)
4. Loads `schema.hery.json` from the entity type directory
5. Composes the schema (entity schema + base HERY schema via `allOf`)
6. Validates the document against the composed schema

## What Is Validated

- `_type` -- must be a valid URI string (required, unless inherited via `_extends`)
- `_extends` -- must be a valid URI string, must resolve to same `_type` (if present)
- `_meta` -- must conform to the entity's meta schema (if present)
- `_body` -- must conform to the entity's body schema (if present)
- `_requires` -- must be an array of valid URI strings (if present); no `../` escape in relative paths
- No unknown `_`-prefixed properties at document root

## Usage

```bash
# Validate a specific entity
hery entity validate github.com/AmadlaOrg/Entity@v1.0.0
```

## Planned

The following validation subcommands are not yet implemented:

- **Entity URI** -- Verify that it is well formatted and exists with the version provided
- **Entity version** -- Only verify the version passed
- **Entity hash** -- Verify if the repository/entity was downloaded correctly
- **Have** -- Verify if the entity exists on the local machine

See [roadmap.md](roadmap.md) for details.
