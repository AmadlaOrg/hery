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
- `.collection`
- `.collection/schema.json`
- `collection.yml`