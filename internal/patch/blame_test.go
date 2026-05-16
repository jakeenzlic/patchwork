package patch

import (
	"strings"
	"testing"
)

func makeBlamePatches() []Patch {
	return []Patch{
		{
			ID:      "patch-001",
			Version: "1.0",
			Ops:     []Op{{Op: "add", Path: "server/port", Value: 8080}},
		},
		{
			ID:      "patch-002",
			Version: "1.0",
			Ops:     []Op{{Op: "replace", Path: "server/port", Value: 9090}},
		},
		{
			ID:      "patch-003",
			Version: "1.0",
			Ops:     []Op{{Op: "add", Path: "db/host", Value: "localhost"}},
		},
	}
}

func makeBlameLog() []AuditEntry {
	return []AuditEntry{
		{PatchID: "patch-001", Status: "success", AppliedAt: "2024-01-01T10:00:00Z"},
		{PatchID: "patch-002", Status: "success", AppliedAt: "2024-01-02T10:00:00Z"},
		{PatchID: "patch-003", Status: "success", AppliedAt: "2024-01-03T10:00:00Z"},
	}
}

func TestBlame_LatestPatchWins(t *testing.T) {
	patches := makeBlamePatches()
	log := makeBlameLog()

	entries, err := Blame(patches, log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	index := make(map[string]BlameEntry)
	for _, e := range entries {
		index[e.ConfigPath] = e
	}

	if e, ok := index["server/port"]; !ok {
		t.Fatal("expected blame entry for server/port")
	} else if e.PatchID != "patch-002" {
		t.Errorf("expected patch-002 to be last, got %s", e.PatchID)
	}
}

func TestBlame_UnrelatedPatchIgnored(t *testing.T) {
	patches := makeBlamePatches()[:1] // only patch-001
	log := makeBlameLog()

	entries, err := Blame(patches, log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if e.PatchID == "patch-002" || e.PatchID == "patch-003" {
			t.Errorf("unexpected patch in blame: %s", e.PatchID)
		}
	}
}

func TestBlame_EmptyLog(t *testing.T) {
	entries, err := Blame(makeBlamePatches(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty log, got %d", len(entries))
	}
}

func TestFormatBlame_ContainsHeader(t *testing.T) {
	out := FormatBlame([]BlameEntry{
		{ConfigPath: "db/host", PatchID: "patch-003", Op: "add", AppliedAt: "2024-01-03T10:00:00Z"},
	})
	if !strings.Contains(out, "PATH") {
		t.Error("expected header row in blame output")
	}
	if !strings.Contains(out, "patch-003") {
		t.Error("expected patch ID in blame output")
	}
}

func TestFormatBlame_Empty(t *testing.T) {
	out := FormatBlame(nil)
	if !strings.Contains(out, "no blame") {
		t.Error("expected empty-state message")
	}
}
