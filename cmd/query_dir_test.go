package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AmadlaOrg/hery/entity/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryFromDir(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600))
	}
	write("cpu.hery", "_type: system/cpu@v1.0.0\n_meta:\n  name: cpu\n_body:\n  cores: 4\n")
	write("memory.hery", "_type: system/memory@v1.0.0\n_meta:\n  name: mem\n_body:\n  size: 8\n")
	write("vm.hery", "_type: infrastructure/vm@v1.0.0\n_body:\n  hostname: web01\n")

	// No selection returns every resolved entity in the directory.
	all, err := queryFromDir(dir, query.SelectionOpts{})
	require.NoError(t, err)
	assert.Len(t, all, 3)

	// Glob selection + jq behaves the same as the piped equivalent.
	got, err := queryFromDir(dir, query.SelectionOpts{Type: "system/*", JQ: "._body"})
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestRunQuery_DirAndFromMutuallyExclusive(t *testing.T) {
	require.NoError(t, QueryCmd.Flags().Set("dir", "/tmp/x"))
	require.NoError(t, QueryCmd.Flags().Set("from", "-"))
	t.Cleanup(func() {
		_ = QueryCmd.Flags().Set("dir", "")
		_ = QueryCmd.Flags().Set("from", "")
	})

	_, _, err := runQuery(QueryCmd)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mutually exclusive")
}
