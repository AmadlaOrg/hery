<img src=".assets/bear.jpg" alt="Electronics photo" style="width: 400px;" align="right">

# `hery`
Hierarchical Entity Relational YAML (HERY)

HERY is an extension to [YAML](https://yaml.org/), leveraging the concept of entities—[YAML](https://yaml.org/)
groupings that can be interconnected similarly to an RDBMS. This CLI utility facilitates the use of HERY, enabling
efficient data organization and interaction.

Additionally, the term "hery" in British English, pronounced /ˈhɛrɪ/, is an obsolete verb meaning "to glorify; praise."
This name reflects the utility's aim to elevate and celebrate structured data management.

## Apropos

HERY differs from [YAML](https://yaml.org/) only by five "reserved" properties: `_type`, `_self`, `_parent`, `_meta`
and `_body`. In other words any `.hery` file can be read by any [YAML](https://yaml.org/) library or editor.

HERY's reserved properties organize content in a [YAML](https://yaml.org/) file into entities:
- `_type` is the URI of the entity type being used, including the version
- `_self` is an optional resolvable identifier for a specific entity content instance (merge discriminator)
- `_parent` is an optional URI to a parent entity instance, enabling deep merge inheritance
- `_meta` contains metadata for the entity, making it easier to query and organize (similar to HTML `<meta>`)
- `_body` contains the entity data (similar to HTML `<body>`)

Entities require a [JSON-Schema](https://json-schema.org/) (`schema.hery.json`) to define the standard for an entity.
When an entity is added, it is validated against its schema.

HERY is also similar to a package manager whereby entities can be added via the CLI or inside the
`.hery` file using `_type` with an entity URI (e.g.: `github.com/AmadlaOrg/Entity@latest`).

Once entities are added to the filesystem, they are cached in an [SQLite3](https://www.sqlite.org/) database.
HERY uses a two-stage query model: selection via CLI flags hitting SQLite indexes, then optional transformation
via [jq](https://jqlang.github.io/jq/) expressions (compiled in via gojq). Output is always [JSON](https://www.json.org/).

Entities support deep merge inheritance via `_parent`: child values override parent values — objects merge recursively,
arrays replace entirely, scalars are overridden by the child.

A simple definition parallel:

| Component             | Parallel                               |
|-----------------------|----------------------------------------|
| **Entity**            | Table                                  |
| **Entity content**    | Row                                    |
| **Meta**              | HTML `<meta>`                          |
| **Self**              | Row identifier / merge discriminator   |

To have an entity it needs to be in a repository that uses [Git](https://git-scm.com/). At the root it needs a
`schema.hery.json` file and one or more `.hery` content files.

## Amadla + HERY
Amadla ecosystem follows as best as possible the [UNIX philosophy](https://en.wikipedia.org/wiki/Unix_philosophy). So any storage sources that can `stdout` will
work. HERY is an optional storage source that is chiefly recommended for the Amadla ecosystem.

It is also possible to use HERY as a library in a custom [Golang](https://go.dev/) project. Or the [JSON](https://www.json.org/) output by HERY can be
piped.

## Why Not Just Use SQLite?
- Entities concept with [YAML](https://yaml.org/) is simpler to use
- Lower learning curve
- It manages the download of separate entities automatically
- Easier to read
- Takes advantages of VCS like Git
- Easy to attach metadata to entities
- Easier validation

It can be thought as an abstraction of a RDBMS and a Package Manager.

## How Fast Is It?
For the downloading of entities it will depend on how heavy the repository is. But generally an entity is just text so
should be quick.

For the query of data via HERY, it should be pretty quick since it uses [SQLite3](https://www.sqlite.org/) in the
backend.

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
HERY does not require a lot of learning to get started. All you need to know is the five reserved properties, the
`.hery` file format, that there is a SQLite caching system, a few of the commands, and the two-stage query model.

### `.hery` File Format
The `.hery` file format is the same as a `.yml`/`.yaml` file format. The reason the extension is different is so that
the `hery` CLI can find it and so that IDEs can have better support.

HERY format has five reserved properties:

| Property   | Description                                                   |
|------------|---------------------------------------------------------------|
| `_type`    | Entity type URI with version (required)                       |
| `_self`    | Resolvable identifier for entity content (optional)           |
| `_parent`  | URI to parent entity for deep merge inheritance (optional)    |
| `_meta`    | Metadata for the entity (optional)                            |
| `_body`    | Contains the content of the entity (optional)                 |

Here is an example:
```yaml
---
_type: github.com/AmadlaOrg/EntityQA/RandomName@latest
_meta:
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

When validation happens it validates what is in the `_body` against the entity's JSON Schema.

### Caching
Since querying [YAML](https://yaml.org/) files directly would be slow, [SQLite3](https://www.sqlite.org/)
is used to cache all the entities at: `~/.cache/hery/entity/`.

### Basic Commands
To verify that it was installed properly:
```bash
hery --version
```

### Download an entity
```bash
# Without version it will get the latest version or generate a pseudo version
# To add a version: @v{version}
hery entity get github.com/Repository/EntityName
```

### Query
```bash
# Two-stage query: selection flags + optional jq transformation
hery query --type "github.com/AmadlaOrg/Entity@latest"
hery query --type "github.com/AmadlaOrg/Entity@latest" --jq '.[] | .name'
```

### More...
To get more details on the functioning and commands: [.docs](.docs).

## Dev
### Developer Benefits
A developer should find the code of this project to be well organized. It also comes with generated mocks that can make
it very easy to write unit tests without needing to make mocks. It also comes with interfaces for each package making
it easy to overwrite.

> PRs are always welcome!

### IDE Plugins
       ![Vim icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/vim.png) [Vim](.editor/.vimrc)

       ![Code icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/vscode.png) [Visual Studio Code](.editor/code.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-code-editor-plugin))

       ![IntelliJ icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/IntelliJ_IDEA.png) [JetBrains](.editor/jetbrains.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-jetbrains-editor-plugin))

       ![Sublime Text icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/sublime.png) [Sublime Text](.editor/sublime.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-sublime-editor-plugin))

       ![GNU Emacs icon](https://raw.githubusercontent.com/SiteNetSoft/resources/master/images/ide/x14/Emacs.png) [GNU Emacs](.editor/emacs.yml) - ([GitHub](https://github.com/AmadlaOrg/hery-emacs-editor-plugin))

## Copyright
- "[The Bear and Honey.](https://www.flickr.com/photos/97123293@N07/29003630251)" by [Swallowtail Garden Seeds](https://www.flickr.com/photos/97123293@N07) is marked with [Public Domain Mark 1.0](https://creativecommons.org/publicdomain/mark/1.0/?ref=openverse).

## License

The license for the code and documentation can be found in the [LICENSE](./LICENSE) file.

---

Made in Québec, Canada!
