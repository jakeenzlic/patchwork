package patch

import (
	"strings"
	"testing"
)

func baseRedactConfig() map[string]any {
	return map[string]any{
		"app": map[string]any{
			"name": "patchwork",
			"secret": "s3cr3t",
		},
		"db": map[string]any{
			"password": "hunter2",
			"host":     "localhost",
		},
	}
}

func TestRedact_SinglePath(t *testing.T) {
	cfg := baseRedactConfig()
	rules := []RedactRule{{Path: "db/password"}}
	res := Redact(cfg, rules)

	dbMap := res.Config["db"].(map[string]any)
	if dbMap["password"] != "***" {
		t.Errorf("expected *** got %v", dbMap["password"])
	}
	if len(res.Masked) != 1 || res.Masked[0] != "db/password" {
		t.Errorf("unexpected masked list: %v", res.Masked)
	}
}

func TestRedact_CustomMask(t *testing.T) {
	cfg := baseRedactConfig()
	rules := []RedactRule{{Path: "app/secret", MaskWith: "<redacted>"}}
	res := Redact(cfg, rules)

	appMap := res.Config["app"].(map[string]any)
	if appMap["secret"] != "<redacted>" {
		t.Errorf("expected <redacted> got %v", appMap["secret"])
	}
}

func TestRedact_MissingPath(t *testing.T) {
	cfg := baseRedactConfig()
	rules := []RedactRule{{Path: "app/nonexistent"}}
	res := Redact(cfg, rules)

	if len(res.Masked) != 0 {
		t.Errorf("expected no masked paths, got %v", res.Masked)
	}
}

func TestRedact_OriginalUnmutated(t *testing.T) {
	cfg := baseRedactConfig()
	rules := []RedactRule{{Path: "db/password"}}
	Redact(cfg, rules)

	dbMap := cfg["db"].(map[string]any)
	if dbMap["password"] != "hunter2" {
		t.Errorf("original config was mutated")
	}
}

func TestRedact_MultiplePaths(t *testing.T) {
	cfg := baseRedactConfig()
	rules := []RedactRule{
		{Path: "app/secret"},
		{Path: "db/password"},
	}
	res := Redact(cfg, rules)

	if len(res.Masked) != 2 {
		t.Errorf("expected 2 masked paths, got %d", len(res.Masked))
	}
}

func TestFormatRedactResult_NoMatches(t *testing.T) {
	res := RedactResult{}
	out := FormatRedactResult(res)
	if !strings.Contains(out, "no paths matched") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatRedactResult_WithMatches(t *testing.T) {
	res := RedactResult{Masked: []string{"db/password", "app/secret"}}
	out := FormatRedactResult(res)
	if !strings.Contains(out, "masked 2") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "db/password") {
		t.Errorf("expected db/password in output")
	}
}
