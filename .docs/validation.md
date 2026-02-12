# Validation | Docs | HERY
To make sure that the entities are used properly they need to be predictable in their structure and content types. The
solution is having a schema definition and one of the solutions to accomplish that is with [JSON Schema](https://json-schema.org/).

The `hery`-cli transforms the YAML entity content into a JSON string that is then validated against the `.<collection name>/schema.hery.json`
file.

## Implemented

- **Entity** - Validate an entity against its schema

Example:
```bash
# Entity
hery entity validate --collection="amadla" github.com/AmadlaOrg/Entity@v1.0.0
```

## Planned

The following validation subcommands are not yet implemented:

- **Entity URI** - Verify that it is well formatted and exists with the version provided
- **Entity version** - Only verify the version passed
- **Entity hash** - Useful to verify if the repository/entity was downloaded correctly
- **Have** - Verify if the entity exists on the local machine

See [roadmap.md](roadmap.md) for details.