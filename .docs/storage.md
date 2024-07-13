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
      - `github.com/` - From 
        - `AmadlaOrg/`
          - `Entity/`
            - `.amadla/`
              - `schema.json`
            - `amadla.yml`
          - `EntityApplication/`
            - ...
    - `amadla.cache` - The [SQLite](https://www.sqlite.org/) caching