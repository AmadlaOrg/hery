# Roadmap | Docs | HERY
Features that are planned but not yet implemented.

## Server / Client Mode
When running the server the [SQLite](https://www.sqlite.org/) database would be loaded in memory for faster access to
cached entities for querying. The server would be accessible via `/tmp/hery.sock` UNIX socket file.

```bash
hery server start --collection=amadla
hery server stop --collection=amadla
hery client --collection=amadla
```

## Compose Lock File
Compose would put entities together and merge them into a `collection.lock` JSON file that is then added to the cache
[SQLite](https://www.sqlite.org/) database for quick querying. Since the lock file is JSON, it would be possible to
extract data from entities without using the caching system (but slower).

## Rich Query Functions
The following query functions are planned:

| Function | Description |
|----------|-------------|
| `equal(<property>, <value>)` | Exact match filter |
| `like(<property>, <pattern>)` | Pattern matching with wildcards (`%` and `_`) |
| `in(<property>, <values>...)` | Filter where value matches any in a list |
| `not_equal(<property>, <value>)` | Inverse exact match |
| `not_like(<property>, <pattern>)` | Inverse pattern matching |
| `not_in(<property>, <values>...)` | Inverse list match |
| `and(<functions>...)` | Combine filters with AND logic |
| `or(<functions>...)` | Combine filters with OR logic |
| `limit(<number>)` | Maximum number of rows to return |
| `offset(<number>)` | Number of rows to skip |
| `order_by(<property> [, ASC\|DESC])` | Order output by a property |
| `group_by(<property> [, ...])` | Group by one or more properties |

## Additional Validation Subcommands
The following validation subcommands are planned:

- `hery entity validate-uri <URI>` - Verify that a URI is well formatted and the version exists
- `hery entity validate-version <version>` - Verify version format
- `hery entity validate-hash <URI>` - Verify repository/entity download integrity
- `hery entity have <URI>` - Verify if the entity exists on the local machine