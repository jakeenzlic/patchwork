package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotPath(t *testing.T) {
	got := SnapshotPath("/etc/app/config.yaml")
	want := "/etc/app/.patchwork/snapshots/config.yaml.snap.json"
	if got != want {
		t.Errorf("SnapshotPath = %q, want %q", got, want)
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	cfg := map[string]interface{}{
		"host": "localhost",
		"port": float64(8080),
	}

	if err := SaveSnapshot(configPath, "v1", cfg); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	snap, err := LoadSnapshot(configPath)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if snap.PatchID != "v1" {
		t.Errorf("PatchID = %q, want %q", snap.PatchID, "v1")
	}
	if snap.Config["host"] != "localhost" {
		t.Errorf("Config[host] = %v, want localhost", snap.Config["host"])
	}
	if snap.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestLoadSnapshot_NotExist(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "missing.json")

	snap, err := LoadSnapshot(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Errorf("expected nil snapshot for missing file, got %+v", snap)
	}
}

func TestSaveSnapshot_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "sub", "config.yaml")

	cfg := map[string]interface{}{"key": "value"}
	if err := SaveSnapshot(configPath, "v2", cfg); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	expectedDir := filepath.Join(dir, "sub", ".patchwork", "snapshots")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("expected directory %q to be created", expectedDir)
	}
}
