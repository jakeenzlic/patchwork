package patch

import (
	"strings"
	"testing"
)

func baseTagPatch(id, tag string) Patch {
	return Patch{
		ID:      id,
		Version: "1.0",
		Tags:    []string{tag},
		Ops:     []Op{{Op: "add", Path: "key", Value: "v"}},
	}
}

func TestValidateTag_Valid(t *testing.T) {
	for _, tag := range []string{"release", "v1", "hot-fix", "env_prod"} {
		if err := ValidateTag(tag); err != nil {
			t.Errorf("expected valid tag %q, got error: %v", tag, err)
		}
	}
}

func TestValidateTag_Empty(t *testing.T) {
	if err := ValidateTag(""); err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestValidateTag_InvalidChars(t *testing.T) {
	for _, tag := range []string{"has space", "dot.tag", "slash/tag"} {
		if err := ValidateTag(tag); err == nil {
			t.Errorf("expected error for tag %q", tag)
		}
	}
}

func TestBuildTagIndex(t *testing.T) {
	patches := []Patch{
		baseTagPatch("p1", "release"),
		baseTagPatch("p2", "release"),
		baseTagPatch("p3", "hotfix"),
	}
	idx, err := BuildTagIndex(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx["release"]) != 2 {
		t.Errorf("expected 2 entries for 'release', got %d", len(idx["release"]))
	}
	if len(idx["hotfix"]) != 1 {
		t.Errorf("expected 1 entry for 'hotfix', got %d", len(idx["hotfix"]))
	}
}

func TestBuildTagIndex_InvalidTag(t *testing.T) {
	patches := []Patch{baseTagPatch("p1", "bad tag")}
	_, err := BuildTagIndex(patches)
	if err == nil {
		t.Error("expected error for invalid tag in patch")
	}
}

func TestFilterByTag(t *testing.T) {
	patches := []Patch{
		baseTagPatch("p1", "release"),
		baseTagPatch("p2", "hotfix"),
		baseTagPatch("p3", "release"),
	}
	result := FilterByTag(patches, "release")
	if len(result) != 2 {
		t.Errorf("expected 2 patches, got %d", len(result))
	}
}

func TestFormatTagIndex(t *testing.T) {
	idx := TagIndex{
		"release": {"p1", "p2"},
		"hotfix":  {"p3"},
	}
	out := FormatTagIndex(idx)
	if !strings.Contains(out, "release") || !strings.Contains(out, "hotfix") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatTagIndex_Empty(t *testing.T) {
	out := FormatTagIndex(TagIndex{})
	if out != "(no tags)" {
		t.Errorf("expected '(no tags)', got %q", out)
	}
}
