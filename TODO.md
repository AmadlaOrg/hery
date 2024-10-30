# TODO
- [ ] Making pseudo versioning portable (will making caching a higher authority)
- [ ] Pseudo version first part (v0.0.0) should this following the previous tagging version

## Notes
### Check dependencies
This would be good to have this in the pipeline.

Useful to check that libraries for testing and development don't find themselves in the build version of the project:
```bash
go list -f '{{.Deps}}' ./main.go
```

> TODO: Create a script to blacklist certain libraries.
> Also useful for security.
