package patch

import (
	"testing"
)

func makePrunePatch(id string) Patch {
	return Patch{
		ID:      id,
		Version: "1",
		Ops: []Op{
			{Op: "add", Path: "key", Value: id},
		},
	}
}

func TestPrune_RemovesApplied(t *testing.T) {
	dir := t.TempDir()
	patches := []Patch{makePrunePatch("p1"), makePrunePatch("p2"), makePrunePatch("p3")}

	hist, _ := LoadHistory(dir)
	hist.Entries = append(hist.Entries, HistoryEntry{PatchID: "p1"}, HistoryEntry{PatchID: "p3"})
	_ = saveHistory(dir, hist)

	result, err := Prune(patches, dir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(result.Removed))
	}
	if len(result.Kept) != 1 || result.Kept[0] != "p2" {
		t.Errorf("expected kept=[p2], got %v", result.Kept)
	}
}

func TestPrune_FilterLimitsCandidates(t *testing.T) {
	dir := t.TempDir()
	patches := []Patch{makePrunePatch("p1"), makePrunePatch("p2")}

	hist, _ := LoadHistory(dir)
	hist.Entries = append(hist.Entries, HistoryEntry{PatchID: "p1"}, HistoryEntry{PatchID: "p2"})
	_ = saveHistory(dir, hist)

	// Only prune p1 even though p2 is also applied
	result, err := Prune(patches, dir, []string{"p1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 1 || result.Removed[0] != "p1" {
		t.Errorf("expected removed=[p1], got %v", result.Removed)
	}
	if len(result.Kept) != 1 || result.Kept[0] != "p2" {
		t.Errorf("expected kept=[p2], got %v", result.Kept)
	}
}

func TestPrune_NoneApplied(t *testing.T) {
	dir := t.TempDir()
	patches := []Patch{makePrunePatch("p1"), makePrunePatch("p2")}

	result, err := Prune(patches, dir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
	if len(result.Kept) != 2 {
		t.Errorf("expected 2 kept, got %d", len(result.Kept))
	}
}

func TestFormatPruneResult_Output(t *testing.T) {
	r := PruneResult{
		Removed: []string{"p1"},
		Kept:    []string{"p2"},
	}
	out := FormatPruneResult(r)
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
	for _, want := range []string{"p1", "p2", "pruned"} {
		if !contains(out, want) {
			t.Errorf("expected output to contain %q", want)
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
