package entity

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEntity(t *testing.T) {
	err := BuildEntity("CPU", true)
	if err != nil {
		return
	}
}

func TestParseEntityArg(t *testing.T) {
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
			gotName, gotVersion, err := parseEntityArg(tt.entityArg)
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
		{
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
		},
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

			gotDir, err := findEntityDirParallel(tmpDir, tt.entityName, tt.version)
			if tt.expectedErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got %v", err)
				assert.Equal(t, tt.expectedDir, gotDir, "expected dir %v, got %v", tt.expectedDir, gotDir)
			}
		})
	}
}
