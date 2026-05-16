package patch

import (
	"testing"
)

func baseAnnotatePatch() Patch {
	return Patch{
		ID:      "p-annotate-001",
		Version: "1.0",
		Ops: []Op{
			{Op: "replace", Path: "server/port", Value: 9090},
			{Op: "add", Path: "feature/flag", Value: true},
		},
	}
}

func TestAnnotate_Valid(t *testing.T) {
	p := baseAnnotatePatch()
	anns := []Annotation{
		{PatchID: "p-annotate-001", Path: "server/port", Note: "bumped for prod", Author: "alice"},
	}
	result, err := Annotate(p, anns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(result))
	}
	if result[0].Note != "bumped for prod" {
		t.Errorf("unexpected note: %s", result[0].Note)
	}
}

func TestAnnotate_PathNotFound(t *testing.T) {
	p := baseAnnotatePatch()
	anns := []Annotation{
		{PatchID: "p-annotate-001", Path: "nonexistent/path", Note: "some note"},
	}
	_, err := Annotate(p, anns)
	if err == nil {
		t.Fatal("expected error for unknown path")
	}
}

func TestAnnotate_WrongPatchID(t *testing.T) {
	p := baseAnnotatePatch()
	anns := []Annotation{
		{PatchID: "wrong-id", Path: "server/port", Note: "note"},
	}
	_, err := Annotate(p, anns)
	if err == nil {
		t.Fatal("expected error for mismatched patch_id")
	}
}

func TestAnnotate_EmptyNote(t *testing.T) {
	p := baseAnnotatePatch()
	anns := []Annotation{
		{PatchID: "p-annotate-001", Path: "server/port", Note: "  "},
	}
	_, err := Annotate(p, anns)
	if err == nil {
		t.Fatal("expected error for empty note")
	}
}

func TestFormatAnnotations_Empty(t *testing.T) {
	out := FormatAnnotations(nil)
	if out != "no annotations" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatAnnotations_ContainsNote(t *testing.T) {
	anns := []Annotation{
		{PatchID: "p-001", Path: "db/host", Note: "updated for staging", Author: "bob"},
	}
	out := FormatAnnotations(anns)
	for _, want := range []string{"p-001", "db/host", "updated for staging", "bob"} {
		if !contains(out, want) {
			t.Errorf("expected %q in output: %s", want, out)
		}
	}
}

func TestFormatAnnotations_DefaultAuthor(t *testing.T) {
	anns := []Annotation{
		{PatchID: "p-001", Path: "db/host", Note: "note"},
	}
	out := FormatAnnotations(anns)
	if !contains(out, "unknown") {
		t.Errorf("expected 'unknown' author in output: %s", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
