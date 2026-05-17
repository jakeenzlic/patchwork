package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckpointPath(t *testing.T) {
	p := CheckpointPath("/srv/cfg")
	want := filepath.Join("/srv/cfg", ".patchwork", "checkpoints.json")
	if p != want {
		t.Fatalf("got %s want %s", p, want)
	}
}

func TestLoadCheckpoints_NotExist(t *testing.T) {
	dir := t.TempDir()
	cps, err := LoadCheckpoints(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cps) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(cps))
	}
}

func TestCreateCheckpoint_Basic(t *testing.T) {
	dir := t.TempDir()
	cp, err := CreateCheckpoint(dir, "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp.Name != "v1" {
		t.Fatalf("expected name v1, got %s", cp.Name)
	}
	if cp.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}

func TestCreateCheckpoint_DuplicateReturnsError(t *testing.T) {
	dir := t.TempDir()
	if _, err := CreateCheckpoint(dir, "v1"); err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err := CreateCheckpoint(dir, "v1")
	if err == nil {
		t.Fatal("expected error for duplicate checkpoint name")
	}
}

func TestCreateCheckpoint_EmptyNameReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := CreateCheckpoint(dir, "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCreateCheckpoint_IncludesAppliedPatches(t *testing.T) {
	dir := t.TempDir()
	// seed a history entry
	h := &History{Applied: map[string]HistoryEntry{
		"patch-001": {PatchID: "patch-001"},
	}}
	data, _ := historyToJSON(h)
	hp := HistoryPath(dir)
	os.MkdirAll(filepath.Dir(hp), 0o755)
	os.WriteFile(hp, data, 0o644)

	cp, err := CreateCheckpoint(dir, "release-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cp.Applied) != 1 || cp.Applied[0] != "patch-001" {
		t.Fatalf("expected applied=[patch-001], got %v", cp.Applied)
	}
}

func TestFindCheckpoint_Found(t *testing.T) {
	dir := t.TempDir()
	if _, err := CreateCheckpoint(dir, "beta"); err != nil {
		t.Fatalf("create: %v", err)
	}
	cp, ok, err := FindCheckpoint(dir, "beta")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if !ok {
		t.Fatal("expected checkpoint to be found")
	}
	if cp.Name != "beta" {
		t.Fatalf("expected name beta, got %s", cp.Name)
	}
}

func TestFindCheckpoint_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, ok, err := FindCheckpoint(dir, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected checkpoint to not be found")
	}
}
