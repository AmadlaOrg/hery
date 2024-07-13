# Cache | Docs | HERY
To be able to have an efficient querying system caching is essentials.

The command compose puts the different entities together to then generate a lock file. Using this lock file and the 
multiple [JSON Schemas](https://json-schema.org/) it generates an [SQLite](https://www.sqlite.org/) DB that is saved
to the root of the collection's storage directory (e.g.: `amadla.cache`).

> When `hery`-cli is run in server mode it will load the DB in memory for quicker reads. Also, to note server mode is 
> not for writing. Creation and updating of the cache is done outside the server mode.

Here is an example what it looks like inside the cache DB: 

![Example of a cache DB structure](./diagram/caching-example.svg)

> `hery`-cli does not support other DBs, and it comes with its own [SQLite](https://github.com/mattn/go-sqlite3).