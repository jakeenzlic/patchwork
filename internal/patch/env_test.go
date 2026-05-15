package patch

import (
	"testing"
)

func makeResolver(vars map[string]string) *EnvResolver {
	return NewEnvResolver(func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	})
}

func TestResolvePatches_NoPlaceholders(t *testing.T) {
	r := makeResolver(nil)
	input := []PatchEntry{{Op: "add", Path: "a/b", Value: "hello"}}
	out, err := r.ResolvePatches(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "hello" {
		t.Errorf("expected 'hello', got %v", out[0].Value)
	}
}

func TestResolvePatches_SubstitutesVar(t *testing.T) {
	r := makeResolver(map[string]string{"ENV": "production"})
	input := []PatchEntry{{Op: "replace", Path: "env", Value: "${ENV}"}}
	out, err := r.ResolvePatches(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "production" {
		t.Errorf("expected 'production', got %v", out[0].Value)
	}
}

func TestResolvePatches_MultipleVarsInValue(t *testing.T) {
	r := makeResolver(map[string]string{"HOST": "localhost", "PORT": "5432"})
	input := []PatchEntry{{Op: "add", Path: "dsn", Value: "${HOST}:${PORT}"}}
	out, err := r.ResolvePatches(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost:5432" {
		t.Errorf("expected 'localhost:5432', got %v", out[0].Value)
	}
}

func TestResolvePatches_MissingVarReturnsError(t *testing.T) {
	r := makeResolver(nil)
	input := []PatchEntry{{Op: "add", Path: "x", Value: "${MISSING}"}}
	_, err := r.ResolvePatches(input)
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestResolvePatches_NonStringValueUnchanged(t *testing.T) {
	r := makeResolver(map[string]string{"X": "1"})
	input := []PatchEntry{{Op: "add", Path: "count", Value: 42}}
	out, err := r.ResolvePatches(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != 42 {
		t.Errorf("expected 42, got %v", out[0].Value)
	}
}

func TestResolvePatches_PreservesOtherFields(t *testing.T) {
	r := makeResolver(map[string]string{"TAG": "v2"})
	input := []PatchEntry{{Op: "replace", Path: "version", Value: "${TAG}"}}
	out, err := r.ResolvePatches(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Op != "replace" || out[0].Path != "version" {
		t.Errorf("op or path mutated unexpectedly")
	}
}
