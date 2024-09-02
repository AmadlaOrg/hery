package compose

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEntity(t *testing.T) {
	/*err := ComposeEntity("CPU", true)
	if err != nil {
		return
	}*/
}

func TestParseEntityArg(t *testing.T) {
	composeService := NewComposeService()

	tests := []struct {
		name        string
		entityArg   string
		wantName    string
		wantVersion string
		wantErr     bool
	}{
		{
			name:        "Valid entity without version",
			entityArg:   "CPU",
			wantName:    "CPU",
			wantVersion: "",
			wantErr:     false,
		},
		{
			name:        "Valid entity with version",
			entityArg:   "CPU@v1.0.0",
			wantName:    "CPU",
			wantVersion: "@v1.0.0",
			wantErr:     false,
		},
		{
			name:        "Invalid entity argument without @v",
			entityArg:   "CPU@1.0.0",
			wantName:    "",
			wantVersion: "",
			wantErr:     true,
		},
		{
			name:        "Empty entity argument",
			entityArg:   "",
			wantName:    "",
			wantVersion: "",
			wantErr:     true,
		},
		{
			name:        "Only version without entity",
			entityArg:   "@v1.0.0",
			wantName:    "",
			wantVersion: "",
			wantErr:     true,
		},
		{
			name:        "Invalid characters in entity name",
			entityArg:   "CPU-123",
			wantName:    "",
			wantVersion: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotVersion, err := composeService.parseEntityArg(tt.entityArg)
			if tt.wantErr {
				assert.Error(t, err, fmt.Sprintf("expected error but got none for %v", tt.entityArg))
			} else {
				assert.NoError(t, err, fmt.Sprintf("expected no error but got %v", err))
				assert.Equal(t, tt.wantName, gotName, fmt.Sprintf("expected name %v, got %v", tt.wantName, gotName))
				assert.Equal(t, tt.wantVersion, gotVersion, fmt.Sprintf("expected version %v, got %v", tt.wantVersion, gotVersion))
			}
		})
	}
}

func TestFindEntityDirParallel(t *testing.T) {
	composeService := NewComposeService()

	tests := []struct {
		name        string
		rootSetup   func(root string) error
		entityName  string
		version     string
		expectedDir string
		expectedErr bool
	}{
		{
			name: "Find specific version",
			rootSetup: func(root string) error {
				testDir := filepath.Join(root, "github.com", "AmadlaOrg")
				err := os.MkdirAll(filepath.Join(testDir, "CPU@v1.0.0"), 0755)
				if err != nil {
					return err
				}
				err = os.MkdirAll(filepath.Join(testDir, "CPU@v2.0.0"), 0755)
				if err != nil {
					return err
				}
				return nil
			},
			entityName:  "CPU",
			version:     "@v1.0.0",
			expectedDir: "github.com/AmadlaOrg/CPU@v1.0.0",
			expectedErr: false,
		},
		/*{
			name: "Find latest version",
			rootSetup: func(root string) error {
				testDir := filepath.Join(root, "github.com", "AmadlaOrg")
				err := os.MkdirAll(filepath.Join(testDir, "CPU@v1.0.0"), 0755)
				if err != nil {
					return err
				}
				err = os.MkdirAll(filepath.Join(testDir, "CPU@v2.0.0"), 0755)
				if err != nil {
					return err
				}
				return nil
			},
			entityName:  "CPU",
			version:     "",
			expectedDir: "github.com/AmadlaOrg/CPU@v2.0.0",
			expectedErr: false,
		},*/
		{
			name: "Entity not found",
			rootSetup: func(root string) error {
				testDir := filepath.Join(root, "github.com", "AmadlaOrg")
				err := os.MkdirAll(filepath.Join(testDir, "Memory@v1.0.0"), 0755)
				if err != nil {
					return err
				}
				return nil
			},
			entityName:  "CPU",
			version:     "@v1.0.0",
			expectedDir: "",
			expectedErr: true,
		},
		{
			name: "Invalid entity name",
			rootSetup: func(root string) error {
				testDir := filepath.Join(root, "github.com", "AmadlaOrg")
				err := os.MkdirAll(filepath.Join(testDir, "CPU@v1.0.0"), 0755)
				if err != nil {
					return err
				}
				return nil
			},
			entityName:  "INVALID@v1.0.0",
			version:     "",
			expectedDir: "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			err := tt.rootSetup(tmpDir)
			assert.NoError(t, err)

			// Convert expectedDir to absolute path
			if tt.expectedDir != "" {
				tt.expectedDir = filepath.Join(tmpDir, tt.expectedDir)
			}

			gotDir, err := composeService.findEntityDirParallel(tmpDir, tt.entityName, tt.version)
			if tt.expectedErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got %v", err)
				assert.Equal(t, tt.expectedDir, gotDir, "expected dir %v, got %v", tt.expectedDir, gotDir)
			}
		})
	}
}

