# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## AI Skills

Follow the practices defined in `~/Projects/SiteNetSoft/ai-skills/`:
- `dev-practices/golang/` — Go style, error handling, functions, testing, linting
- `dev-practices/git/` — Git authorship rules, multi-repo workspace patterns

## Project Overview

HERY (Hierarchical Entity Relational YAML) is a Go CLI utility that extends YAML with entity management capabilities. It combines concepts from RDBMS databases and package managers, enabling structured data organization with schema validation, Git-based versioning, and SQLite caching.

**Core concept (Draft 3.5):** HERY adds five reserved YAML properties to organize content into entities:
- `_type` - Entity type URI with version (e.g., `github.com/AmadlaOrg/Application@latest`)
- `_extends` - Optional extended entity URI (enables deep merge inheritance, purely data/merge — no execution ordering)
- `_meta` - Metadata for the entity (like HTML `<meta>`)
- `_body` - Contains entity data (like HTML `<body>`)
- `_requires` - Optional hard dependencies on other entities for execution ordering (amadla builds DAG, topological sort)

Entity identity is derived from the git path (directory position in repo), not declared in the document.

**Terminology parallel:**
| HERY Component | RDBMS Equivalent |
|----------------|------------------|
| Entity | Table |
| Entity content | Row |

**Note:** Collections have been removed in Draft 3.2. All entities share a flat global namespace under `~/.cache/hery/entity/`.

## Build Commands

```bash
make build              # Build for Linux amd64 (default)
make build-all          # Build for all platforms
make test               # Run tests with coverage
make lint               # Run golangci-lint
make lint-fix           # Auto-fix linting issues
make generate           # Generate mocks via mockery
make install-deps       # Install dev tools (golangci-lint, mockery, etc.)
make clean              # Clean build artifacts
```

Output binaries go to `bin/` directory.

## Architecture

### Package Structure

```
cmd/                    # Cobra CLI commands (entity, query, compose, settings)
entity/                 # Core entity logic
├── build/              # Entity building
├── cmd/                # Entity subcommands (get, list, validate)
├── compose/            # Multi-entity composition
├── get/                # Entity retrieval
├── merge/              # Deep merge engine (Draft 3.2)
├── query/              # Two-stage query engine (selection + jq)
├── schema/             # JSON Schema handling
├── validation/         # Entity validation
└── version/            # Version management
cache/                  # SQLite caching layer
├── database/           # SQLite operations
└── parser/             # Cache parsing
storage/                # Filesystem abstraction (Storage, Entities, Cache paths)
message/                # Error types
```

### Key Patterns

**Interface-based design:** Every package defines interfaces using idiomatic Go naming:
- Interfaces use concept names: `Storage`, `Cache`, `Service`, `Builder`, `Getter`, `Schema`, `Validator`, `Version`, etc.
- Implementation structs are unexported (e.g., `service`, `builder`, `cacheImpl`, `dbImpl`)
- Constructor pattern: `New()` (one primary type per package)

**Mock generation:** Uses mockery v2 configured in `.mockery.yaml`. Mocks are generated in-package with `mock_*.go` naming. Run `make generate` after interface changes. Mockery binary is at `/home/jn/go/bin/mockery`.

**Testing:** Uses testify/assert and testify/mock. Table-driven tests are standard. E2e tests use Ginkgo v2 + Gomega.

### Dependencies

The project uses local replace directives for sibling libraries:
```
github.com/AmadlaOrg/LibraryUtils => ../LibraryUtils
github.com/AmadlaOrg/LibraryFramework => ../LibraryFramework
```

Key dependencies: Cobra (CLI), github.com/goccy/go-yaml (YAML), jsonschema/v6 (validation), go-sqlite3 (caching), gojq (queries).

### CLI Structure

Entry point in `main.go` uses `cli.New()` from LibraryFramework. Commands registered:
- `hery entity` - Entity operations (get, list, validate)
- `hery query` - Query entities (two-stage: selection flags + `--jq` transformation). Default source: the current directory's `.hery` files, served from the project-level `.hery.cache` when fresh; or a file/stdin with `--from <file>` / `-f -`. Output: `-o table|json|yaml` (default table, like the rest of the suite — pipelines pass `-o json`), `--hery` to wrap in a HERY envelope.
- `hery compose` - Compose multiple entities
- `hery settings` - Configuration

### Query Model (Draft 3.2)

Two-stage query:
1. **Selection** — CLI flags (`--type`, `--meta`, `--tag`)
2. **Transformation** — `--jq` flag applies jq expressions via gojq (compiled in, no external binary)

**Data sources (mutually exclusive):**
- Default — the current directory's `.hery` files (the source of truth), resolved like `--dir '.'` and served from the project-level `.hery.cache` (SQLite, gitignored, derived). The cache carries a manifest of its sources (file mtime+size, per-directory `.hery` filename listing); any mismatch triggers a transparent re-resolve and rebuild. Cache writes are best-effort (a failure warns and falls back to in-memory), an empty project writes no cache file, and selection always runs through the same in-memory predicates as the other sources — identical semantics everywhere.
- `--from <file>` / `-f -` — read entities from a file or stdin instead of the cache. Format is auto-detected (YAML multi-doc stream, YAML single doc, JSON array, JSON object, or NDJSON). This is the consumer side of `hery compose --dir | hery query --from -` (the Option F pipeline: plugins query the composed graph without flattening). For file/dir input, `--type` matches the doc's `_type` (version-stripped unless the pattern carries `@version`); `*` matches across `/`, matching is case-insensitive; `--meta`/`--tag` substring-match the JSON-encoded `_meta`.
- `--dir <dir>` — convenience flag: resolves a directory of `.hery` files (`_extends` merge, `_requires` ordering, layering) exactly like `compose --dir`, then queries the resolved graph in memory — no shell pipe needed. Semantically identical to `hery compose --dir <dir> | hery query --from -`. Mutually exclusive with `--from`.

**Output:** `-o table|json|yaml` (default `table`, matching the rest of the suite; pipelines/plugins pass `-o json`). `--hery` wraps results in a HERY envelope (`amadla.org/entity/query/result@v1.0.0`, json/yaml only — with default table it promotes to json, with explicit `-o table` it errors). Exit codes: `0` = results, `2` = no match (grep-style), `1` = error (reported to stderr).

**Note:** `compose --dir` emits a `---`-separated multi-doc YAML stream so its output round-trips through `query --from`.

## File Locations

- Schema definitions: `.schema/`
- Documentation: `.docs/`
- SQL resources: `resources/sql/`
- Coverage reports: `.reports/`
- Linter config: `.golangci.yml` (excludes `mock_*.go` files)
- Test fixtures: `test/fixture/`
