package query_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/AmadlaOrg/hery/entity/query"
	"github.com/AmadlaOrg/hery/entity/resolve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPipeline_ComposeDirToQueryFrom proves the critical-path pipeline:
// the multi-doc YAML emitted by `hery compose --dir` (resolve.MarshalAll)
// is re-parseable by `hery query --from -` (query.LoadDocs) and queryable
// in memory (query.QueryDocs). This is the contract behind
//
//	hery compose --dir ./graph | hery query --from - --type system/cpu --jq ._body
func TestPipeline_ComposeDirToQueryFrom(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600))
	}
	write("cpu.hery", "_type: system/cpu@v1.0.0\n_meta:\n  name: cpu\n_body:\n  cores: 4\n")
	write("memory.hery", "_type: system/memory@v1.0.0\n_meta:\n  name: mem\n_body:\n  size: 8\n")
	write("vm.hery", "_type: infrastructure/vm@v1.0.0\n_body:\n  hostname: web01\n")

	// Stage 1: compose --dir
	res, err := resolve.New().Resolve(dir)
	require.NoError(t, err)
	out, err := resolve.MarshalAll(res.Layers)
	require.NoError(t, err)

	// Stage 2: query --from - (parse the piped stream)
	docs, err := query.LoadDocs(bytes.NewReader(out))
	require.NoError(t, err)
	require.Len(t, docs, 3, "all three composed entities must survive the round-trip")

	// Stage 3: select + transform, just like the plugin would
	got, err := query.QueryDocs(docs, query.SelectionOpts{Type: "system/cpu", JQ: "._body"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.EqualValues(t, 4, got[0]["cores"])

	// Glob across the System sub-tree returns both system entities.
	sys, err := query.QueryDocs(docs, query.SelectionOpts{Type: "system/*"})
	require.NoError(t, err)
	assert.Len(t, sys, 2)
}
