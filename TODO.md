# TODO
- [ ] Making pseudo versioning portable (what??? will making caching a higher authority) meaning no more date time
- [ ] Pseudo version first part (v0.0.0) should this following the previous tagging version
- [ ] In the pipeline add a check if the GitHub of library is in archive... Maybe throw a warning
- [ ] In the pipeline check the license and if they are compatible (throw a warning if none found or ot open source)
- [ ] Check if it is possible to generate different binaries (client, and server or both)

## Dev
- [ ] Parsers - 10h
- [ ] Schema - Change the "name" and switch the datatype formating place from parser to schema package
- [ ] Database - 3h
- [ ] Caching - 5h
- [ ] Collection - 10h
- [ ] Entity - 10h
- [ ] Query - 10h
- [ ] Documentation - 10h
- [ ] BDD testing 12h

~ 70h (~ 18 days)

## Notes
### Check dependencies
This would be good to have this in the pipeline.

Useful to check that libraries for testing and development don't find themselves in the build version of the project:
```bash
go list -f '{{.Deps}}' ./main.go
```

> TODO: Create a script to blacklist certain libraries.
> Also useful for security.
