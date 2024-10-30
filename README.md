<img src=".assets/bear.jpg" alt="Electronics photo" style="width: 400px;" align="right">

# hery 🐻
🐻 Hierarchical Entity Relational YAML (HERY) 🐻

HERY is an extension to [YAML](https://yaml.org/), leveraging the concept of entities—[YAML](https://yaml.org/)
groupings that can be interconnected similarly to an RDBMS. This CLI utility facilitates the use of HERY, enabling
efficient data organization and interaction.

Additionally, the term "hery" in British English, pronounced /ˈhɛrɪ/, is an obsolete verb meaning "to glorify; praise."
This name reflects the utility's aim to elevate and celebrate structured data management.

## Install
### With Go
```bash
go install github.com/AmadlaOrg/hery
```
### Build
```bash
go build -o hery
```

## Quickstart
HERY does not require a lot of learning to get started. All you need to know is the four reserved properties, the
`.hery` file format that the reserved properties are found, that there is a SQLite caching system, a few of the
commands and understand some of the basics of the query language.

### `.hery` File Format
The `.hery` file format is the same as a `.yml`/`.yaml` file format. The reason the extension is different is so that
`hery` CLI is able to find it and so that IDEs can have better support.

HERY format comes with four different reserved properties:
- `_meta` - For metadata for the relative entity
- `_entity` - Contains the entity URI with the version (is optional)
- `_body` - Contains the content of the entity
- `_id` - To be able to make reference to a specific entity content

Here is an example:
```yaml
---
_meta:
  _entity: github.com/AmadlaOrg/Entity@latest
  _body:
    name: RandomName
    description: Entity Pseudo Version definitions.
    category: QA
    tags:
      - QA
      - fixture
      - test
_body:
  name: Random Name
```

There are different structures that are valid. For example `_meta` is an optional reserved property. And at the root of
the entity content `_entity` is not allowed since it takes the value from the repository itself.

When validation happens it takes for account the entity URI and validates what is in the `_body` whilst everything else
is ignored.

### Caching
Since the querying on YAML files would be a bit slow and resource demanding, SQLite is used to store all the entities.

Each entity has its own table, and it is found: `~/.hery/collection/<collection name>/<collection name>.cache`.

### Basic Commands
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
### IDE Plugins
       ![Vim icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/vim.png) [Vim](.editor/.vimrc)

       ![Code icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/vscode.png) [Visual Studio Code](.editor/code.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-code-editor-plugin))

       ![IntelliJ icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/IntelliJ_IDEA.png) [JetBrains](.editor/jetbrains.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-jetbrains-editor-plugin))

       ![Zed icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/zed.png) [Zed](.editor/zed.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-zed-editor-plugin))

       ![Sublime Text icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/sublime.png) [Sublime Text](.editor/sublime.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-sublime-editor-plugin))

       ![GNU Emacs icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/Emacs.png) [GNU Emacs](.editor/emacs.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-emacs-editor-plugin))

       ![Brackets icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/brackets.png) [Brackets](.editor/brackets.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-brackets-editor-plugin))

       ![Atom icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/atom.png) [Atom](.editor/atom.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-atom-editor-plugin))

### Giving Credit
The `.pairs` file is a great way to add details about the people that contributed to the entity. There is one used in
this repository that can be used as reference.

> TODO: Create a script to blacklist certain libraries.
> Also useful for security.

## ©️ Copyright
- "[The Bear and Honey.](https://www.flickr.com/photos/97123293@N07/29003630251)" by [Swallowtail Garden Seeds](https://www.flickr.com/photos/97123293@N07) is marked with [Public Domain Mark 1.0](https://creativecommons.org/publicdomain/mark/1.0/?ref=openverse).

## :scroll: License

The license for the code and documentation can be found in the [LICENSE](./LICENSE) file.

---

Made in Québec 🏴󠁣󠁡󠁱󠁣󠁿, Canada 🇨🇦!
