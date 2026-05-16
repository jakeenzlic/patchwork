package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func writeStashConfig(t *testing.T, dir string) string {
	t.Helper()
	cfgPath := filepath.Join(dir, "config.json")
	data := `{"version": "1", "debug": false}`
	if err := os.WriteFile(cfgPath, []byte(data), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return cfgPath
}

func TestStashPath(t *testing.T) {
	p := StashPath("/etc/app/config.json")
	want := "/etc/app/.patchwork/config.json.stash.json"
	if p != want {
		t.Errorf("got %q, want %q", p, want)
	}
}

func TestLoadStash_NotExist(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	entries, err := LoadStash(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty stash, got %d entries", len(entries))
	}
}

func TestStash_CreatesEntry(t *testing.T) {
	dir := t.TempDir()
	cfgPath := writeStashConfig(t, dir)

	entry, err := Stash(cfgPath, "before migration")
	if err != nil {
		t.Fatalf("Stash error: %v", err)
	}
	if entry.Message != "before migration" {
		t.Errorf("expected message %q, got %q", "before migration", entry.Message)
	}
	if entry.ID != "stash@{0}" {
		t.Errorf("unexpected ID %q", entry.ID)
	}
}

func TestStash_Accumulates(t *testing.T) {
	dir := t.TempDir()
	cfgPath := writeStashConfig(t, dir)

	if _, err := Stash(cfgPath, "first"); err != nil {
		t.Fatalf("first stash: %v", err)
	}
	if _, err := Stash(cfgPath, "second"); err != nil {
		t.Fatalf("second stash: %v", err)
	}

	entries, err := LoadStash(cfgPath)
	if err != nil {
		t.Fatalf("LoadStash: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[1].ID != "stash@{1}" {
		t.Errorf("unexpected second ID %q", entries[1].ID)
	}
}

func TestStashPop_RestoresConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := writeStashConfig(t, dir)

	if _, err := Stash(cfgPath, "snapshot"); err != nil {
		t.Fatalf("Stash: %v", err)
	}

	entry, err := StashPop(cfgPath)
	if err != nil {
		t.Fatalf("StashPop: %v", err)
	}
	if entry.Message != "snapshot" {
		t.Errorf("expected message %q, got %q", "snapshot", entry.Message)
	}

	remaining, _ := LoadStash(cfgPath)
	if len(remaining) != 0 {
		t.Errorf("expected empty stash after pop, got %d", len(remaining))
	}
}

func TestStashPop_EmptyReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	_, err := StashPop(cfgPath)
	if err == nil {
		t.Error("expected error on empty stash pop")
	}
}

func TestFormatStash_Empty(t *testing.T) {
	out := FormatStash([]StashEntry{})
	if out != "No stash entries.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
