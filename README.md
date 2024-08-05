<img src=".assets/bear.jpg" alt="Electronics photo" style="width: 400px;" align="right">

# hery 🐻
🐻 Hierarchical Entity Relational YAML (HERY) 🐻

HERY is an extension to [YAML](https://yaml.org/), leveraging the concept of entities—[YAML](https://yaml.org/) 
groupings that can be interconnected similarly to an RDBMS. This CLI utility facilitates the use of HERY, enabling 
efficient data organization and interaction.

Additionally, the term "hery" in British English, pronounced /ˈhɛrɪ/, is an obsolete verb meaning "to glorify; praise." 
This name reflects the utility's aim to elevate and celebrate structured data management.

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
### IDEs
- [Vim](.editor/.vimrc)
- [Visual Studio Code](.editor/code.yml)
- [JetBrains](.editor/jetbrains.yml)
- [Zed](.editor/zed.yml)
- [Sublime Text](.editor/sublime.yml)
- [GNU Emacs](.editor/emacs.yml)
- [Brackets](.editor/brackets.yml)
- [Atom](.editor/atom.yml)

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