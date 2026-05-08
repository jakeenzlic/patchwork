package patch

import (
	"testing"
	"time"
)

func basePatch(ops ...Operation) *Patch {
	return &Patch{
		Version:     "v1",
		Description: "test patch",
		CreatedAt:   time.Now(),
		Ops:         ops,
	}
}

func TestValidate_Valid(t *testing.T) {
	p := basePatch(Operation{Op: OpAdd, Path: "/db/host", Value: "localhost"})
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingVersion(t *testing.T) {
	p := basePatch()
	p.Version = ""
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestValidate_UnknownOp(t *testing.T) {
	p := basePatch(Operation{Op: "patch", Path: "/a", Value: 1})
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestApply_Add(t *testing.T) {
	cfg := map[string]interface{}{"app": map[string]interface{}{"port": 8080}}
	p := basePatch(Operation{Op: OpAdd, Path: "/app/debug", Value: true})
	result, err := Apply(cfg, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := result["app"].(map[string]interface{})
	if app["debug"] != true {
		t.Errorf("expected debug=true, got %v", app["debug"])
	}
}

func TestApply_Replace(t *testing.T) {
	cfg := map[string]interface{}{"app": map[string]interface{}{"port": 8080}}
	p := basePatch(Operation{Op: OpReplace, Path: "/app/port", Value: 9090})
	result, err := Apply(cfg, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := result["app"].(map[string]interface{})
	if app["port"] != 9090 {
		t.Errorf("expected port=9090, got %v", app["port"])
	}
}

func TestApply_Remove(t *testing.T) {
	cfg := map[string]interface{}{"app": map[string]interface{}{"port": 8080, "debug": true}}
	p := basePatch(Operation{Op: OpRemove, Path: "/app/debug"})
	result, err := Apply(cfg, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := result["app"].(map[string]interface{})
	if _, exists := app["debug"]; exists {
		t.Error("expected debug key to be removed")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	cfg := map[string]interface{}{"key": "original"}
	p := basePatch(Operation{Op: OpReplace, Path: "/key", Value: "changed"})
	_, err := Apply(cfg, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["key"] != "original" {
		t.Error("original config was mutated")
	}
}
