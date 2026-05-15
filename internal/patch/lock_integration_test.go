package patch_test

import (
	"path/filepath"
	"testing"

	"patchwork/internal/patch"
)

// TestLock_PinAndCheckRoundtrip ensures Pin + CheckLock work end-to-end.
func TestLock_PinAndCheckRoundtrip(t *testing.T) {
	dir := t.TempDir()
	p := patch.Patch{
		ID:      "int-lock-001",
		Version: "2.1.0",
		Ops: []patch.Op{
			{Op: "replace", Path: "db/pool", Value: 10},
		},
	}

	if err := patch.Pin(dir, p); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	if err := patch.CheckLock(dir, p); err != nil {
		t.Fatalf("CheckLock after Pin: %v", err)
	}
}

// TestLock_MultiplePatches verifies independent entries per patch ID.
func TestLock_MultiplePatches(t *testing.T) {
	dir := t.TempDir()
	patches := []patch.Patch{
		{ID: "p1", Version: "1.0.0", Ops: []patch.Op{{Op: "add", Path: "a", Value: 1}}},
		{ID: "p2", Version: "1.0.0", Ops: []patch.Op{{Op: "add", Path: "b", Value: 2}}},
	}
	for _, p := range patches {
		if err := patch.Pin(dir, p); err != nil {
			t.Fatalf("Pin %s: %v", p.ID, err)
		}
	}
	lf, err := patch.LoadLockFile(dir)
	if err != nil {
		t.Fatalf("LoadLockFile: %v", err)
	}
	if len(lf.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lf.Entries))
	}
}

// TestLock_RePinUpdatesChecksum verifies that re-pinning a mutated patch updates the entry.
func TestLock_RePinUpdatesChecksum(t *testing.T) {
	dir := t.TempDir()
	p := patch.Patch{
		ID:      "repin-001",
		Version: "1.0.0",
		Ops:     []patch.Op{{Op: "add", Path: "x", Value: 1}},
	}
	if err := patch.Pin(dir, p); err != nil {
		t.Fatalf("first Pin: %v", err)
	}

	// Mutate and re-pin
	p.Ops[0].Value = 99
	if err := patch.Pin(dir, p); err != nil {
		t.Fatalf("second Pin: %v", err)
	}

	// CheckLock should now pass with the new value
	if err := patch.CheckLock(dir, p); err != nil {
		t.Fatalf("CheckLock after re-pin: %v", err)
	}
}

// TestLock_FileLocation confirms the lock file is at the expected path.
func TestLock_FileLocation(t *testing.T) {
	dir := t.TempDir()
	p := patch.Patch{
		ID: "loc-001", Version: "1.0.0",
		Ops: []patch.Op{{Op: "add", Path: "k", Value: true}},
	}
	if err := patch.Pin(dir, p); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	expected := filepath.Join(dir, ".patchwork.lock")
	if got := patch.LockPath(dir); got != expected {
		t.Fatalf("LockPath: got %s, want %s", got, expected)
	}
}
