package patch

import (
	"testing"
)

func baseScopePatch(id, path, op string) Patch {
	return Patch{
		ID:      id,
		Version: "1.0",
		Operations: []Operation{
			{Op: op, Path: path, Value: "v"},
		},
	}
}

func TestNewScope_Valid(t *testing.T) {
	s, err := NewScope("database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Prefix != "database" {
		t.Errorf("expected prefix 'database', got %q", s.Prefix)
	}
}

func TestNewScope_EmptyPrefix(t *testing.T) {
	_, err := NewScope("")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNewScope_LeadingSlash(t *testing.T) {
	_, err := NewScope("/database")
	if err == nil {
		t.Fatal("expected error for leading slash")
	}
}

func TestInScope_DirectMatch(t *testing.T) {
	if !InScope("db", "db") {
		t.Error("expected 'db' to be in scope 'db'")
	}
}

func TestInScope_NestedPath(t *testing.T) {
	if !InScope("db", "db/host") {
		t.Error("expected 'db/host' to be in scope 'db'")
	}
}

func TestInScope_UnrelatedPath(t *testing.T) {
	if InScope("db", "app/name") {
		t.Error("expected 'app/name' NOT to be in scope 'db'")
	}
}

func TestInScope_PrefixSubstring(t *testing.T) {
	// 'database' should not match scope 'db'
	if InScope("db", "database/host") {
		t.Error("expected 'database/host' NOT to be in scope 'db'")
	}
}

func TestScope_Filter_KeepsMatchingOps(t *testing.T) {
	s, _ := NewScope("db")
	patches := []Patch{
		baseScopePatch("p1", "db/host", OpReplace),
		baseScopePatch("p2", "app/name", OpReplace),
	}
	result := s.Filter(patches)
	if len(result) != 1 {
		t.Fatalf("expected 1 patch, got %d", len(result))
	}
	if result[0].ID != "p1" {
		t.Errorf("expected patch p1, got %s", result[0].ID)
	}
}

func TestScope_Filter_NoMatch(t *testing.T) {
	s, _ := NewScope("db")
	patches := []Patch{
		baseScopePatch("p1", "app/name", OpAdd),
	}
	result := s.Filter(patches)
	if len(result) != 0 {
		t.Errorf("expected 0 patches, got %d", len(result))
	}
}

func TestFormatScope_Output(t *testing.T) {
	out := FormatScope("db", 10, 4)
	for _, want := range []string{"db", "10", "4", "6"} {
		if !containsStr(out, want) {
			t.Errorf("FormatScope output missing %q: %s", want, out)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
