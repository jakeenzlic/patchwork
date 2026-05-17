package patch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRedact_ExportRoundtrip(t *testing.T) {
	cfg := map[string]any{
		"service": map[string]any{
			"api_key": "abc123",
			"url":     "https://example.com",
		},
	}

	rules := []RedactRule{{Path: "service/api_key"}}
	res := Redact(cfg, rules)

	dir := t.TempDir()
	out := filepath.Join(dir, "config.json")

	if err := Export(res.Config, out, "json"); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	var loaded map[string]any
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	svc := loaded["service"].(map[string]any)
	if svc["api_key"] != "***" {
		t.Errorf("expected *** in exported file, got %v", svc["api_key"])
	}
	if svc["url"] != "https://example.com" {
		t.Errorf("url should be unmasked")
	}
}

func TestRedact_DeepNesting(t *testing.T) {
	cfg := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep-secret",
			},
		},
	}

	rules := []RedactRule{{Path: "a/b/c"}}
	res := Redact(cfg, rules)

	aMap := res.Config["a"].(map[string]any)
	bMap := aMap["b"].(map[string]any)
	if bMap["c"] != "***" {
		t.Errorf("expected *** at a/b/c, got %v", bMap["c"])
	}
}
