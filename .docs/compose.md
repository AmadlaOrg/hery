# Compose | Docs | HERY

> [!NOTE]
> Compose is not yet fully implemented. See [roadmap.md](roadmap.md) for planned features.

With compose the entities are put together and merged into a lock file that is then added to the cache
SQLite database for quick querying. Compose starts with the main entity and from that entity
resolves all dependencies.

Since the `hery.lock` is a JSON file it is possible to extract the data from the entities without using the caching
system, but it will be slower.
