# Storage | Docs | HERY

## Project Level

A project's entity content lives at the project root:

```
my-project/
  webserver.hery           # Entity instances (source of truth, committed)
  network.hery             # Entity instances
  database.hery            # Entity instances
  hery.lock                # Lock file (JSON, committed to Git)
  .hery.cache              # SQLite cache (gitignored, rebuilt from .hery files)
```

- **`*.hery`** -- entity content files (committed)
- **`hery.lock`** -- portable JSON snapshot of the merged state (committed). Same data as .hery.cache but in a portable, debuggable format.
- **`.hery.cache`** -- SQLite database for fast queries (gitignored, derived)

## Global Entity Cache

Resolved entity types are cached globally at `~/.cache/hery/entity/`, organized by git path:

```
~/.cache/hery/
  entity/
    amadla.org/
      entity/
        application@v1.0.0/
          application.hery.json
          default.hery
        network@v1.0.0/
          network.hery.json
          default.hery
    github.com/
      AmadlaOrg/
        EntityApplication@v1.0.0/
          application.hery.json
          default.hery
      jnbdz/
        personal-website/        # Cloned for _extends resolution
          webserver/
            webserver.hery
```

The cache is disposable -- hery re-fetches from Git when needed.

## Global Config (optional)

```
~/.config/hery/
  config.yaml              # Global hery settings
```
