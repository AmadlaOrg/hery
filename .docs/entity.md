# Entity | Docs | HERY
An entity is:
- A block of YAML
- A YAML that comes with a [JSON Schema](https://json-schema.org/)
- A YAML block that is named by a URI that contains a version number and points to related resources

In an entity directory there is a basic file and directory structure. The file at the root of the entity repository
is named with a `.hery` file extension.

The schema file `schema.hery.json` is located at the root of the entity type directory. It is possible
to have multiple [JSON Schemas](https://json-schema.org/), but they won't be connected automatically to the entity
definition. That will have to be done manually following the [JSON Schema](https://json-schema.org/) documentation by
adding a full URL to the schema file in the schema file of choosing.

The key files for an entity:
- `schema.hery.json` — entity JSON Schema at the entity type directory root
- `*.hery` — entity content files

## Properties
- `_type` — Contains the URI (without the protocol) to the entity repository with the version (e.g.: `github.com/AmadlaOrg/Entity@latest`)
- `_self` — Optional self-referencing identifier for a specific entity content instance
- `_parent` — Optional URI to parent entity instance for deep merge inheritance
- `_meta` — Contains metadata for the entity
- `_body` — Contains the entity content

### `_type`
Identifies the entity type and points to where to get the entity definition and schema. Required in every `.hery` file.

### `_self`
A resolvable URI identifier for a specific entity content instance. Acts as the merge discriminator: when entities
of the same type are merged across layers, same `_self` means override (child wins), different `_self` means
accumulate (new entry). If omitted, a UUID v4 is auto-generated.

### `_parent`
Points to a parent entity instance of the same `_type`. Enables deep merge inheritance where the child's values
override the parent's: objects merge recursively, arrays replace entirely, scalars are overridden by the child.

### `_meta`
Optional metadata for the entity. Useful for querying and organizing entities.

### `_body`
Contains the entity data. Optional — when omitted, the entity inherits all default values from the parent entity
definition. When present, its content is validated against the entity's JSON Schema.

## Entity Examples:

1. Basic entity with metadata
```hery
# yaml-language-server: $schema=amadla.org/schemas/hery@v1.0.0
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
# yaml-language-server: $schema=amadla.org/schemas/hery@v1.0.0
---
_type: github.com/AmadlaOrg/Entity@latest
_body:
  name: Random Name
```

3. Entity with parent inheritance
```hery
# yaml-language-server: $schema=amadla.org/schemas/hery@v1.0.0
---
_type: github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0
_self: "my-webserver"
_parent: "github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0#base"
_meta:
  name: MyWebServer
  description: Custom web server configuration.
  category: Application
_body:
  subject: Some random subject.
  listing:
    - Apple
    - Orange
    - Grape
  external:
    _type: github.com/AmadlaOrg/EntitySystem/Net@v1.0.0
    _body:
      ports:
        - 80
        - 443
```

### Create Entity
1. Create repository
2. Create a `schema.hery.json` [JSON-Schema](https://json-schema.org/) at the repository root
3. Create one or more `.hery` content files
4. Add it in git:
    - `git add .`
    - `git commit -m "Initial entity"`
