# Entity | Docs | HERY
An entity is:
- A block of YAML
- A YAML that comes with a [JSON Schema](https://json-schema.org/)
- A YAML block that is named by a URI that contains a version number and points to related resources

In an entity directory there is a basic file and directory structure. The file at the root of the entity repository
is named with a `.hery` file extension.

The schema file is a `<name>.hery.json` file located at the root of the entity type directory (e.g., `application.hery.json`, `package.hery.json`). There must be exactly one `.hery.json` file per entity type directory.

The key files for an entity:
- `<name>.hery.json` — entity JSON Schema at the entity type directory root (e.g., `application.hery.json`)
- `*.hery` — entity content files

## Properties
- `_type` — Contains the URI (without the protocol) to the entity repository with the version (e.g.: `github.com/AmadlaOrg/Entity@latest`). The version maps directly to a Git tag (e.g., `@v1.0.0` → git tag `v1.0.0`)
- `_extends` — Optional URI to parent entity instance for deep merge inheritance
- `_meta` — Contains metadata for the entity
- `_body` — Contains the entity content
- `_requires` — Optional list of hard dependencies on other entities for execution ordering

### `_type`
Identifies the entity type and points to where to get the entity definition and schema. Required in every `.hery` file (unless inherited via `_extends`).

When used with `_extends`, both can be present: `_type` declares the schema (what this entity *is*), `_extends` declares data inheritance (who it inherits from). When only `_extends` is present, `_type` is inherited from the extended entity. If both are present and conflict, hery warns but allows it.

### `_extends`
Points to an entity instance of the same `_type` to inherit from. Purely a data/merge operation — no execution ordering implied. Deep merge applies to `_body`, `_meta`, and `_requires` — the child's values override the extended entity's: objects merge recursively, arrays replace entirely, scalars are overridden by the child. `_type` is never merged.

Targets specific elements via the filename as a path segment: `_extends: github.com/SomeOrg/WordPress/database.hery`. One addressable element = one file.

### `_meta`
Optional metadata for the entity. Useful for querying and organizing entities.

### `_body`
Contains the entity data. Optional — when omitted, the entity inherits all default values from the extended entity
definition. When present, its content is validated against the entity's JSON Schema.

### `_requires`
Declares hard dependencies on other entities (list). Used by amadla to build a dependency graph (DAG) and determine execution order via topological sort. hery validates syntax at parse time (valid URIs, no `../` escape).

References can be:
- Entity type URIs: `github.com/AmadlaOrg/Entities/Application/DB/RDBMS@^v1.0.0` (at least one must exist)
- Specific elements: `github.com/SomeOrg/WordPress#php.hery`
- Local elements: `#database.hery` (same entity directory)

Version constraints are supported: `@v1.0.0` (exact), `@^v1.0.0` (compatible range). Relative paths are sandboxed to the entity directory (no `../` escape). `_requires: []` is valid and equivalent to omitting the property.

## Entity Examples:

1. Basic entity with metadata
```hery
# yaml-language-server: $schema=https://amadla.org/entity/hery/v1.0.0/schema.hery.json
---
_type: github.com/AmadlaOrg/EntityQA/RandomName@latest
_meta:
  name: RandomName
  description: Entity Pseudo Version definitions.
  category: QA
  tags:
    - QA
    - fixture
    - test
_body:
  name: Random Name
```

2. Minimal entity without metadata
```hery
# yaml-language-server: $schema=https://amadla.org/entity/hery/v1.0.0/schema.hery.json
---
_type: github.com/AmadlaOrg/Entity@latest
_body:
  name: Random Name
```

3. Entity with extends inheritance
```hery
# yaml-language-server: $schema=https://amadla.org/entity/hery/v1.0.0/schema.hery.json
---
_type: github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0
_extends: github.com/some-org/base-configs/webserver
_meta:
  name: MyWebServer
  description: Custom web server configuration.
  category: Application
_body:
  server_name: myserver.dev
  port: 443
  network:
    ports:
      - 80
      - 443
```

### Create Entity
1. Create repository
2. Create a `<name>.hery.json` [JSON-Schema](https://json-schema.org/) at the repository root (e.g., `application.hery.json`)
3. Create one or more `.hery` content files
4. Add it in git:
    - `git add .`
    - `git commit -m "Initial entity"`
