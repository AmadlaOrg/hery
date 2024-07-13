# Storage | Docs | HERY
All the entities are stored in two possible places: 
- The root of a project
  - `{project root}/.henry/`
- The home directory: 
  - Unix: `~/.hery/`
  - Windows: `{APPDATA}\Hery\`

Inside the storage directory `./hery` the next layer of directories are of collections.

Example of a storage structure: 
- `hery/` - The root of the storage for collections of entities
  - `amadla/` - An example of a collection
    - `entity/` - All the entities are stored inside of this directory
      - `github.com/` - From this point forward the path to the entity is broken down into directories
        - `AmadlaOrg/`
          - `Entity@v1.0.0/` - Entity name with the version
            - `.amadla/` - The directory is always named after the collection
              - `schema.json` - The schema for the entity
            - `amadla.yml` - Is an entity with data, normally used to define an entity
          - `EntityApplication@v1.0.0/` - This is just another entity
            - ...
    - `amadla.cache` - The [SQLite](https://www.sqlite.org/) caching