# Query | Docs | HERY
[`jq`](https://jqlang.github.io/jq/) and [`JSONPATH`](https://jsonpath.com/) are great utilities to extract information from json content. But with the added
dimensions of relational entities it becomes a bit too limiting and hard to only use these tools.

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

## Examples
### Get all the entities
```eql
entities
```

JSON:
```json
[
  {
    "_id": "1d286661-b922-4bf4-a317-5530abfe57d5",
    "_entity": "github.com/AmadlaOrg/Entity@v1.0.0",
    "version": "v1.0.0"
  },
  {
    "_id": "5b9efb38-be80-4eeb-9167-1747b52ca72b",
    "_entity": "github.com/AmadlaOrg/EntityApplication@v1.0.0",
    "version": "v1.0.0"
  }
]
```

Table:
```
+--------------------------------------+-----------------------------------------------+---------+
| _id                                  | _entity                                       | Version |
+--------------------------------------+-----------------------------------------------+---------+
| 1d286661-b922-4bf4-a317-5530abfe57d5 | github.com/AmadlaOrg/Entity@v1.0.0            | v1.0.0  |
| 5b9efb38-be80-4eeb-9167-1747b52ca72b | github.com/AmadlaOrg/EntityApplication@v1.0.0 | v1.0.0  |
+--------------------------------------+-----------------------------------------------+---------+
```

### Get entities property
```eql
entities.property("_id")
```

JSON:
```json
[
  {
    "_id": "1d286661-b922-4bf4-a317-5530abfe57d5"
  },
  {
    "_id": "5b9efb38-be80-4eeb-9167-1747b52ca72b"
  }
]
```

Table:
```
+--------------------------------------+
| _id                                  |
+--------------------------------------+
| 1d286661-b922-4bf4-a317-5530abfe57d5 |
| 5b9efb38-be80-4eeb-9167-1747b52ca72b |
+--------------------------------------+
```

### Get entities multiple property
```eql
entities.property("_id", "Version")
```

JSON:
```json
[
  {
    "_id": "1d286661-b922-4bf4-a317-5530abfe57d5",
    "version": "v1.0.0"
  },
  {
    "_id": "5b9efb38-be80-4eeb-9167-1747b52ca72b",
    "version": "v1.0.0"
  }
]
```

Table:
```
+--------------------------------------+---------+
| _id                                  | Version |
+--------------------------------------+---------+
| 1d286661-b922-4bf4-a317-5530abfe57d5 | v1.0.0  |
| 5b9efb38-be80-4eeb-9167-1747b52ca72b | v1.0.0  |
+--------------------------------------+---------+
```

### Get entity by filter
```eql
entities.contains("_id", "5530")
```

JSON:
```json
[
  {
    "_id": "1d286661-b922-4bf4-a317-5530abfe57d5",
    "_entity": "github.com/AmadlaOrg/Entity@v1.0.0",
    "version": "v1.0.0"
  }
]
```

Table:
```
+--------------------------------------+-----------------------------------------------+---------+
| _id                                  | _entity                                       | Version |
+--------------------------------------+-----------------------------------------------+---------+
| 1d286661-b922-4bf4-a317-5530abfe57d5 | github.com/AmadlaOrg/Entity@v1.0.0            | v1.0.0  |
+--------------------------------------+-----------------------------------------------+---------+
```

```eql
entities.and(contains("_id", "17"), contains("_entity", "Entity@v1.0.0"))
```

```eql
entities.or(contains("_id", "17"), contains("_entity", "Entity@v1.0.0"))
```

```eql
entities.contains("_entity", "Entity@v1.0.0")._body
```

JSON:

```json
[
  {
    "server_name": "amadla.com",
    "directory": "/var/www/",
    "listen": [
      {
          "ports": [80, 443]
      }
    ]
  }
]
```


### Function

| Function                                                    | Description                                                                                                                           |
|-------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| equal(&lt;property name&gt;, &lt;filter with&gt;)           |                                                                                                                                       |
| like(&lt;property name&gt;, &lt;filter with&gt;)            | Used for pattern matching in string comparisons, often with wildcards (% for any sequence of characters and _ for a single character) |
| in(&lt;property name&gt;, &lt;values&gt;...)                | Used to filter rows where a column's value matches any value in a specified list                                                      |
| not_equal(&lt;property name&gt;, &lt;filter with&gt;)       |                                                                                                                                       |
| not_like(&lt;property name&gt;, &lt;filter with&gt;)        |                                                                                                                                       |
| not_in(&lt;property name&gt;, &lt;values&gt;...)            |                                                                                                                                       |
| and(&lt;functions&gt;...)                                   |                                                                                                                                       |
| or(&lt;functions&gt;...)                                    |                                                                                                                                       |
| limit(&lt;number&gt;)                                       | Limit specifies the maximum number of rows to return in a query                                                                       |
| offset(&lt;number&gt;)                                      | Offset specifies the number of rows to skip before starting to return results in a query                                              |
| order_by(&lt;property name&gt; [, &lt;ASC\|DESC&gt;])       | Using a property to order the output and optionally choose `DESC` or `ASC`                                                            |
| group_by(&lt;property name&gt; [, &lt;property name&gt...]) | Group a property or multiple-properties                                                                                               |
| jq()                                                        | Used to query JSON content with [jq](https://jqlang.github.io/jq/)                                                                    |