//func TestMergeYamlFiles(t *testing.T) {
// Create a temporary directory for testing
/*dir, err := os.MkdirTemp("", "test-yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create sample YAML files
	file1 := filepath.Join(dir, "file1.yml")
	file2 := filepath.Join(dir, "file2.yaml")

	content1 := `
key1:
  subkey1: value1
key2: value2
`
	content2 := `
key1:
  subkey2: value3
key3: value4
`

	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatal(err)
	}

	// Call the function to merge YAML files
	mergedBytes, err := mergeYamlFiles(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the merged result
	var merged map[string]interface{}
	err = yaml.Unmarshal(mergedBytes, &merged)
	if err != nil {
		t.Fatal(err)
	}

	// Expected result
	expected := map[string]interface{}{
		"key1": map[interface{}]interface{}{
			"subkey1": "value1",
			"subkey2": "value3",
		},
		"key2": "value2",
		"key3": "value4",
	}

	// Use testify's assert package to compare the results
	assert.Equal(t, expected, merged)*/
//}

// Sample merge function from utilObjectPkg (to be replaced with actual implementation)
/*func MergeMultilevel(dst, src map[string]interface{}, overwrite bool) map[string]interface{} {
	for k, v := range src {
		if vMap, ok := v.(map[string]interface{}); ok {
			if dstV, ok := dst[k].(map[string]interface{}); ok {
				dst[k] = MergeMultilevel(dstV, vMap, overwrite)
				continue
			}
		}
		if overwrite {
			dst[k] = v
		} else {
			if _, ok := dst[k]; !ok {
				dst[k] = v
			}
		}
	}
	return dst
}*/

func TestMergeYamlFiles(t *testing.T) {
	// Create temporary directory for test YAML files
	tempDir, err := os.MkdirTemp("", "test_mergeYamlFiles")
	assert.NoError(t, err, "Failed to create temp dir")
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
	}(tempDir)

	// Create test YAML files
	yamlFiles := map[string]string{
		"file1.yaml": `
key1: value1
key2:
  subkey1: subvalue1
  subkey2: subvalue2
`,
		"file2.yml": `
key1: newValue1
key2:
  subkey2: newSubvalue2
  subkey3: subvalue3
key3: value3
`,
	}

	for filename, content := range yamlFiles {
		err := os.WriteFile(filepath.Join(tempDir, filename), []byte(content), 0644)
		assert.NoError(t, err, "Failed to write file %s", filename)
	}

	composeService := NewComposeService()
	// Call mergeYamlFiles function
	mergedContent, err := composeService.mergeYamlFiles(tempDir)
	assert.NoError(t, err, "mergeYamlFiles returned an error")

	// Unmarshal merged content
	var mergedData map[string]interface{}
	err = yaml.Unmarshal(mergedContent, &mergedData)
	assert.NoError(t, err, "Failed to unmarshal merged content")

	// Expected merged content
	expectedData := map[string]interface{}{
		"key1": "newValue1",
		"key2": map[string]interface{}{
			"subkey1": "subvalue1",
			"subkey2": "newSubvalue2",
			"subkey3": "subvalue3",
		},
		"key3": "value3",
	}

	// Compare merged data with expected data
	assert.True(t, compareYamlMaps(convertYamlMap(mergedData), expectedData), "Merged content does not match expected content.\nGot: %#v\nExpected: %#v", mergedData, expectedData)
}

func compareYamlMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, vA := range a {
		vB, ok := b[k]
		if !ok {
			return false
		}
		switch vA := vA.(type) {
		case map[string]interface{}:
			vBMap, ok := vB.(map[string]interface{})
			if !ok {
				return false
			}
			if !compareYamlMaps(vA, vBMap) {
				return false
			}
		default:
			if vA != vB {
				return false
			}
		}
	}
	return true
}

func convertYamlMap(m map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for k, v := range m {
		switch v := v.(type) {
		case map[interface{}]interface{}:
			converted[k] = convertYamlMap(convertYamlMapKeys(v))
		case map[string]interface{}:
			converted[k] = convertYamlMap(v)
		default:
			converted[k] = v
		}
	}
	return converted
}

func convertYamlMapKeys(m map[interface{}]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for k, v := range m {
		converted[k.(string)] = v
	}
	return converted
}
