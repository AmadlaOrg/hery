# Validation | Docs | HERY
To make sure that the entities are used properly they need to be predictable in their structure and content types. The 
solution is having a schema definition and one of the solutions to accomplish that is with [JSON Schema](https://json-schema.org/).

The `hery`-cli transforms the YAML entity content into a JSON string that is then validated against the `.collection/schema.json`
file.

There is also other subcommands to validate other parts of HERY: 
- **Entity** - An entity itself
- **Entity URI** - Verifies that it is well formated and exists with the version provided
- **Entity version** - Only verifies the version pass
- **Entity hash** - Useful to verify if the repository/entity was well downloaded
- **Have** - To verify if the entity exist on the local machine

Example: 
```bash
# Entity
hery validation entity --collection="amadla" github.com/AmadlaOrg/Entity@v1.0.0 # The collection needs to be provided

# URI
hery validation uri github.com/AmadlaOrg/Entity@v1.0.0
hery validation uri github.com/AmadlaOrg/Entity

# Version
hery validation version v1.0.0
hery validation version github.com/AmadlaOrg/Entity@v1.0.0 # This way it will check if the version provided actually exist

# Hash
hery validation hash --collection="amadla" github.com/AmadlaOrg/Entity@v1.0.0 # The collection needs to be provided

# Have
hery validation hash --collection="amadla" github.com/AmadlaOrg/Entity@v1.0.0 # The collection needs to be provided
```