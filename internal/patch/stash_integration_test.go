package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStash_RoundtripPreservesValues(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	original := `{"host": "localhost", "port": 8080}`
	if err := os.WriteFile(cfgPath, []byte(original), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Stash(cfgPath, "original state")
	if err != nil {
		t.Fatalf("Stash: %v", err)
	}

	// Mutate config on disk
	modified := `{"host": "prod.example.com", "port": 443}`
	if err := os.WriteFile(cfgPath, []byte(modified), 0o644); err != nil {
		t.Fatalf("write modified config: %v", err)
	}

	_, err = StashPop(cfgPath)
	if err != nil {
		t.Fatalf("StashPop: %v", err)
	}

	restored, err := parseConfig(cfgPath)
	if err != nil {
		t.Fatalf("parse restored: %v", err)
	}
	if restored["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", restored["host"])
	}
}

func TestStash_MultiplePopOrder(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	for i, content := range []string{
		`{"step": 1}`,
		`{"step": 2}`,
		`{"step": 3}`,
	} {
		if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
			t.Fatalf("write step %d: %v", i, err)
		}
		if _, err := Stash(cfgPath, content); err != nil {
			t.Fatalf("Stash step %d: %v", i, err)
		}
	}

	entry, err := StashPop(cfgPath)
	if err != nil {
		t.Fatalf("StashPop: %v", err)
	}
	if entry.ID != "stash@{2}" {
		t.Errorf("expected last entry popped, got %q", entry.ID)
	}

	remaining, _ := LoadStash(cfgPath)
	if len(remaining) != 2 {
		t.Errorf("expected 2 remaining entries, got %d", len(remaining))
	}
}

func TestStash_StashPathUniqueness(t *testing.T) {
	p1 := StashPath("/app/prod/config.json")
	p2 := StashPath("/app/staging/config.json")
	if p1 == p2 {
		t.Error("stash paths should differ for different directories")
	}
}
