package patch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestExport_JSON(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "config.json")

	cfg := map[string]any{"version": "1.0", "debug": true}

	if err := Export(cfg, dest, FormatJSON); err != nil {
		t.Fatalf("Export JSON failed: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read exported file: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("exported file is not valid JSON: %v", err)
	}

	if result["version"] != "1.0" {
		t.Errorf("expected version=1.0, got %v", result["version"])
	}
}

func TestExport_YAML(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "config.yaml")

	cfg := map[string]any{"env": "production", "replicas": 3}

	if err := Export(cfg, dest, FormatYAML); err != nil {
		t.Fatalf("Export YAML failed: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read exported file: %v", err)
	}

	var result map[string]any
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Fatalf("exported file is not valid YAML: %v", err)
	}

	if result["env"] != "production" {
		t.Errorf("expected env=production, got %v", result["env"])
	}
}

func TestExport_InferFormat_JSON(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.json")

	cfg := map[string]any{"key": "value"}
	if err := Export(cfg, dest, ""); err != nil {
		t.Fatalf("Export with inferred JSON format failed: %v", err)
	}

	data, _ := os.ReadFile(dest)
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Errorf("inferred JSON output is not valid JSON: %v", err)
	}
}

func TestExport_InferFormat_YAML(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.yml")

	cfg := map[string]any{"key": "value"}
	if err := Export(cfg, dest, ""); err != nil {
		t.Fatalf("Export with inferred YAML format failed: %v", err)
	}

	data, _ := os.ReadFile(dest)
	var result map[string]any
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Errorf("inferred YAML output is not valid YAML: %v", err)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.toml")

	cfg := map[string]any{"key": "value"}
	if err := Export(cfg, dest, "toml"); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
