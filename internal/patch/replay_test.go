package patch

import (
	"testing"
)

func makeReplayPatch(id, op, path string, value any) Patch {
	return Patch{
		ID:      id,
		Version: "1.0",
		Ops: []Op{
			{Op: op, Path: path, Value: value},
		},
	}
}

func TestReplay_AppliesAllPatches(t *testing.T) {
	base := map[string]any{"env": "dev"}
	patches := []Patch{
		makeReplayPatch("p1", "add", "region", "us-east-1"),
		makeReplayPatch("p2", "replace", "env", "prod"),
	}

	res, err := Replay(base, patches, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 2 {
		t.Fatalf("expected 2 applied, got %d", len(res.Applied))
	}
	if res.Final["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", res.Final["env"])
	}
	if res.Final["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", res.Final["region"])
	}
}

func TestReplay_StopsAtTarget(t *testing.T) {
	base := map[string]any{}
	patches := []Patch{
		makeReplayPatch("p1", "add", "a", "1"),
		makeReplayPatch("p2", "add", "b", "2"),
		makeReplayPatch("p3", "add", "c", "3"),
	}

	res, err := Replay(base, patches, "p2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 2 {
		t.Fatalf("expected 2 applied, got %d", len(res.Applied))
	}
	if _, ok := res.Final["c"]; ok {
		t.Error("expected patch p3 to not be applied")
	}
}

func TestReplay_SkipsFailingPatch(t *testing.T) {
	base := map[string]any{"x": "1"}
	patches := []Patch{
		makeReplayPatch("p1", "replace", "nonexistent", "val"),
		makeReplayPatch("p2", "add", "y", "2"),
	}

	res, err := Replay(base, patches, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "p1" {
		t.Errorf("expected p1 skipped, got %v", res.Skipped)
	}
	if res.Final["y"] != "2" {
		t.Errorf("expected y=2, got %v", res.Final["y"])
	}
}

func TestReplay_InvalidPatchReturnsError(t *testing.T) {
	base := map[string]any{}
	bad := Patch{ID: "bad", Version: "", Ops: []Op{{Op: "add", Path: "k", Value: "v"}}}

	_, err := Replay(base, []Patch{bad}, "")
	if err == nil {
		t.Fatal("expected error for invalid patch")
	}
}

func TestFormatReplayResult_Output(t *testing.T) {
	res := ReplayResult{
		Applied: []string{"p1", "p2"},
		Skipped: []string{"p3"},
	}
	out := FormatReplayResult(res)
	if out == "" {
		t.Error("expected non-empty output")
	}
	for _, id := range []string{"p1", "p2", "p3"} {
		if !contains(out, id) {
			t.Errorf("expected output to contain %q", id)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
