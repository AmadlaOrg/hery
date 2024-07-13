# Collection | Docs | HERY
A collection is a grouping of entities or entity type.

For example [Amadla](https://github.com/AmadlaOrg/) project has its own collection: `amadla`. [Amadla](https://github.com/AmadlaOrg/)
has its own entities. Each has its own structure to store information.

Having collections makes it possible to have other project have their own grouping of entities. So to not confused
different standards and projects.

## `hery`-cli collection utils
With the cli it is possible to list and add a collection.

## Where?
`~/.henry/` or `{project root}/.henry/`.

## How to set?
With the environment variable: `HERY_COLLECTION=amadla`

It is also possible to pass it as a command flag: `--collection=amadla` or `-c amadla`