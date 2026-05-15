package patch

import (
	"strings"
	"testing"
)

func makeReorderPatch(id string, priority int) Patch {
	return Patch{
		ID:       id,
		Version:  "1.0",
		Priority: priority,
		Ops: []Op{
			{Op: "add", Path: "key", Value: "v"},
		},
	}
}

func TestReorder_AlreadySorted(t *testing.T) {
	patches := []Patch{
		makeReorderPatch("a", 1),
		makeReorderPatch("b", 2),
		makeReorderPatch("c", 3),
	}
	res, err := Reorder(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Moved) != 0 {
		t.Errorf("expected no moves, got %v", res.Moved)
	}
}

func TestReorder_SortsByPriority(t *testing.T) {
	patches := []Patch{
		makeReorderPatch("c", 3),
		makeReorderPatch("a", 1),
		makeReorderPatch("b", 2),
	}
	res, err := Reorder(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	order := []string{"a", "b", "c"}
	for i, p := range res.Reordered {
		if p.ID != order[i] {
			t.Errorf("position %d: want %s, got %s", i, order[i], p.ID)
		}
	}
}

func TestReorder_ZeroPriorityLast(t *testing.T) {
	patches := []Patch{
		makeReorderPatch("no-pri", 0),
		makeReorderPatch("first", 1),
	}
	res, err := Reorder(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Reordered[0].ID != "first" {
		t.Errorf("expected 'first' at position 0, got %s", res.Reordered[0].ID)
	}
	if res.Reordered[1].ID != "no-pri" {
		t.Errorf("expected 'no-pri' at position 1, got %s", res.Reordered[1].ID)
	}
}

func TestReorder_StableByID(t *testing.T) {
	patches := []Patch{
		makeReorderPatch("z", 1),
		makeReorderPatch("a", 1),
	}
	res, err := Reorder(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Reordered[0].ID != "a" {
		t.Errorf("expected 'a' first for equal priority, got %s", res.Reordered[0].ID)
	}
}

func TestReorder_Empty(t *testing.T) {
	res, err := Reorder([]Patch{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Reordered) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestFormatReorderResult_NoMoves(t *testing.T) {
	res := ReorderResult{Moved: []string{}}
	out := FormatReorderResult(res)
	if !strings.Contains(out, "no reordering") {
		t.Errorf("expected 'no reordering' message, got: %s", out)
	}
}

func TestFormatReorderResult_WithMoves(t *testing.T) {
	res := ReorderResult{
		Reordered: []Patch{makeReorderPatch("a", 1), makeReorderPatch("b", 2)},
		Moved:     []string{"b"},
	}
	out := FormatReorderResult(res)
	if !strings.Contains(out, "1 patch(es) reordered") {
		t.Errorf("expected reorder count, got: %s", out)
	}
	if !strings.Contains(out, "a") {
		t.Errorf("expected patch id 'a' in output")
	}
}
