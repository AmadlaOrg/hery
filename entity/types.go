package entity

import "github.com/santhosh-tekuri/jsonschema/v6"

const (
	// The entity name format
	// TODO: It also does not seemed to be used
	EntityNameMatch = `^[a-zA-Z0-9]+$`

	// The entity name and version format
	// @deprecated: the version is wrong since it changed
	// TODO: It also does not seemed to be used
	EntityNameAndVersionMatch = `^([a-zA-Z0-9]+)(@v\d+\.\d+\.\d+)$`

	// Used to identify the entities that are stored in a collection
	FormatEntityPathAndNameVersion = `^(.+[/\/]\.[A-z0-9-_]+[/\/]entity[/\/])(.+)([/\/].+@).+$`
)

// Entity holds the origin and version of an entity
//
// There are multiple attributes that are attached to an Entity. They are used for selecting and working with entities.
type Entity struct {
	Id              string             // uuid (e.g.: 97d4b783-f448-483c-8111-380d6082ae1c)
	Entity          string             // Reserved (e.g.: github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0)
	Name            string             // The simple name of an entity (e.g.: WebServer)
	RepoUrl         string             // The full repository URL (e.g.: https://github.com/AmadlaOrg/EntityApplication)
	Origin          string             // The entity URL path (it can also be used as a relative path) (e.g.: github.com/AmadlaOrg/EntityApplication)
	Version         string             // The entity version (what is after `@`) (e.g.: v1.0.0)
	IsPseudoVersion bool               // True if the version was generated
	AbsPath         string             // The absolute path to where the entity is stored (e.g.: /home/user/.hery/amadla/entity/github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0)
	Have            bool               // True if the entity is downloaded and false if not (e.g.: true)
	Hash            string             // The hash of the entity to verify if the repository on the local environment was corrupted or not (e.g.: c7e9911d38b263a69c664b8e0b5d4f27e607554d)
	Exist           bool               // True if it was found and false if not found with Git remote (e.g.: true)
	Schema          *jsonschema.Schema // The entity's JSON-Schema
	Config          map[string]any     // From the `.hery` config file
}
