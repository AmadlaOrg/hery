package entity

/*func createTestDirectoryStructure(t *testing.T, root string) {
	// Create test directories and files
	err := os.MkdirAll(filepath.Join(root, "origin1", "entity1@v1.0.0"), 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(root, "origin1", "entity2@v2.1.0"), 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(root, "origin2", "entity3@v0.9.1"), 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(filepath.Join(root, "origin2", "entity4@v1.1.0"), 0755)
	assert.NoError(t, err)

	// Create some files (to be ignored by the crawler)
	_, err = os.Create(filepath.Join(root, "origin1", "file1.txt"))
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(root, "origin2", "file2.txt"))
	assert.NoError(t, err)
}

func TestCrawlDirectoriesParallel(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Set up the test directory structure
	createTestDirectoryStructure(t, tmpDir)

	// Define the expected entities
	expectedEntities := map[string]Entity{
		"entity1": {Origin: "origin1", Version: "v1.0.0"},
		"entity2": {Origin: "origin1", Version: "v2.1.0"},
		"entity3": {Origin: "origin2", Version: "v0.9.1"},
		"entity4": {Origin: "origin2", Version: "v1.1.0"},
	}

	// Run the function under test
	entities, err := CrawlDirectoriesParallel(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, expectedEntities, entities)
}*/
