package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/AmadlaOrg/hery/entity/resolve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeHery(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func resolveDir(t *testing.T, dir string) *resolve.Result {
	t.Helper()
	res, err := resolve.New().Resolve(dir)
	require.NoError(t, err)
	return res
}

func TestProjectCache_MissWithoutCacheFile(t *testing.T) {
	pc := NewProject(t.TempDir())
	_, ok := pc.LoadFresh()
	assert.False(t, ok)
}

func TestProjectCache_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	writeHery(t, dir, "cpu.hery",
		"_type: system/cpu@v1.0.0\n_meta:\n  name: cpu\n_body:\n  cores: 4\n")
	writeHery(t, dir, "vm.hery",
		"_type: infrastructure/vm@v1.0.0\n_body:\n  hostname: web01\n")

	pc := NewProject(dir)
	require.NoError(t, pc.Rebuild(resolveDir(t, dir)))
	require.FileExists(t, filepath.Join(dir, ProjectCacheFile))

	docs, ok := pc.LoadFresh()
	require.True(t, ok)
	require.Len(t, docs, 2)

	types := []string{}
	for _, d := range docs {
		typ, _ := d["_type"].(string)
		types = append(types, typ)
	}
	assert.Contains(t, types, "system/cpu@v1.0.0")
	assert.Contains(t, types, "infrastructure/vm@v1.0.0")

	for _, d := range docs {
		if d["_type"] == "system/cpu@v1.0.0" {
			body, bok := d["_body"].(map[string]any)
			require.True(t, bok)
			assert.EqualValues(t, 4, body["cores"])
		}
	}
}

func TestProjectCache_StaleOnFileChange(t *testing.T) {
	dir := t.TempDir()
	path := writeHery(t, dir, "cpu.hery", "_type: system/cpu@v1.0.0\n_body:\n  cores: 4\n")

	pc := NewProject(dir)
	require.NoError(t, pc.Rebuild(resolveDir(t, dir)))

	future := time.Now().Add(2 * time.Second)
	require.NoError(t, os.Chtimes(path, future, future))

	_, ok := pc.LoadFresh()
	assert.False(t, ok)
}

func TestProjectCache_StaleOnNewFile(t *testing.T) {
	dir := t.TempDir()
	writeHery(t, dir, "cpu.hery", "_type: system/cpu@v1.0.0\n_body:\n  cores: 4\n")

	pc := NewProject(dir)
	require.NoError(t, pc.Rebuild(resolveDir(t, dir)))

	writeHery(t, dir, "memory.hery", "_type: system/memory@v1.0.0\n_body:\n  size: 8\n")

	_, ok := pc.LoadFresh()
	assert.False(t, ok)
}

func TestProjectCache_StaleOnRemovedFile(t *testing.T) {
	dir := t.TempDir()
	writeHery(t, dir, "cpu.hery", "_type: system/cpu@v1.0.0\n_body:\n  cores: 4\n")
	victim := writeHery(t, dir, "memory.hery", "_type: system/memory@v1.0.0\n_body:\n  size: 8\n")

	pc := NewProject(dir)
	require.NoError(t, pc.Rebuild(resolveDir(t, dir)))

	require.NoError(t, os.Remove(victim))

	_, ok := pc.LoadFresh()
	assert.False(t, ok)
}

// An unchanged source set is served from the cache: mutate a file's content
// but restore its mtime and size so the manifest still matches — the cached
// (old) content must come back, proving no re-resolve happened.
func TestProjectCache_FreshServesCachedContent(t *testing.T) {
	dir := t.TempDir()
	path := writeHery(t, dir, "cpu.hery", "_type: system/cpu@v1.0.0\n_meta:\n  name: aaaa\n")

	pc := NewProject(dir)
	require.NoError(t, pc.Rebuild(resolveDir(t, dir)))

	info, err := os.Stat(path)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(path, []byte("_type: system/cpu@v1.0.0\n_meta:\n  name: bbbb\n"), 0o600))
	require.NoError(t, os.Chtimes(path, info.ModTime(), info.ModTime()))

	docs, ok := pc.LoadFresh()
	require.True(t, ok)
	require.Len(t, docs, 1)
	meta, mok := docs[0]["_meta"].(map[string]any)
	require.True(t, mok)
	assert.Equal(t, "aaaa", meta["name"])
}

func TestProjectCache_RebuildReplacesPreviousRows(t *testing.T) {
	dir := t.TempDir()
	path := writeHery(t, dir, "cpu.hery", "_type: system/cpu@v1.0.0\n_body:\n  cores: 4\n")

	pc := NewProject(dir)
	require.NoError(t, pc.Rebuild(resolveDir(t, dir)))

	require.NoError(t, os.WriteFile(path, []byte("_type: system/cpu@v1.0.0\n_body:\n  cores: 8\n"), 0o600))
	require.NoError(t, pc.Rebuild(resolveDir(t, dir)))

	docs, ok := pc.LoadFresh()
	require.True(t, ok)
	require.Len(t, docs, 1)
	body, bok := docs[0]["_body"].(map[string]any)
	require.True(t, bok)
	assert.EqualValues(t, 8, body["cores"])
}
