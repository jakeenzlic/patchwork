package patch

import (
	"strings"
	"testing"
)

func makeTrimPatch(id string, ops []Op) Patch {
	return Patch{ID: id, Version: "1.0", Ops: ops}
}

var baseTrimConfig = map[string]any{
	"database": map[string]any{
		"host": "localhost",
		"port": 5432,
	},
	"app": map[string]any{
		"debug": true,
	},
}

func TestTrim_KeepsValidPaths(t *testing.T) {
	p := makeTrimPatch("p1", []Op{
		{Op: "replace", Path: "database.host", Value: "prod-db"},
	})
	trimmed, result, err := Trim(p, baseTrimConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
	if len(trimmed.Ops) != 1 {
		t.Errorf("expected 1 op kept, got %d", len(trimmed.Ops))
	}
}

func TestTrim_RemovesStalePaths(t *testing.T) {
	p := makeTrimPatch("p2", []Op{
		{Op: "replace", Path: "database.host", Value: "prod-db"},
		{Op: "delete", Path: "legacy.setting"},
	})
	trimmed, result, err := Trim(p, baseTrimConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed))
	}
	if result.Removed[0] != "legacy.setting" {
		t.Errorf("unexpected removed path: %s", result.Removed[0])
	}
	if len(trimmed.Ops) != 1 {
		t.Errorf("expected 1 op in trimmed patch, got %d", len(trimmed.Ops))
	}
}

func TestTrim_AddOpsAlwaysKept(t *testing.T) {
	p := makeTrimPatch("p3", []Op{
		{Op: "add", Path: "new.feature", Value: true},
	})
	trimmed, result, err := Trim(p, baseTrimConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 0 {
		t.Errorf("add op should never be removed")
	}
	if len(trimmed.Ops) != 1 {
		t.Errorf("expected 1 op kept, got %d", len(trimmed.Ops))
	}
}

func TestTrim_InvalidPatchReturnsError(t *testing.T) {
	p := Patch{ID: "", Version: "", Ops: []Op{{Op: "replace", Path: "x"}}}
	_, _, err := Trim(p, baseTrimConfig)
	if err == nil {
		t.Error("expected error for invalid patch")
	}
}

func TestFormatTrimResult_NoRemovals(t *testing.T) {
	r := TrimResult{Kept: []string{"database.host"}}
	out := FormatTrimResult(r)
	if !strings.Contains(out, "no stale") {
		t.Errorf("expected 'no stale' message, got: %s", out)
	}
}

func TestFormatTrimResult_WithRemovals(t *testing.T) {
	r := TrimResult{
		Removed: []string{"legacy.key", "old.setting"},
		Kept:    []string{"database.host"},
	}
	out := FormatTrimResult(r)
	if !strings.Contains(out, "legacy.key") {
		t.Errorf("expected removed path in output, got: %s", out)
	}
	if !strings.Contains(out, "2 stale") {
		t.Errorf("expected count in output, got: %s", out)
	}
}
