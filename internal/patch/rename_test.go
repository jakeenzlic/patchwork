package patch

import (
	"testing"
)

func makeRenamePatch(id string, ops []Op) Patch {
	return Patch{ID: id, Version: "1.0", Ops: ops}
}

func TestRename_RewritesMatchingOps(t *testing.T) {
	p := makeRenamePatch("p1", []Op{
		{Op: "replace", Path: "database/host", Value: "localhost"},
		{Op: "replace", Path: "database/port", Value: 5432},
		{Op: "replace", Path: "cache/host", Value: "redis"},
	})

	updated, result, err := Rename(p, "database", "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Renamed {
		t.Fatal("expected Renamed=true")
	}
	if updated.Ops[0].Path != "db/host" {
		t.Errorf("expected db/host, got %s", updated.Ops[0].Path)
	}
	if updated.Ops[1].Path != "db/port" {
		t.Errorf("expected db/port, got %s", updated.Ops[1].Path)
	}
	if updated.Ops[2].Path != "cache/host" {
		t.Errorf("cache/host should be unchanged, got %s", updated.Ops[2].Path)
	}
}

func TestRename_NoMatchReturnsRenamedFalse(t *testing.T) {
	p := makeRenamePatch("p2", []Op{
		{Op: "add", Path: "app/name", Value: "patchwork"},
	})

	_, result, err := Rename(p, "database", "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Renamed {
		t.Error("expected Renamed=false when no ops match")
	}
}

func TestRename_EmptyOldPrefixReturnsError(t *testing.T) {
	p := makeRenamePatch("p3", []Op{})
	_, _, err := Rename(p, "", "db")
	if err == nil {
		t.Fatal("expected error for empty oldPrefix")
	}
}

func TestRename_EmptyNewPrefixReturnsError(t *testing.T) {
	p := makeRenamePatch("p4", []Op{})
	_, _, err := Rename(p, "database", "")
	if err == nil {
		t.Fatal("expected error for empty newPrefix")
	}
}

func TestRename_OriginalUnmutated(t *testing.T) {
	p := makeRenamePatch("p5", []Op{
		{Op: "replace", Path: "database/host", Value: "localhost"},
	})

	_, _, err := Rename(p, "database", "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Ops[0].Path != "database/host" {
		t.Errorf("original patch mutated: got %s", p.Ops[0].Path)
	}
}

func TestFormatRenameResult_Renamed(t *testing.T) {
	r := RenameResult{OldPath: "database", NewPath: "db", PatchID: "p1", Renamed: true}
	out := FormatRenameResult(r)
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormatRenameResult_NotRenamed(t *testing.T) {
	r := RenameResult{OldPath: "database", NewPath: "db", PatchID: "p1", Renamed: false}
	out := FormatRenameResult(r)
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestPathHasPrefix_Exact(t *testing.T) {
	if !pathHasPrefix("database", "database") {
		t.Error("exact match should return true")
	}
}

func TestPathHasPrefix_Nested(t *testing.T) {
	if !pathHasPrefix("database/host", "database") {
		t.Error("nested path should match prefix")
	}
}

func TestPathHasPrefix_NoPrefixMatch(t *testing.T) {
	if pathHasPrefix("databasex/host", "database") {
		t.Error("partial word match should not count as prefix")
	}
}
