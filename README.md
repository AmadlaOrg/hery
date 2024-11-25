<img src=".assets/bear.jpg" alt="Electronics photo" style="width: 400px;" align="right">

# hery 🐻
🐻 Hierarchical Entity Relational YAML (HERY) 🐻

HERY is an extension to [YAML](https://yaml.org/), leveraging the concept of entities—[YAML](https://yaml.org/)
groupings that can be interconnected similarly to an RDBMS. This CLI utility facilitates the use of HERY, enabling
efficient data organization and interaction.

Additionally, the term "hery" in British English, pronounced /ˈhɛrɪ/, is an obsolete verb meaning "to glorify; praise."
This name reflects the utility's aim to elevate and celebrate structured data management.

## 👐 Apropos

HERY differs with [YAML](https://yaml.org/) only by the four "reserved" properties: `_meta`, `_entity`, `_id` and `_body`. In other
words any `.hery` can be read by any [YAML](https://yaml.org/) library or editor. HERY's reserved properties are there to organize
the content in a [YAML](https://yaml.org/) file into entities. The property `_meta` contains one or more entities that attach metadata
information to an entity to make it easier to query and organize them. Similar to HTML `<meta>` element. The `_entity`
is the URI of the entity being used. It also contains the version of the entity. The `_id` is to link a property to
another entity similar to relationship in a RDBMS database. The `_body` contains the entity data, similar to HTML
`<body>` element.

Entities also require a [JSON-Schema](https://json-schema.org/) to define the standard for an entity. Everytime an
entity is added it validates it against its own schema or any other entity schemas. In the context of HTML it would the
[DTD](https://en.wikipedia.org/wiki/Document_type_definition).

HERY is also similar to a package manager whereby the entity's that are required can be added via the CLI or inside the
`.hery` file configuration using the reserved property `_entity` that contains the entity URI
(e.g.: `github.com/AmadlaOrg/Entity@latest`).

Once that the entities are added to the filesystem, it is also added to a cache system that is an
[SQLite3](https://www.sqlite.org/) database. From this cache system HERY can query the different entities. HERY comes with
its own query language. Everything that is output will be in [JSON](https://www.json.org/) format that can then be used
with [jq](https://jqlang.github.io/jq/) or any other tools that support [JSON](https://www.json.org/).

An entity can overwrite another entity of the same type. Entities can be thought as a table and the data set in the `_body`
as row in a RDBMS database. To overwrite a specific entity content ("row"), the `_id` property needs to be used.

Entities are grouped by collections. A collection can be thought as a database.

✅ A simple definition list:
- 🛒 **Collection** -> Database
- 📦 **Entity** -> Table
- 🎁 **Entity content** -> Row
- 🔖 **Meta** -> HTML `<meta>`
- 🪪 **Id** -> An `id` for a row in a RDBMS database

To have an entity it needs to be in a repository that uses [Git](https://git-scm.com/). At the root it needs a file
that is named after the collection and with the extension `.hery`. This means that it is possible to have multiple
collections in the same repository using `.hery` files. Inside is the definition of the entity.

The other component required is the collection directory. It is named after the collection with a dot at the beginning.
For example: `.amadla/`. Inside is the `schema.hery.json` file that is a [JSON-Schema](https://json-schema.org/) definition.
The directory can also contain any files that an entity might need.

## Amadla 🐰 ❤️ HERY 🐻
Amadla uses HERY as storage system. All the other tools in the Amadla ecosystem uses it for storage. It is also possible
to use custom tools since HERY can easily be used as a library in a custom [Golang](https://go.dev/) project. Or the
[JSON](https://www.json.org/) output by HERY can be piped into anything.

Amadla tries to follow the [UNIX philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) as much as possible.

## :suspect: Why Not Just Use SQLite 🐘?
- Entities concept with [YAML](https://yaml.org/) is simpler to use
- Lower learning curve
- It manages the download of separate entities automatically
- Easier to read
- Takes advantages of CVS like Git
- It easy to attach metadata to entity ("table")
- Easier validation

It can be thought as an abstraction of a RDBMS and a Package Manager.

## 🏎 How Fast Is It?
For the downloading of entities it will depend on how heavy the repository is. But generally an entity is just text so
should be quick.

For the query of data via HERY, it should be pretty quick since it uses [SQLite3](https://www.sqlite.org/) in the
backend.

## 📥 Install
### 🐹 With Go
```bash
go install github.com/AmadlaOrg/hery
```
### 🔨 Build
```bash
go build -o hery
```

## 🚀 Quickstart
HERY does not require a lot of learning to get started. All you need to know is the four reserved properties, the
`.hery` file format that the reserved properties are found, that there is a SQLite caching system, a few of the
commands and understand some of the basics of the query language.

### 📑 `.hery` File Format
The `.hery` file format is the same as a `.yml`/`.yaml` file format. The reason the extension is different is so that
`hery` CLI is able to find it and so that IDEs can have better support.

HERY format comes with four different reserved properties:
- `_meta` - For metadata for the relative entity
- `_entity` - Contains the entity URI with the version
- `_id` - To be able to make reference to a specific entity content
- `_body` - Contains the content of the entity

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

### 🚛 Caching
Since the querying on [YAML](https://yaml.org/) files would be a bit slow and resource demanding, [SQLite3](https://www.sqlite.org/)
is used to store all the entities.

Each entity has its own table, and it is found: `~/.hery/collection/<collection name>/<collection name>.cache`.

### 🖥️ Basic Commands
To verify that it was installed properly:
```bash
hery --version
```

### 🔩 Create a collection
```bash
hery collection init collection_name
```

### 📃 List collections
```bash
hery collection list
```

### 📥 Download an entity
```bash
# Without version it will get the latest version or if no version found then it will generate a pseudo version number using the commit hash
# To add a version: @v{version}
hery entity --collection amadla get github.com/Repository/EntityName
```

### 🔍 Query
```bash
hery entity --collection {collection name} query 'entities'
```

### ➕ More...
To get more details on the functioning and commands: [.docs](.docs).

## ⌨ Dev
### 🔥 Developer Benefits
A developer should find the code of this project to be well organized. It also comes with generated mocks that can make
it very easy to write unit tests without needing to make mocks. It also comes with interfaces for each packages making
it easy to overwrite. The only functions that are not in the interface are generally just simple utility functions.

There are also a lot of text and code editor plugins that make it a breeze for the `.hery` files to be used.

> PRs are always welcome!

### 📝 IDE Plugins
       ![Vim icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/vim.png) [Vim](.editor/.vimrc)

       ![Code icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/vscode.png) [Visual Studio Code](.editor/code.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-code-editor-plugin))

       ![IntelliJ icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/IntelliJ_IDEA.png) [JetBrains](.editor/jetbrains.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-jetbrains-editor-plugin))

       ![Sublime Text icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/sublime.png) [Sublime Text](.editor/sublime.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-sublime-editor-plugin))

       ![GNU Emacs icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/Emacs.png) [GNU Emacs](.editor/emacs.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-emacs-editor-plugin))

## ©️ Copyright
- "[The Bear and Honey.](https://www.flickr.com/photos/97123293@N07/29003630251)" by [Swallowtail Garden Seeds](https://www.flickr.com/photos/97123293@N07) is marked with [Public Domain Mark 1.0](https://creativecommons.org/publicdomain/mark/1.0/?ref=openverse).

## :scroll: License

The license for the code and documentation can be found in the [LICENSE](./LICENSE) file.

---

Made in Québec 🏴󠁣󠁡󠁱󠁣󠁿, Canada 🇨🇦!
