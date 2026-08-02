package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The default data source (no --from/--dir) queries the current directory's
// .hery files, exactly like --dir '.'.
func TestRunQuery_DefaultQueriesCwd(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "cpu.hery"),
		[]byte("_type: system/cpu@v1.0.0\n_meta:\n  name: cpu\n_body:\n  cores: 4\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "vm.hery"),
		[]byte("_type: infrastructure/vm@v1.0.0\n_body:\n  hostname: web01\n"), 0o600))
	t.Chdir(dir)

	out, code, err := runQuery(QueryCmd)
	require.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Contains(t, out, "system/cpu@v1.0.0")
	assert.Contains(t, out, "infrastructure/vm@v1.0.0")
}

// A directory with no .hery files is an empty selection (exit 2), not an error.
func TestRunQuery_DefaultEmptyDir(t *testing.T) {
	t.Chdir(t.TempDir())

	_, code, err := runQuery(QueryCmd)
	require.NoError(t, err)
	assert.Equal(t, 2, code)
}
