package entity

const (
	// The entity name format
	entityNameMatch = `^[a-zA-Z0-9]+$`

	// The entity name and version format
	// @deprecated: the version is wrong since it changed
	entityNameAndVersionMatch = `^([a-zA-Z0-9]+)(@v\d+\.\d+\.\d+)$`

	// Used to identify the entities that are stored in a collection
	formatEntityPathAndNameVersion = `^(.+[/\/]\.[A-z0-9-_]+[/\/]entity[/\/])(.+)([/\/].+@).+$`
)

// Entity holds the origin and version of an entity
//
// There are multiple attributes that are attached to an Entity. They are used for selecting and working with entities.
type Entity struct {
	Name    string // The simple name of an entity
	Uri     string //
	Origin  string // The entity URL path (it can also be used as a relative path)
	Version string // The entity version (what is after `@`)
	AbsPath string // The absolute path to where the entity is stored
	Have    bool   // True if the entity is downloaded and false if not
	Hash    string // The hash of the entity to verify if the repository on the local environment was corrupted or not
	Exist   bool   // True if it was found and false if not found
}
