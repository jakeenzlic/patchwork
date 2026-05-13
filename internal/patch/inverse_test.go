package patch

import (
	"testing"
)

func TestInvertPatch_Add(t *testing.T) {
	p := Patch{
		Version: "1",
		Ops:     []Op{{Op: "add", Path: "feature", Value: true}},
	}
	inv := invertPatch(p)
	if len(inv.Ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(inv.Ops))
	}
	if inv.Ops[0].Op != "delete" {
		t.Errorf("expected delete, got %s", inv.Ops[0].Op)
	}
	if inv.Ops[0].Path != "feature" {
		t.Errorf("unexpected path %s", inv.Ops[0].Path)
	}
}

func TestInvertPatch_Delete(t *testing.T) {
	p := Patch{
		Version: "1",
		Ops:     []Op{{Op: "delete", Path: "old.key", OldValue: "preserved"}},
	}
	inv := invertPatch(p)
	if inv.Ops[0].Op != "add" {
		t.Errorf("expected add, got %s", inv.Ops[0].Op)
	}
	if inv.Ops[0].Value != "preserved" {
		t.Errorf("expected value 'preserved', got %v", inv.Ops[0].Value)
	}
}

func TestInvertPatch_Replace(t *testing.T) {
	p := Patch{
		Version: "1",
		Ops:     []Op{{Op: "replace", Path: "port", Value: float64(9090), OldValue: float64(8080)}},
	}
	inv := invertPatch(p)
	if inv.Ops[0].Op != "replace" {
		t.Errorf("expected replace, got %s", inv.Ops[0].Op)
	}
	if inv.Ops[0].Value != float64(8080) {
		t.Errorf("expected value 8080, got %v", inv.Ops[0].Value)
	}
	if inv.Ops[0].OldValue != float64(9090) {
		t.Errorf("expected old_value 9090, got %v", inv.Ops[0].OldValue)
	}
}

func TestInvertPatch_ReverseOrder(t *testing.T) {
	p := Patch{
		Version: "1",
		Ops: []Op{
			{Op: "add", Path: "a", Value: 1},
			{Op: "add", Path: "b", Value: 2},
		},
	}
	inv := invertPatch(p)
	if inv.Ops[0].Path != "b" || inv.Ops[1].Path != "a" {
		t.Error("expected inverse ops in reverse order")
	}
}
