package patch_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"patchwork/internal/patch"
)

func writeChainPatch(t *testing.T, dir, id, op, path string, value any) {
	t.Helper()
	p := patch.Patch{
		ID:      id,
		Version: "1",
		Ops:     []patch.Op{{Op: op, Path: path, Value: value}},
	}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, id+".json"), b, 0o644))
}

func TestChain_FullRoundtrip(t *testing.T) {
	patchDir := t.TempDir()
	histDir := t.TempDir()

	writeChainPatch(t, patchDir, "01-add-debug", "add", "app.debug", false)
	writeChainPatch(t, patchDir, "02-set-port", "replace", "app.port", 9000)

	cfg := map[string]any{"app": map[string]any{"port": 8080}}

	patches, err := patch.LoadDir(patchDir)
	require.NoError(t, err)

	res, err := patch.Chain(patches, cfg, histDir)
	require.NoError(t, err)
	assert.Len(t, res.Applied, 2)
	app := res.Config["app"].(map[string]any)
	assert.Equal(t, 9000, app["port"])
	assert.Equal(t, false, app["debug"])
}

func TestChain_IdempotentOnRerun(t *testing.T) {
	patchDir := t.TempDir()
	histDir := t.TempDir()

	writeChainPatch(t, patchDir, "p1", "add", "x", 1)

	cfg := map[string]any{}
	patches, err := patch.LoadDir(patchDir)
	require.NoError(t, err)

	// First run
	res1, err := patch.Chain(patches, cfg, histDir)
	require.NoError(t, err)
	assert.Equal(t, []string{"p1"}, res1.Applied)

	// Record applied so second run sees it
	hist, err := patch.LoadHistory(histDir)
	require.NoError(t, err)
	require.NoError(t, hist.Record("p1"))
	require.NoError(t, hist.Save(histDir))

	// Second run — should skip
	res2, err := patch.Chain(patches, res1.Config, histDir)
	require.NoError(t, err)
	assert.Empty(t, res2.Applied)
	assert.Equal(t, []string{"p1"}, res2.Skipped)
}
