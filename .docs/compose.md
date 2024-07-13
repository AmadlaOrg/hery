# Compose | Docs | HERY
With compose the entities are put to getter and merged into a lock file that is then added to the cache 
[SQLite](https://www.sqlite.org/) database for quick querying. Compose starts with the main entity and from that entity 

Since the `collection.lock` is a JSON file it is possible to extract the data from the entities without using the caching 
system, but it will be slower.