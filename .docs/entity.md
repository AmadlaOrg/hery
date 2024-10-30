# Entity | Docs | HERY
An entity is:
- A block of YAML
- A YAML that comes with a [JSON Schema](https://json-schema.org/)
- A YAML block that is named by a URI that contains a version number and points to related resources

In an entity directory there is a basic file and a directory. The file is at the root of the directory or repository of
the entity and is named after the collection with a `.yml` or `.yaml` file extension.

For the directory it is dot and the name of the collection (e.g.: `.amadla/`). This directory contains the `schema.json`
[JSON Schema](https://json-schema.org/) file. It is always named: `schema.json`. It is possible to have multiple
[JSON Schemas](https://json-schema.org/), but they won't be connected automatically to the entity definition. That will
have to be done manually following the [JSON Schema](https://json-schema.org/) documentation by adding a full URL to
the schema file in the schema file of choosing. The directory can contain other files that are connected to the entity.

The only important detail to retain is to not conflict with these three basic files and directory:
- `.<collection name>`
- `.<collection name>/schema.hery.json`
- `<collection name>.hery`

## Properties
- `_entity` - Contains the URI (without the protocol) to the entity repository with the version (e.g.: `github.com/AmadlaOrg/Entity@latest`)
- `_body` - Contains the entity content
- `_id` - Used as reference to a specific entity content
- `_meta` - Contains metadata for a related entity

### `_entity`
Useful to identify an entity and point to where to get the entity.

### `_body`

### `_id`

### `_meta`
Useful for querying entities.

## Entity Examples:

1. With metadata for the entity
A basic example.
```hery
---
_meta:
  _entity: github.com/AmadlaOrg/Entity@latest
  _body:
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

2. Without metadata entity
In this example the entity is: github.com/AmadlaOrg/Entity@latest.
```hery
---
name: Random Name
```

3. Multi-layered entity
This example is to show a multi-layered entity.
```hery
---
_meta:
  _entity: github.com/AmadlaOrg/Entity@latest
  _body:
    name: RandomName
    description: Some description.
    category: QA
_id: "random:ID:1234"
_body:
  subject: Some random subject.
  listing:
    - Apple
    - Orange
    - Grape
  external:
    _entity: github.com/AmadlaOrg/QAFixturesSubEntityWithMultiSubEntities@latest
    _body:
      message: Some message...
  external-list:
    - _entity: github.com/AmadlaOrg/QAFixturesSubEntityWithMultiSubEntities@latest
      _body:
        message: Another random message.
    - _entity: github.com/AmadlaOrg/QAFixturesSubEntityWithMultiSubEntities@latest
      _body:
        message: Again, another random message.
    - _entity: github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@latest
      _body:
        title: Hello World!
```

### Create Entity
1. Create repository
2. Make a dot directory with the collection name at the root of the repository:
    - `mkdir -p .<collection name>`
3. Make `.hery` configuration file with the collection name at the root of the repository
    - `touch <collection name>.hery`
4. Make a `schema.json` configuration file in the `.<collection name>` directory
    - `touch ./.<collection name>/schema.json`
5. Add the content in `.hery` configuration file and the `schema.json` [JSON-Schema](https://json-schema.org/)
6. Add it in git:
    - `git add .`
    - `git commit -m "Batman"`
