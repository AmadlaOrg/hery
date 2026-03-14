# Schema | Docs | HERY

Every entity type has a JSON Schema file (`schema.hery.json`) at the root of its repository directory. This schema validates `_meta` and `_body` content for entities of that type.

## Schema File Location

```
application/
  schema.hery.json          # JSON Schema (the type definition)
  default.hery              # Default values (optional)
```

The schema file is **visible** (not in a hidden directory). One schema per entity type directory.

## Schema ID Format

Schema IDs use the HERY URN format:

```
urn:hery:<type-uri-with-/-and-@-replaced-by-:>
```

Example: `urn:hery:amadla.org:entity:application:v1.0.0`

For sub-types, the sub-type name is appended: `urn:hery:amadla.org:entity:application:db:v1.0.0`

## Base Schema

All entity schemas must extend the base HERY schema via `allOf`:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "urn:hery:amadla.org:entity:application:v1.0.0",
  "allOf": [
    { "$ref": "amadla.org/entity/hery@v1.0.0" }
  ]
}
```

The base schema (`amadla.org/entity/hery@v1.0.0`) defines the five reserved properties that all entities inherit: `_type`, `_extends`, `_meta`, `_body`, `_requires`.

## Entity Composition

For nested entity data, the schema uses standard JSON Schema `$ref` to reference sub-entity schemas:

```json
{
  "properties": {
    "_body": {
      "type": "object",
      "properties": {
        "database": {
          "$ref": "urn:hery:amadla.org:entity:application:db:v1.0.0#/properties/_body"
        }
      }
    }
  }
}
```

This tells parsers and downstream tools (weaver, judge) which parts of `_body` correspond to which entity types.

## Validation Flow

1. Parse YAML (resolve anchors/aliases and merge keys)
2. Extract `_type` to determine the entity type
3. Resolve `_type` URI to fetch the entity type (Git clone/pull)
4. Load `schema.hery.json` from the entity type directory
5. Compose the schema (entity schema + base HERY schema via `allOf`)
6. Validate the document against the composed schema
