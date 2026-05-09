package patch

import (
	"testing"
)

func TestDiff_AddedKey(t *testing.T) {
	src := map[string]interface{}{"version": "1.0"}
	dst := map[string]interface{}{"version": "1.0", "debug": true}

	ops, err := Diff(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}
	if ops[0].Op != "add" || ops[0].Path != "debug" {
		t.Errorf("unexpected op: %+v", ops[0])
	}
}

func TestDiff_DeletedKey(t *testing.T) {
	src := map[string]interface{}{"version": "1.0", "legacy": "yes"}
	dst := map[string]interface{}{"version": "1.0"}

	ops, err := Diff(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}
	if ops[0].Op != "delete" || ops[0].Path != "legacy" {
		t.Errorf("unexpected op: %+v", ops[0])
	}
}

func TestDiff_ReplacedValue(t *testing.T) {
	src := map[string]interface{}{"timeout": 30}
	dst := map[string]interface{}{"timeout": 60}

	ops, err := Diff(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}
	if ops[0].Op != "replace" || ops[0].Path != "timeout" {
		t.Errorf("unexpected op: %+v", ops[0])
	}
}

func TestDiff_NestedChange(t *testing.T) {
	src := map[string]interface{}{"db": map[string]interface{}{"host": "localhost", "port": 5432}}
	dst := map[string]interface{}{"db": map[string]interface{}{"host": "prod-db", "port": 5432}}

	ops, err := Diff(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}
	if ops[0].Op != "replace" || ops[0].Path != "db.host" {
		t.Errorf("unexpected op: %+v", ops[0])
	}
}

func TestDiff_NoChanges(t *testing.T) {
	src := map[string]interface{}{"key": "value"}
	dst := map[string]interface{}{"key": "value"}

	ops, err := Diff(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ops) != 0 {
		t.Errorf("expected no ops, got %d: %+v", len(ops), ops)
	}
}
