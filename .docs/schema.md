# Schema | Docs | HERY
The `.hery` files are YAML files with reserved properties: `_entity`, `_id`, `_self`. So not to complicated. That said,
every entity has its own standard set of properties and data format. And those are standardized in a [JSON Schema](https://json-schema.org/)
file named: `schema.entity.json`. This schema file is found in the `.<collection name>/schema.entity.json` directory.

Inside this [JSON Schema](https://json-schema.org/) file there are standards that need to be followed.

1. The naming of the schema file.
  - It has the word `hery` in it for better IDE integration
  - An IDE plugin can quickly identify that it is part of a HERY entity schema definition since they vary from a normal [JSON Schema](https://json-schema.org/)
2. Inside the schema file there is the `id` property
   - It is required, compared to the [JSON Schema](https://json-schema.org/) standard
   - It uses a HERY standard URN that is prefixed with: `urn:hery:<collection name>:`
   - The URN is used so that the content of the schema can be identified as a HERY schema and so to not be confused with any other [JSON Schema](https://json-schema.org/)
   - This also helps with IDE plugins
   - The rest of the content is the same as `_entity` but `/` and `@` is replaced with `:`
3. The schema of an entity is always merged with `.schema/entity.schema.json` to make sure that when validation happens the standard HERY properties are also validated

> HERY is for developers and developers need the technology they use to have a good integration with their IDEs.
> So this is why some of these standard and rules exists. It adds to the responsibility of using HERY but makes working with it
> with developer tools a breeze.
