package patch

import (
	"testing"
)

func TestLint_NoWarnings(t *testing.T) {
	p := &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "add", Path: "server/port", Value: 9090},
			{Op: "replace", Path: "server/host", Value: "localhost"},
		},
	}
	warnings := Lint(p)
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got: %v", warnings)
	}
}

func TestLint_DuplicatePath(t *testing.T) {
	p := &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "add", Path: "server/port", Value: 8080},
			{Op: "add", Path: "server/port", Value: 9090},
		},
	}
	warnings := Lint(p)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning for duplicate path, got %d: %v", len(warnings), warnings)
	}
}

func TestLint_LeadingSlash(t *testing.T) {
	p := &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "add", Path: "/server/port", Value: 8080},
		},
	}
	warnings := Lint(p)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning for leading slash, got %d", len(warnings))
	}
}

func TestLint_TrailingSlash(t *testing.T) {
	p := &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "replace", Path: "server/port/", Value: 8080},
		},
	}
	warnings := Lint(p)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning for trailing slash, got %d", len(warnings))
	}
}

func TestLint_EmptySegment(t *testing.T) {
	p := &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "add", Path: "server//port", Value: 8080},
		},
	}
	warnings := Lint(p)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning for empty segment, got %d", len(warnings))
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	p := &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "add", Path: "/server//port/", Value: 1},
			{Op: "add", Path: "/server//port/", Value: 2},
		},
	}
	warnings := Lint(p)
	if len(warnings) < 3 {
		t.Fatalf("expected at least 3 warnings, got %d: %v", len(warnings), warnings)
	}
}
