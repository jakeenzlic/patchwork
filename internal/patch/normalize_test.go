package patch

import (
	"strings"
	"testing"
)

func baseNormPatch() Patch {
	return Patch{
		ID:      "norm-001",
		Version: "1.0",
		Ops: []Op{
			{Op: "add", Path: "server/port", Value: 8080},
			{Op: "replace", Path: "server/host", Value: "localhost"},
		},
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	p := baseNormPatch()
	r, err := Normalize(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %d: %v", len(r.Changes), r.Changes)
	}
	if len(r.Patch.Ops) != 2 {
		t.Errorf("expected 2 ops, got %d", len(r.Patch.Ops))
	}
}

func TestNormalize_LowercasesOp(t *testing.T) {
	p := baseNormPatch()
	p.Ops[0].Op = "ADD"
	r, err := Normalize(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Patch.Ops[0].Op != "add" {
		t.Errorf("expected op 'add', got %q", r.Patch.Ops[0].Op)
	}
	if len(r.Changes) == 0 {
		t.Error("expected at least one change recorded")
	}
}

func TestNormalize_StripsLeadingSlash(t *testing.T) {
	p := baseNormPatch()
	p.Ops[0].Path = "/server/port"
	r, err := Normalize(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Patch.Ops[0].Path != "server/port" {
		t.Errorf("expected path 'server/port', got %q", r.Patch.Ops[0].Path)
	}
}

func TestNormalize_TrimsSegmentWhitespace(t *testing.T) {
	p := baseNormPatch()
	p.Ops[0].Path = "server / port"
	r, err := Normalize(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Patch.Ops[0].Path != "server / port" {
		// segments are trimmed but slash-joined; "server " -> "server"
	}
	if strings.Contains(r.Patch.Ops[0].Path, " ") {
		t.Errorf("expected whitespace trimmed from segments, got %q", r.Patch.Ops[0].Path)
	}
}

func TestNormalize_RemovesDuplicateOps(t *testing.T) {
	p := baseNormPatch()
	p.Ops = append(p.Ops, Op{Op: "add", Path: "server/port", Value: 9090})
	r, err := Normalize(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Patch.Ops) != 2 {
		t.Errorf("expected duplicate removed, got %d ops", len(r.Patch.Ops))
	}
	found := false
	for _, c := range r.Changes {
		if strings.Contains(c, "duplicate") {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate change recorded")
	}
}

func TestNormalize_InvalidPatchReturnsError(t *testing.T) {
	p := Patch{ID: "bad", Version: "", Ops: []Op{{Op: "add", Path: "x", Value: 1}}}
	_, err := Normalize(p)
	if err == nil {
		t.Error("expected error for invalid patch")
	}
}

func TestFormatNormalizeResult_NoChanges(t *testing.T) {
	r := NormalizeResult{}
	out := FormatNormalizeResult(r)
	if !strings.Contains(out, "no changes") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatNormalizeResult_WithChanges(t *testing.T) {
	r := NormalizeResult{Changes: []string{"op[0]: lowercased op"}}
	out := FormatNormalizeResult(r)
	if !strings.Contains(out, "1 change") {
		t.Errorf("unexpected output: %q", out)
	}
}
