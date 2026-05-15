package patch

import (
	"strings"
	"testing"
)

func makeConflictPatch(id string, ops []Op) Patch {
	return Patch{ID: id, Version: 1, Ops: ops}
}

func TestDetectConflicts_NoConflict(t *testing.T) {
	a := makeConflictPatch("a", []Op{{Op: "add", Path: "server.port", Value: 8080}})
	b := makeConflictPatch("b", []Op{{Op: "add", Path: "server.host", Value: "localhost"}})

	got := DetectConflicts([]Patch{a, b})
	if len(got) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(got))
	}
}

func TestDetectConflicts_BothReplace(t *testing.T) {
	a := makeConflictPatch("a", []Op{{Op: "replace", Path: "server.port", Value: 8080}})
	b := makeConflictPatch("b", []Op{{Op: "replace", Path: "server.port", Value: 9090}})

	got := DetectConflicts([]Patch{a, b})
	if len(got) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(got))
	}
	if got[0].PatchA != "a" || got[0].PatchB != "b" {
		t.Errorf("unexpected patch IDs: %+v", got[0])
	}
	if !strings.Contains(got[0].Reason, "replace") {
		t.Errorf("reason should mention replace, got: %s", got[0].Reason)
	}
}

func TestDetectConflicts_DeleteThenReplace(t *testing.T) {
	a := makeConflictPatch("a", []Op{{Op: "delete", Path: "feature.flag"}})
	b := makeConflictPatch("b", []Op{{Op: "replace", Path: "feature.flag", Value: true}})

	got := DetectConflicts([]Patch{a, b})
	if len(got) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(got))
	}
	if !strings.Contains(got[0].Reason, "delete") {
		t.Errorf("reason should mention delete, got: %s", got[0].Reason)
	}
}

func TestDetectConflicts_BothDelete(t *testing.T) {
	a := makeConflictPatch("a", []Op{{Op: "delete", Path: "legacy.key"}})
	b := makeConflictPatch("b", []Op{{Op: "delete", Path: "legacy.key"}})

	got := DetectConflicts([]Patch{a, b})
	if len(got) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(got))
	}
}

func TestDetectConflicts_MultipleOps(t *testing.T) {
	a := makeConflictPatch("a", []Op{
		{Op: "add", Path: "db.host", Value: "primary"},
		{Op: "replace", Path: "db.port", Value: 5432},
	})
	b := makeConflictPatch("b", []Op{
		{Op: "replace", Path: "db.port", Value: 5433},
	})

	got := DetectConflicts([]Patch{a, b})
	if len(got) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(got))
	}
	if got[0].Path != "db.port" {
		t.Errorf("expected conflict on db.port, got %s", got[0].Path)
	}
}

func TestFormatConflicts_Empty(t *testing.T) {
	out := FormatConflicts(nil)
	if !strings.Contains(out, "no conflicts") {
		t.Errorf("expected no-conflict message, got: %s", out)
	}
}

func TestFormatConflicts_NonEmpty(t *testing.T) {
	conflicts := []Conflict{
		{PatchA: "a", PatchB: "b", Path: "x.y", Reason: "both replace"},
	}
	out := FormatConflicts(conflicts)
	if !strings.Contains(out, "1 conflict") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "x.y") {
		t.Errorf("expected path in output, got: %s", out)
	}
}
