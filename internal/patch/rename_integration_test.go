package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRenamePatch(t *testing.T, dir, name string, p Patch) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := Export(p, path); err != nil {
		t.Fatalf("writeRenamePatch: %v", err)
	}
	return path
}

func TestRename_ExportRoundtrip(t *testing.T) {
	dir := t.TempDir()

	p := makeRenamePatch("rename-rt", []Op{
		{Op: "replace", Path: "service/host", Value: "localhost"},
		{Op: "add", Path: "service/port", Value: 8080},
	})

	src := writeRenamePatch(t, dir, "rename-rt.json", p)

	loaded, err := LoadFromFile(src)
	if err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}

	updated, result, err := Rename(loaded, "service", "svc")
	if err != nil {
		t.Fatalf("Rename: %v", err)
	}
	if !result.Renamed {
		t.Fatal("expected rename to succeed")
	}

	out := filepath.Join(dir, "renamed.json")
	if err := Export(updated, out); err != nil {
		t.Fatalf("Export: %v", err)
	}

	reloaded, err := LoadFromFile(out)
	if err != nil {
		t.Fatalf("LoadFromFile reloaded: %v", err)
	}
	if reloaded.Ops[0].Path != "svc/host" {
		t.Errorf("expected svc/host, got %s", reloaded.Ops[0].Path)
	}
	if reloaded.Ops[1].Path != "svc/port" {
		t.Errorf("expected svc/port, got %s", reloaded.Ops[1].Path)
	}
}

func TestRename_PreservesUnrelatedOps(t *testing.T) {
	dir := t.TempDir()

	p := makeRenamePatch("rename-preserve", []Op{
		{Op: "replace", Path: "alpha/key", Value: "v1"},
		{Op: "replace", Path: "beta/key", Value: "v2"},
		{Op: "delete", Path: "alpha/old"},
	})

	_ = dir
	updated, result, err := Rename(p, "alpha", "a")
	if err != nil {
		t.Fatalf("Rename: %v", err)
	}
	if !result.Renamed {
		t.Fatal("expected rename")
	}
	if updated.Ops[1].Path != "beta/key" {
		t.Errorf("beta/key should be unchanged, got %s", updated.Ops[1].Path)
	}
	if updated.Ops[0].Path != "a/key" {
		t.Errorf("expected a/key, got %s", updated.Ops[0].Path)
	}
	if updated.Ops[2].Path != "a/old" {
		t.Errorf("expected a/old, got %s", updated.Ops[2].Path)
	}
}

func TestRename_FileNotFound(t *testing.T) {
	_, err := LoadFromFile(filepath.Join(os.TempDir(), "nonexistent-rename.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
