package file

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestIsValidMagic(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		inputMagic     []byte
		internalOsStat func() (os.FileInfo, error)
		internalOsOpen func() (io.ReadCloser, error)
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
			expect:    false,
			hasError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
