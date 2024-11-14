# TODO
- [ ] Making pseudo versioning portable (will making caching a higher authority)
- [ ] Pseudo version first part (v0.0.0) should this following the previous tagging version
- [ ] In the pipeline add a check if the Github of library is in archive... Maybe throw a warning
- [ ] In the pipeline check the license and if they are compatible (throw a warning if none found)
- [ ] Check if it is possible to generate different binaries (client, and server or both)

## Notes
### Check dependencies
This would be good to have this in the pipeline.

Useful to check that libraries for testing and development don't find themselves in the build version of the project:
```bash
go list -f '{{.Deps}}' ./main.go
```

> TODO: Create a script to blacklist certain libraries.
> Also useful for security.
