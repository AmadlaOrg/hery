package resolve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
}

func TestResolve_SingleFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "pkg.hery", `_type: example/package@v1.0.0
_body:
  name: nginx
`)
	res, err := New().Resolve(dir)
	require.NoError(t, err)
	require.Len(t, res.Layers, 1)
	require.Len(t, res.Layers[0].Docs, 1)
	assert.Equal(t, "example/package@v1.0.0", res.Layers[0].Docs[0].Type)
}

func TestResolve_LocalRequiresOrdering(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.hery", `_type: a@v1
_requires:
  - ./b.hery
_body:
  x: 1
`)
	writeFile(t, dir, "b.hery", `_type: b@v1
_requires:
  - ./c.hery
_body:
  x: 2
`)
	writeFile(t, dir, "c.hery", `_type: c@v1
_body:
  x: 3
`)

	res, err := New().Resolve(dir)
	require.NoError(t, err)
	require.Len(t, res.Layers, 1)
	docs := res.Layers[0].Docs
	require.Len(t, docs, 3)
	assert.Equal(t, "a@v1", docs[0].Type)
	assert.Equal(t, "b@v1", docs[1].Type)
	assert.Equal(t, "c@v1", docs[2].Type)
}

func TestResolve_ExtendsLocalMergesBody(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "parent.hery", `_type: parent@v1
_body:
  a: 1
  b: 2
`)
	writeFile(t, dir, "child.hery", `_type: child@v1
_extends: ./parent.hery
_body:
  b: 99
  c: 3
`)

	res, err := New().Resolve(dir)
	require.NoError(t, err)
	require.Len(t, res.Layers, 1)

	out, err := MarshalAll(res.Layers)
	require.NoError(t, err)
	// Intra-layer order is alphabetical by filename: child.hery before parent.hery.
	// _extends is a merge edge, not a dependency edge, so it doesn't drive ordering.
	want := `---
_type: child@v1
_extends: ./parent.hery
_body:
  a: 1
  b: 99
  c: 3
---
_type: parent@v1
_body:
  a: 1
  b: 2
`
	assert.Equal(t, want, string(out))
}

func TestResolve_ExternalRequiresBecomesLayer2(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "input")
	extDir := filepath.Join(root, "ext")
	writeFile(t, inputDir, "main.hery", `_type: main@v1
_requires:
  - ../ext/
_body:
  k: v
`)
	writeFile(t, extDir, "dep.hery", `_type: dep@v1
_body:
  n: 1
`)

	res, err := New().Resolve(inputDir)
	require.NoError(t, err)
	require.Len(t, res.Layers, 2)
	assert.Equal(t, "main@v1", res.Layers[0].Docs[0].Type)
	assert.Equal(t, "dep@v1", res.Layers[1].Docs[0].Type)
}

func TestResolve_ExtendsExternalAddsLayerAndPreservesBoth(t *testing.T) {
	root := t.TempDir()
	inputDir := filepath.Join(root, "input")
	extDir := filepath.Join(root, "ext")
	writeFile(t, extDir, "base.hery", `_type: base@v1
_body:
  base: true
  shared: parent
`)
	writeFile(t, inputDir, "child.hery", `_type: child@v1
_extends: ../ext/base.hery
_body:
  shared: child
`)

	res, err := New().Resolve(inputDir)
	require.NoError(t, err)
	require.Len(t, res.Layers, 2)

	childBody := extractBody(res.Layers[0].Docs[0].Raw)
	assert.NotNil(t, childBody)

	// External parent appears unmerged at layer 2
	assert.Equal(t, "base@v1", res.Layers[1].Docs[0].Type)
}

func TestResolve_ErrorOnNonLocalReference(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "x.hery", `_type: x@v1
_requires:
  - github.com/foo/bar@v1.0.0
_body: {}
`)
	_, err := New().Resolve(dir)
	assert.Error(t, err)
}
