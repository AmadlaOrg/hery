package file

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsFile(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		internalOsStat func(string) (os.FileInfo, error)
		expected       bool
		expectErr      error
		hasError       bool
	}{
		{
			name:      "file exists",
			inputPath: "./testdata/test.txt",
			internalOsStat: func(string) (os.FileInfo, error) {
				mockFileInfo := NewMockFileInfo(t)
				mockFileInfo.EXPECT().IsDir().Return(false)
				return mockFileInfo, nil
			},
			expected:  true,
			expectErr: nil,
			hasError:  false,
		},
		//
		// Error
		//
		{
			name:      "file exists",
			inputPath: "./testdata/test.txt",
			internalOsStat: func(string) (os.FileInfo, error) {
				mockFileInfo := NewMockFileInfo(t)
				return mockFileInfo, errors.New("stat ./testdata/test.txt: no such file or directory")
			},
			expected:  false,
			expectErr: errors.New("stat ./testdata/test.txt: no such file or directory"),
			hasError:  true,
		},
		{
			name:      "file exists",
			inputPath: "./testdata/test.txt",
			internalOsStat: func(string) (os.FileInfo, error) {
				mockFileInfo := NewMockFileInfo(t)
				mockFileInfo.EXPECT().IsDir().Return(true)
				return mockFileInfo, nil
			},
			expected:  false,
			expectErr: errors.New("not a file"),
			hasError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalOsStat := osStat
			defer func() { osStat = originalOsStat }()
			osStat = tt.internalOsStat

			got, err := IsFile(tt.inputPath)
			if tt.hasError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestIsValidMagic(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		inputMagic     []byte
		internalOsOpen func(name string) (IFile, error)
		expect         bool
		expectedError  error
		hasError       bool
	}{
		/*{
			name:      "normal",
			inputPath: `testdata/valid.db`,
			expect:    true,
			hasError:  false,
		},*/
		//
		// Error
		//
		{
			name:      "Error: normal",
			inputPath: `testdata/valid.db`,
			internalOsOpen: func(name string) (IFile, error) {
				return os.Open("testdata/valid.db")
			},
			expect:        false,
			expectedError: errors.New("not a valid magic file"),
			hasError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalOsStat := osStat
			defer func() { osStat = originalOsStat }()
			osStat = func(string) (os.FileInfo, error) {
				mockFileInfo := NewMockFileInfo(t)
				mockFileInfo.EXPECT().IsDir().Return(false)
				return mockFileInfo, nil
			}

			originalOsOpen := osOpen
			defer func() { osOpen = originalOsOpen }()
			osOpen = tt.internalOsOpen

			got, err := IsValidMagic(tt.inputPath, tt.inputMagic)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expect, got)
		})
	}
}
