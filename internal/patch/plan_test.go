package patch

import (
	"strings"
	"testing"
)

func basePlanPatch() *Patch {
	return &Patch{
		Version: "v1.2.0",
		Ops: []Op{
			{Op: "add", Path: "server/timeout", Value: 30},
			{Op: "replace", Path: "server/host", Value: "prod.example.com"},
			{Op: "delete", Path: "debug"},
		},
	}
}

func TestPlan_ReturnsPendingEntries(t *testing.T) {
	p := basePlanPatch()
	target := map[string]any{"server": map[string]any{"host": "localhost"}, "debug": true}

	entries, err := Plan(p, target, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Applied {
			t.Errorf("expected Applied=false for entry %d", e.Index)
		}
	}
}

func TestPlan_MarksAlreadyApplied(t *testing.T) {
	p := basePlanPatch()
	target := map[string]any{}

	h := &History{}
	h.Record(p.Version)

	entries, err := Plan(p, target, h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if !e.Applied {
			t.Errorf("expected Applied=true for entry %d", e.Index)
		}
	}
}

func TestPlan_InvalidPatchReturnsError(t *testing.T) {
	p := &Patch{Version: "", Ops: []Op{{Op: "add", Path: "x", Value: 1}}}
	_, err := Plan(p, map[string]any{}, nil)
	if err == nil {
		t.Fatal("expected error for invalid patch")
	}
}

func TestFormatPlan_ContainsOpAndPath(t *testing.T) {
	entries := []PlanEntry{
		{Index: 0, Op: "add", Path: "foo/bar", Value: 42, Applied: false},
		{Index: 1, Op: "delete", Path: "baz", Applied: true},
	}
	out := FormatPlan(entries)
	if !strings.Contains(out, "add") {
		t.Error("expected 'add' in output")
	}
	if !strings.Contains(out, "foo/bar") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "already applied") {
		t.Error("expected 'already applied' status in output")
	}
	if !strings.Contains(out, "pending") {
		t.Error("expected 'pending' status in output")
	}
}

func TestFormatPlan_EmptyEntries(t *testing.T) {
	out := FormatPlan(nil)
	if !strings.Contains(out, "no operations") {
		t.Errorf("expected empty message, got: %q", out)
	}
}
