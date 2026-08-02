package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AmadlaOrg/hery/cache"
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

// A directory with no .hery files is an empty selection (exit 2), not an
// error — and no .hery.cache is littered into it.
func TestRunQuery_DefaultEmptyDir(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, code, err := runQuery(QueryCmd)
	require.NoError(t, err)
	assert.Equal(t, 2, code)
	assert.NoFileExists(t, filepath.Join(dir, cache.ProjectCacheFile))
}

// The default source materializes the project cache and serves from it while
// the sources are unchanged. Cache use is proven by mutating a file while
// restoring its mtime and size: the cached (old) content must come back.
func TestRunQuery_DefaultUsesProjectCache(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cpu.hery")
	require.NoError(t, os.WriteFile(path,
		[]byte("_type: system/cpu@v1.0.0\n_meta:\n  name: aaaa\n"), 0o600))
	t.Chdir(dir)

	out, code, err := runQuery(QueryCmd)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	assert.Contains(t, out, "aaaa")
	require.FileExists(t, filepath.Join(dir, cache.ProjectCacheFile))

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path,
		[]byte("_type: system/cpu@v1.0.0\n_meta:\n  name: bbbb\n"), 0o600))
	require.NoError(t, os.Chtimes(path, info.ModTime(), info.ModTime()))

	out, code, err = runQuery(QueryCmd)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	assert.Contains(t, out, "aaaa", "unchanged sources must be served from .hery.cache")

	// A visible change (new mtime) rebuilds the cache.
	future := time.Now().Add(2 * time.Second)
	require.NoError(t, os.Chtimes(path, future, future))

	out, code, err = runQuery(QueryCmd)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	assert.Contains(t, out, "bbbb")
}
