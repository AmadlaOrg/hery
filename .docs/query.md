# Query | Docs | HERY
[`jq`](https://jqlang.github.io/jq/) and [`JSONPATH`](https://jsonpath.com/) are great utilities to extract information from
json content. But with the added dimensions of relational entities it becomes a bit too limiting and hard to only use 
these tools.

So for that reason HERY comes with its own query language that is even simpler than SQL or [`jq`](https://jqlang.github.io/jq/) and [`JSONPATH`](https://jsonpath.com/).

## Selecting all entities
```heryquery
entities
```
Yes, that simple!

## Selecting by `_entity` property
```heryquery
entities._entity("github.com/AmadlaOrg/EntityApplication@v*")
```

Again, yes, that simple!

> In `("")` it is possible to use [glob](https://en.wikipedia.org/wiki/Glob_(programming)).

## Selecting by any entity property
```heryquery
entities._entity("github.com/AmadlaOrg/Entity@v*").name("EntityWebserver")
```
This will return all the entity web servers.

```heryquery
entities._entity("github.com/AmadlaOrg/Entity@v*").name("EntityWebserver").server_name("example.com")
```
This one will return the web server with server name *example.com*.

## Piping?
In UNIX systems to pass data to the next command you pipe with `|` but since this might cause conflicts in the command line 
the piping character for HERY is: `!`. Just like [gstreamer](https://gstreamer.freedesktop.org/).

## When is it possible to use `jq` and `JSONPATH`?
After selecting one or many entities it is then possible to use `jq` and `JSONPATH`. Since the output from HERY queries
are always in JSON format. It is also possible with a flag to indicate to output everything in YAML.

```heryquery
entities._entity("github.com/AmadlaOrg/Entity@v*").name("EntityWebserver") ! jq(".[0]") ! jsonpath("$.listen")
```

The documentation for [`jq`](https://jqlang.github.io/jq/) and [`JSONPATH`](https://jsonpath.com/) are found on their 
respected websites. For HERY queries there isn't more to it than what was shown.