# Version | Docs | HERY
The version numbering follows the same as [Golang's](https://go.dev/doc/modules/version-numbers) standard.

The versions are set in the tags in a Git repo. If no version is found then a pseudo-version number will be generated
just like in [Golang](https://go.dev/doc/modules/version-numbers#pseudo-version-number).

## Entity version numbering
HERY follows the same [versioning standard as Golang](https://go.dev/doc/modules/version-numbers).

| Version stage                         | Example                                                              | Description                                                                                                                                                       |
|---------------------------------------|----------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| In development                        | Automatic pseudo-version number `v0.0.0-20170915032832-14c0d48ead0c` | Indicates that the module is currently **in development and is considered unstable**. This release does not offer backward compatibility or stability assurances. |
| Major version                         | `v1.x.x`                                                             | Indicates **public API changes that are not backward-compatible**. This release does not guarantee compatibility with previous major versions.                    |
| Minor version                         | `vx.4.x`                                                             | Indicates **public API changes that are backward-compatible**. This release guarantees both backward compatibility and stability.                                 |
| Patch version                         | `vx.x.1`                                                             | Indicates **changes that do not impact the module's public API or its dependencies**. This release ensures backward compatibility and stability.                  |
| Pre-release version (alpha, beta, rc) | `vx.x.x-beta.2`                                                      | Indicates that this is a **pre-release milestone, such as an alpha or beta version**. This release does not offer any stability guarantees.                       |

## Manage Tag Versioning
Steps to add a tag:
```bash
git tag -a v1.0.0-alpha.1
git push --tags
```
