package patch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeEnvPatch(t *testing.T, dir, env, id string, p Patch) {
	t.Helper()
	envDir := filepath.Join(dir, env)
	if err := os.MkdirAll(envDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(envDir, id+".json"), b, 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestPromote_CopiesNewPatches(t *testing.T) {
	tmp := t.TempDir()
	patchDir := filepath.Join(tmp, "patches")
	configDir := filepath.Join(tmp, "config")

	p := Patch{
		ID:      "001-add-feature",
		Version: "1.0",
		Ops:     []Op{{Op: "add", Path: "feature/enabled", Value: true}},
	}
	writeEnvPatch(t, patchDir, "staging", p.ID, p)

	result, err := Promote(patchDir, "staging", "production", configDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Copied) != 1 || result.Copied[0] != p.ID {
		t.Errorf("expected copied=[%s], got %v", p.ID, result.Copied)
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected no skipped, got %v", result.Skipped)
	}

	dest := filepath.Join(patchDir, "production", p.ID+".json")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("expected promoted patch file to exist at %s", dest)
	}
}

func TestPromote_SkipsAlreadyApplied(t *testing.T) {
	tmp := t.TempDir()
	patchDir := filepath.Join(tmp, "patches")
	configDir := filepath.Join(tmp, "config")

	p := Patch{
		ID:      "001-add-feature",
		Version: "1.0",
		Ops:     []Op{{Op: "add", Path: "feature/enabled", Value: true}},
	}
	writeEnvPatch(t, patchDir, "staging", p.ID, p)

	// Pre-record the patch as applied in the target history.
	prodConfigDir := filepath.Join(configDir, "production")
	if err := os.MkdirAll(prodConfigDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	h := &History{Applied: []string{p.ID}}
	if err := h.Save(HistoryPath(prodConfigDir)); err != nil {
		t.Fatalf("save history: %v", err)
	}

	result, err := Promote(patchDir, "staging", "production", configDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Copied) != 0 {
		t.Errorf("expected no copied patches, got %v", result.Copied)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != p.ID {
		t.Errorf("expected skipped=[%s], got %v", p.ID, result.Skipped)
	}
}

func TestPromote_MissingSourceDir(t *testing.T) {
	tmp := t.TempDir()
	_, err := Promote(filepath.Join(tmp, "patches"), "nonexistent", "production", tmp)
	if err == nil {
		t.Error("expected error for missing source directory")
	}
}

func TestFormatPromoteResult(t *testing.T) {
	r := &PromoteResult{
		SourceEnv: "staging",
		TargetEnv: "production",
		Copied:    []string{"001-init"},
		Skipped:   []string{"000-base"},
	}
	out := FormatPromoteResult(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
	for _, want := range []string{"staging", "production", "001-init", "000-base"} {
		if !contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
