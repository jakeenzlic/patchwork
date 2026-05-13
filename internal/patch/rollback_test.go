package patch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writePatchFile(t *testing.T, dir string, name string, p Patch) {
	t.Helper()
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatalf("marshal patch: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), b, 0644); err != nil {
		t.Fatalf("write patch file: %v", err)
	}
}

func TestRollback_SingleStep(t *testing.T) {
	dir := t.TempDir()

	// Initial config.
	cfgPath := filepath.Join(dir, "config.json")
	initial := map[string]any{"host": "localhost", "port": float64(8080)}
	b, _ := json.Marshal(initial)
	os.WriteFile(cfgPath, b, 0644)

	// Patch: add a new key.
	p := Patch{
		Version:     "1",
		Description: "add debug flag",
		Ops: []Op{
			{Op: "add", Path: "debug", Value: true},
		},
	}
	writePatchFile(t, dir, "001_add_debug.json", p)

	// Record patch as applied.
	hisFile := filepath.Join(dir, "history.json")
	h, _ := LoadHistory(hisFile)
	h.Applied = append(h.Applied, "001_add_debug.json")
	h.Save(hisFile)

	// Apply the patch manually so the config reflects it.
	current, _ := parseConfig(b, cfgPath)
	result, _ := Apply(current, p)
	out, _ := Export(result, "json")
	os.WriteFile(cfgPath, out, 0644)

	// Now rollback.
	err := Rollback(RollbackOptions{
		PatchDir:    dir,
		HistoryFile: hisFile,
		Target:      cfgPath,
		Steps:       1,
	})
	if err != nil {
		t.Fatalf("Rollback: %v", err)
	}

	// Verify debug key is gone.
	raw, _ := os.ReadFile(cfgPath)
	var got map[string]any
	json.Unmarshal(raw, &got)
	if _, ok := got["debug"]; ok {
		t.Error("expected debug key to be removed after rollback")
	}

	// Verify history trimmed.
	h2, _ := LoadHistory(hisFile)
	if len(h2.Applied) != 0 {
		t.Errorf("expected empty history, got %v", h2.Applied)
	}
}

func TestRollback_NoHistory(t *testing.T) {
	dir := t.TempDir()
	err := Rollback(RollbackOptions{
		PatchDir:    dir,
		HistoryFile: filepath.Join(dir, "history.json"),
		Target:      filepath.Join(dir, "config.json"),
	})
	if err == nil {
		t.Error("expected error when no history exists")
	}
}
