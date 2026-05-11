package patch

import (
	"testing"
)

func baseValidPatch() *Patch {
	return &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "add", Path: "server/port", Value: 8080},
		},
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := Validate(baseValidPatch()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingVersion(t *testing.T) {
	p := baseValidPatch()
	p.Version = ""
	err := Validate(p)
	if err == nil {
		t.Fatal("expected error for missing version")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(ve.Errors))
	}
}

func TestValidate_UnknownOp(t *testing.T) {
	p := baseValidPatch()
	p.Ops[0].Op = "upsert"
	err := Validate(p)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestValidate_MissingPath(t *testing.T) {
	p := baseValidPatch()
	p.Ops[0].Path = ""
	err := Validate(p)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestValidate_MissingValueForAdd(t *testing.T) {
	p := baseValidPatch()
	p.Ops[0].Value = nil
	err := Validate(p)
	if err == nil {
		t.Fatal("expected error when add op has no value")
	}
}

func TestValidate_RemoveNoValueOK(t *testing.T) {
	p := &Patch{
		Version: "1.0.0",
		Ops: []Op{
			{Op: "remove", Path: "server/port", Value: nil},
		},
	}
	if err := Validate(p); err != nil {
		t.Fatalf("remove without value should be valid, got: %v", err)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	p := &Patch{
		Version: "",
		Ops: []Op{
			{Op: "bad", Path: "", Value: nil},
		},
	}
	err := Validate(p)
	if err == nil {
		t.Fatal("expected multiple errors")
	}
	ve := err.(*ValidationError)
	if len(ve.Errors) < 3 {
		t.Fatalf("expected at least 3 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}
