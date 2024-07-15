package entity

func BuildMeta(entityUri string) Entity {

	return Entity{
		Name:    "",
		Uri:     entityUri,
		Origin:  "",
		Version: "",
		AbsPath: "",
		Have:    false,
		Hash:    "",
		Exist:   true,
	}
}
