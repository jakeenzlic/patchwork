package patch

import (
	"testing"
)

var baseExtractConfig = map[string]any{
	"database": map[string]any{
		"host": "localhost",
		"port": 5432,
		"credentials": map[string]any{
			"user": "admin",
			"pass": "secret",
		},
	},
	"app": map[string]any{
		"name": "patchwork",
		"debug": true,
	},
}

func TestExtract_TopLevelPrefix(t *testing.T) {
	r, err := Extract(baseExtractConfig, "database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Prefix != "database" {
		t.Errorf("expected prefix %q, got %q", "database", r.Prefix)
	}
	if _, ok := r.Data["host"]; !ok {
		t.Error("expected key 'host' in extracted data")
	}
	if _, ok := r.Data["credentials"]; !ok {
		t.Error("expected key 'credentials' in extracted data")
	}
}

func TestExtract_NestedPrefix(t *testing.T) {
	r, err := Extract(baseExtractConfig, "database/credentials")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Data["user"] != "admin" {
		t.Errorf("expected user=admin, got %v", r.Data["user"])
	}
}

func TestExtract_LeadingSlashStripped(t *testing.T) {
	r, err := Extract(baseExtractConfig, "/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Data["name"] != "patchwork" {
		t.Errorf("expected name=patchwork, got %v", r.Data["name"])
	}
}

func TestExtract_MissingPath(t *testing.T) {
	_, err := Extract(baseExtractConfig, "missing/path")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestExtract_NonMapValue(t *testing.T) {
	_, err := Extract(baseExtractConfig, "database/host")
	if err == nil {
		t.Fatal("expected error when target is not a map")
	}
}

func TestExtract_EmptyPrefix(t *testing.T) {
	_, err := Extract(baseExtractConfig, "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestExtract_IsolatesData(t *testing.T) {
	r, err := Extract(baseExtractConfig, "database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Mutating extracted data must not affect original.
	r.Data["injected"] = "evil"
	db := baseExtractConfig["database"].(map[string]any)
	if _, ok := db["injected"]; ok {
		t.Error("extract mutated the original config")
	}
}

func TestFormatExtractResult_Output(t *testing.T) {
	r := ExtractResult{Prefix: "app", Data: map[string]any{"name": "pw", "debug": true}}
	out := FormatExtractResult(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
