<img src=".assets/bear.jpg" alt="Electronics photo" style="width: 400px;" align="right">

# hery 🐻
🐻 Hierarchical Entity Relational YAML (HERY) 🐻

Is an extension to YAML that uses the concept of entities that are YAML groupings that can be interconnected similarly 
to an RDBMS.

This cli utility makes it possible to use HERY.

## Install
### Build
```bash
go build -o hery
```

## Quickstart
To verify that it was installed properly: 
```bash
hery --version
```

### Create a collection
```bash
hery collection init collection_name
```

### List collections
```bash
hery collection list
```

### Download an entity
```bash
# Without version it will get the latest version or if no version found then it will generate a pseudo version number using the commit hash
# To add a version: @v{version}
hery entity --collection amadla get github.com/Repository/EntityName
```

### Query
```bash
hery entity --collection {collection name} query 'entities'
```

### More...
To get more details on the functioning and commands: [.docs](.docs).

## Dev
### Check dependencies
Useful to check that libraries for testing and development don't find themselves in the build version of the project: 
```bash
go list -f '{{.Deps}}' ./main.go
```

> TODO: Create a script to blacklist certain libraries.
> Also useful for security.

## ©️ Copyright
- "[The Bear and Honey.](https://www.flickr.com/photos/97123293@N07/29003630251)" by [Swallowtail Garden Seeds](https://www.flickr.com/photos/97123293@N07) is marked with [Public Domain Mark 1.0](https://creativecommons.org/publicdomain/mark/1.0/?ref=openverse).

## :scroll: License

The license for the code and documentation can be found in the [LICENSE](./LICENSE) file.

---

Made in Québec :fleur_de_lis:, Canada 🇨🇦!