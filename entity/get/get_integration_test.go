package get

/*func Test_Integration_GetInTmp(t *testing.T) {
	entities := []string{"github.com/AmadlaOrg/Entity"}
	entityBuild := NewGetService()
	paths, err := entityBuild.GetInTmp("amadla", entities)
	if err != nil {
		t.Fatal(err)
	}

	println(paths.Entities)
}*/

// FIXME:
/*func Test_Integration_Get(t *testing.T) {
tempDir, err := os.MkdirTemp("", "hery_test_*")
if err != nil {
	t.Fatal(err)
}

// Clean up after test
defer func(path string) {
	err = os.RemoveAll(path)
	if err != nil {
		t.Fatal(err)
	}
}(tempDir)

paths := storage.AbsPaths{
	Storage:    filepath.Join(tempDir, ".hery"),
	Catalog:    filepath.Join(tempDir, ".hery", "collection"),
	Collection: filepath.Join(tempDir, ".hery", "collection", "amadla"),
	Entities:   filepath.Join(tempDir, ".hery", "collection", "amadla", "entity"),
	Cache:      filepath.Join(tempDir, ".hery", "collection", "amadla", "amadla.cache"),
}

tests := []struct {
	name           string
	collectionName string
	paths          storage.AbsPaths
	entityURIs     []string
	collision      bool
	hasError       bool
}{
	/*{
		name:           "Get One",
		collectionName: "amadla",
		paths:          paths,
		entityURIs: []string{
			"github.com/AmadlaOrg/EntityApplication",
		},
	},
	{
		name:           "Get Multiple different URIs",
		collectionName: "amadla",
		paths:          paths,
		entityURIs: []string{
			"github.com/AmadlaOrg/Entity",
			"github.com/AmadlaOrg/EntityApplication",
		},
	},
	{
		name:           "Get Multiple identical URIs (pseudo versions)",
		collectionName: "amadla",
		paths:          paths,
		entityURIs: []string{
			"github.com/AmadlaOrg/Entity",
			"github.com/AmadlaOrg/Entity",
		},
	},
	{
		name:           "Get Multiple identical URIs (static versions)",
		collectionName: "amadla",
		paths:          paths,
		entityURIs: []string{
			"github.com/AmadlaOrg/Entity@v1.0.0",
			"github.com/AmadlaOrg/Entity@v1.0.0",
		},
	},*/
/*{
	name:           "Get Multiple different URIs (with none-existing version for QAFixturesEntityPseudoVersion)",
	collectionName: "amadla",
	paths:          paths,
	entityURIs: []string{
		"github.com/AmadlaOrg/QAFixturesEntityPseudoVersion@v1.0.0",
		"github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@v2.0.0",
	},
	hasError: true,
},
/*{
	name:           "Get Multiple different URIs",
	collectionName: "amadla",
	paths:          paths,
	entityURIs: []string{
		"github.com/AmadlaOrg/QAFixturesEntityPseudoVersion",
		"github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@v1.0.0",
	},
	hasError: false,
},*/
/*	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entityBuild := NewGetService()
			err = entityBuild.Get(test.collectionName, &test.paths, test.entityURIs)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Perform other assertions and checks as needed

			// Clean up
			err = os.RemoveAll(test.paths.Storage)
			if err != nil {
				t.Fatalf("Cleanup failed: %v", err)
			}
		})
	}
}*/
