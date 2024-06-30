package entity

const (
	entityNameMatch                = `^[a-zA-Z0-9]+$`
	entityNameAndVersionMatch      = `^([a-zA-Z0-9]+)(@v\d+\.\d+\.\d+)$`
	formatEntityPathAndNameVersion = `^(.+[/\/]\.[A-z0-9-_]+[/\/]entity[/\/])(.+)([/\/].+@).+$`
)

// Entity holds the origin and version of an entity
type Entity struct {
	Name    string
	Origin  string
	Version string
}
