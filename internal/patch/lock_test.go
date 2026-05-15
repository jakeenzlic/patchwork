package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func baseLockPatch() Patch {
	return Patch{
		ID:      "lock-test-001",
		Version: "1.0.0",
		Ops: []Op{
			{Op: "add", Path: "feature/enabled", Value: true},
		},
	}
}

func TestLockPath(t *testing.T) {
	dir := "/tmp/cfg"
	got := LockPath(dir)
	want := filepath.Join(dir, ".patchwork.lock")
	if got != want {
		t.Fatalf("LockPath: got %s, want %s", got, want)
	}
}

func TestLoadLockFile_NotExist(t *testing.T) {
	dir := t.TempDir()
	lf, err := LoadLockFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf.Entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(lf.Entries))
	}
}

func TestPin_CreatesEntry(t *testing.T) {
	dir := t.TempDir()
	p := baseLockPatch()
	if err := Pin(dir, p); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	lf, err := LoadLockFile(dir)
	if err != nil {
		t.Fatalf("LoadLockFile: %v", err)
	}
	entry, ok := lf.Entries[p.ID]
	if !ok {
		t.Fatal("expected entry for patch ID")
	}
	if entry.Version != p.Version {
		t.Fatalf("version: got %s, want %s", entry.Version, p.Version)
	}
	if entry.Checksum == "" {
		t.Fatal("expected non-empty checksum")
	}
}

func TestCheckLock_Valid(t *testing.T) {
	dir := t.TempDir()
	p := baseLockPatch()
	if err := Pin(dir, p); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	if err := CheckLock(dir, p); err != nil {
		t.Fatalf("CheckLock should pass for unchanged patch: %v", err)
	}
}

func TestCheckLock_Mismatch(t *testing.T) {
	dir := t.TempDir()
	p := baseLockPatch()
	if err := Pin(dir, p); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	// Mutate the patch after pinning
	p.Ops[0].Value = false
	if err := CheckLock(dir, p); err == nil {
		t.Fatal("expected checksum mismatch error")
	}
}

func TestCheckLock_NoPinReturnsNil(t *testing.T) {
	dir := t.TempDir()
	p := baseLockPatch()
	if err := CheckLock(dir, p); err != nil {
		t.Fatalf("CheckLock with no pin should return nil: %v", err)
	}
}

func TestSaveLockFile_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	lf := &LockFile{Entries: make(map[string]LockEntry)}
	if err := SaveLockFile(dir, lf); err != nil {
		t.Fatalf("SaveLockFile: %v", err)
	}
	if _, err := os.Stat(LockPath(dir)); err != nil {
		t.Fatalf("lock file not created: %v", err)
	}
}
