package patch

import (
	"testing"
)

func TestRenderTemplate_NoPlaceholders(t *testing.T) {
	out, err := RenderTemplate("database.host", TemplateContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "database.host" {
		t.Errorf("expected unchanged string, got %q", out)
	}
}

func TestRenderTemplate_SubstitutesVar(t *testing.T) {
	out, err := RenderTemplate("{{ENV}}.database.host", TemplateContext{"ENV": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "prod.database.host" {
		t.Errorf("expected 'prod.database.host', got %q", out)
	}
}

func TestRenderTemplate_MissingVar(t *testing.T) {
	_, err := RenderTemplate("{{MISSING}}.key", TemplateContext{})
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestRenderTemplate_MultipleVars(t *testing.T) {
	out, err := RenderTemplate("{{A}}/{{B}}", TemplateContext{"A": "foo", "B": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "foo/bar" {
		t.Errorf("expected 'foo/bar', got %q", out)
	}
}

func TestRenderPatches_SubstitutesPathAndValue(t *testing.T) {
	p := Patch{
		Version: "1",
		Ops: []Op{
			{Op: "replace", Path: "{{ENV}}.timeout", Value: "{{TIMEOUT}}"},
		},
	}
	ctx := TemplateContext{"ENV": "prod", "TIMEOUT": "30s"}
	out, err := RenderPatches(p, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Ops[0].Path != "prod.timeout" {
		t.Errorf("expected path 'prod.timeout', got %q", out.Ops[0].Path)
	}
	if out.Ops[0].Value != "30s" {
		t.Errorf("expected value '30s', got %v", out.Ops[0].Value)
	}
}

func TestRenderPatches_MissingVarReturnsError(t *testing.T) {
	p := Patch{
		Version: "1",
		Ops: []Op{
			{Op: "add", Path: "{{UNDEFINED}}.key", Value: "x"},
		},
	}
	_, err := RenderPatches(p, TemplateContext{})
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestExtractTemplateVars(t *testing.T) {
	p := Patch{
		Version: "1",
		Ops: []Op{
			{Op: "replace", Path: "{{ENV}}.host", Value: "{{HOST}}"},
			{Op: "add", Path: "{{ENV}}.port", Value: 8080},
		},
	}
	vars := ExtractTemplateVars(p)
	set := map[string]bool{}
	for _, v := range vars {
		set[v] = true
	}
	if !set["ENV"] {
		t.Error("expected ENV in vars")
	}
	if !set["HOST"] {
		t.Error("expected HOST in vars")
	}
	if len(vars) != 2 {
		t.Errorf("expected 2 unique vars, got %d", len(vars))
	}
}
