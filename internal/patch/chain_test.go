package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeChainPatch(id, op, path string, value any) Patch {
	return Patch{
		ID:      id,
		Version: "1",
		Ops: []Op{
			{Op: op, Path: path, Value: value},
		},
	}
}

func baseChainConfig() map[string]any {
	return map[string]any{"app": map[string]any{"port": 8080}}
}

func TestChain_AppliesAllPatches(t *testing.T) {
	dir := t.TempDir()
	patches := []Patch{
		makeChainPatch("p1", "add", "app.debug", true),
		makeChainPatch("p2", "replace", "app.port", 9090),
	}
	res, err := Chain(patches, baseChainConfig(), dir)
	require.NoError(t, err)
	assert.Equal(t, []string{"p1", "p2"}, res.Applied)
	assert.Empty(t, res.Skipped)
	assert.Equal(t, 9090, res.Config["app"].(map[string]any)["port"])
}

func TestChain_SkipsAppliedPatches(t *testing.T) {
	dir := t.TempDir()
	hist, _ := LoadHistory(dir)
	_ = hist.Record("p1")
	_ = hist.Save(dir)

	patches := []Patch{
		makeChainPatch("p1", "add", "app.debug", true),
		makeChainPatch("p2", "add", "app.name", "svc"),
	}
	res, err := Chain(patches, baseChainConfig(), dir)
	require.NoError(t, err)
	assert.Equal(t, []string{"p2"}, res.Applied)
	assert.Equal(t, []string{"p1"}, res.Skipped)
}

func TestChain_StopsOnInvalidPatch(t *testing.T) {
	dir := t.TempDir()
	bad := Patch{ID: "bad", Version: "", Ops: []Op{{Op: "add", Path: "x", Value: 1}}}
	res, err := Chain([]Patch{bad}, baseChainConfig(), dir)
	require.Error(t, err)
	assert.Equal(t, "bad", res.Failed)
	assert.Empty(t, res.Applied)
}

func TestChain_StopsOnApplyError(t *testing.T) {
	dir := t.TempDir()
	// replace on non-existent path should error
	bad := makeChainPatch("p1", "replace", "does.not.exist", 42)
	res, err := Chain([]Patch{bad}, baseChainConfig(), dir)
	require.Error(t, err)
	assert.Equal(t, "p1", res.Failed)
}

func TestFormatChainResult_Summary(t *testing.T) {
	r := ChainResult{
		Applied: []string{"p1"},
		Skipped: []string{"p2"},
		Failed:  "p3",
	}
	out := FormatChainResult(r)
	assert.Contains(t, out, "applied  p1")
	assert.Contains(t, out, "skipped  p2")
	assert.Contains(t, out, "failed   p3")
}

func TestFormatChainResult_Empty(t *testing.T) {
	out := FormatChainResult(ChainResult{})
	assert.Contains(t, out, "nothing to do")
}
