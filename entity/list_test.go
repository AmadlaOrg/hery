package entity

import (
	"testing"
)

func TestCrawlDirectoriesParallel(t *testing.T) {
	// Create a temporary directory structure for testing
	/*tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "github.com", "AmadlaOrg")
	err := os.MkdirAll(filepath.Join(testDir, "Entity", "CPU@v1.0.0"), 0755)
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Join(testDir, "Entity", "Memory@v1.0.0"), 0755)
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Join(testDir, "Entity", "Net@v1.0.0"), 0755)
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Join(testDir, "Entity", "Storage@v1.0.0"), 0755)
	if err != nil {
		return
	}

	entities, err := CrawlDirectoriesParallel(tmpDir)
	assert.NoError(t, err)
	expectedEntities := map[string]string{
		"CPU":     "1.0.0",
		"Memory":  "1.0.0",
		"Net":     "1.0.0",
		"Storage": "1.0.0",
	}

	assert.Equal(t, expectedEntities, entities)*/
}
