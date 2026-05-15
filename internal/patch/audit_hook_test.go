package patch

import (
	"path/filepath"
	"testing"
)

func TestRecordSuccess_WritesEntry(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")
	opts := AuditOptions{Enabled: true, ConfigPath: cfg}

	recordSuccess(opts, "001.json", "1", 2)

	log, err := LoadAuditLog(cfg)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	if log.Entries[0].Status != "applied" {
		t.Fatalf("expected applied, got %s", log.Entries[0].Status)
	}
}

func TestRecordFailure_WritesEntry(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")
	opts := AuditOptions{Enabled: true, ConfigPath: cfg}

	recordFailure(opts, "002.json", "2", 1, "apply error")

	log, _ := LoadAuditLog(cfg)
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	e := log.Entries[0]
	if e.Status != "failed" || e.Message != "apply error" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestRecordRollback_WritesEntry(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")
	opts := AuditOptions{Enabled: true, ConfigPath: cfg}

	recordRollback(opts, "001.json", "1", 2)

	log, _ := LoadAuditLog(cfg)
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if log.Entries[0].Status != "rolled_back" {
		t.Fatalf("expected rolled_back")
	}
}

func TestRecordSuccess_Disabled(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")
	opts := AuditOptions{Enabled: false, ConfigPath: cfg}

	recordSuccess(opts, "001.json", "1", 2)

	log, _ := LoadAuditLog(cfg)
	if len(log.Entries) != 0 {
		t.Fatalf("expected no entries when disabled")
	}
}
