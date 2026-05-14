package patch

import (
	"testing"
)

func baseConfig() map[string]any {
	return map[string]any{
		"app": map[string]any{
			"name":    "patchwork",
			"version": float64(2),
			"debug":   false,
		},
		"database": map[string]any{
			"host": "localhost",
		},
	}
}

func TestValidateSchema_NoViolations(t *testing.T) {
	schema := Schema{Rules: []SchemaRule{
		{Path: "app.name", Type: "string", Required: true},
		{Path: "app.version", Type: "number", Required: true},
		{Path: "app.debug", Type: "bool"},
	}}
	violations := ValidateSchema(baseConfig(), schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got: %v", violations)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	schema := Schema{Rules: []SchemaRule{
		{Path: "app.missing", Required: true},
	}}
	violations := ValidateSchema(baseConfig(), schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_WrongType(t *testing.T) {
	schema := Schema{Rules: []SchemaRule{
		{Path: "app.name", Type: "number"},
	}}
	violations := ValidateSchema(baseConfig(), schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_OptionalMissing(t *testing.T) {
	schema := Schema{Rules: []SchemaRule{
		{Path: "app.optional", Type: "string", Required: false},
	}}
	violations := ValidateSchema(baseConfig(), schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for optional missing path, got: %v", violations)
	}
}

func TestValidateSchema_NestedObject(t *testing.T) {
	schema := Schema{Rules: []SchemaRule{
		{Path: "database", Type: "object", Required: true},
		{Path: "database.host", Type: "string", Required: true},
	}}
	violations := ValidateSchema(baseConfig(), schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got: %v", violations)
	}
}

func TestValidateSchema_MultipleViolations(t *testing.T) {
	schema := Schema{Rules: []SchemaRule{
		{Path: "missing.a", Required: true},
		{Path: "missing.b", Required: true},
		{Path: "app.name", Type: "bool"},
	}}
	violations := ValidateSchema(baseConfig(), schema)
	if len(violations) != 3 {
		t.Fatalf("expected 3 violations, got %d: %v", len(violations), violations)
	}
}
