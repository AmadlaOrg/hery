# Server | Docs | HERY
When running the server the [SQLite](https://www.sqlite.org/) database is loaded in memory. This makes it faster to 
access the cached entities for querying. The other benefit is the sending requests without having to call cli command 
making things more verbose.

The server is accessible via `/tmp/hery.sock` UNIX socket file.

## Start Server
```bash
hery server start --collection=amadla
```

## Stop Server
```bash
hery server stop --collection=amadla
```
OR
`Ctrl+c`

## Client Access
```bash
hery client --collection=amadla
```