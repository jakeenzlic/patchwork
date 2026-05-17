package patch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeExtractConfig(t *testing.T, dir string, data map[string]any) string {
	t.Helper()
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	p := filepath.Join(dir, "config.json")
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestExtract_RoundtripViaExport(t *testing.T) {
	dir := t.TempDir()
	cfg := map[string]any{
		"service": map[string]any{
			"host": "0.0.0.0",
			"port": 8080,
		},
		"logging": map[string]any{
			"level": "info",
		},
	}

	r, err := Extract(cfg, "service")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	out := filepath.Join(dir, "service.json")
	if err := Export(r.Data, out, "json"); err != nil {
		t.Fatalf("export: %v", err)
	}

	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var loaded map[string]any
	if err := json.Unmarshal(b, &loaded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if loaded["host"] != "0.0.0.0" {
		t.Errorf("expected host=0.0.0.0, got %v", loaded["host"])
	}
}

func TestExtract_DeepNesting(t *testing.T) {
	cfg := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"value": 42,
				},
			},
		},
	}

	r, err := Extract(cfg, "a/b/c")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if r.Data["value"] != 42 {
		t.Errorf("expected value=42, got %v", r.Data["value"])
	}
}

func TestExtract_MutationIsolation(t *testing.T) {
	cfg := map[string]any{
		"ns": map[string]any{
			"key": "original",
		},
	}
	r, err := Extract(cfg, "ns")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	r.Data["key"] = "mutated"
	if cfg["ns"].(map[string]any)["key"] != "original" {
		t.Error("original config was mutated")
	}
}
