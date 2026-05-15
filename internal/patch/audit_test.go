package patch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAuditPath(t *testing.T) {
	p := AuditPath("/configs/app.json")
	want := "/configs/.patchwork/audit.json"
	if p != want {
		t.Fatalf("got %q, want %q", p, want)
	}
}

func TestLoadAuditLog_NotExist(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")
	log, err := LoadAuditLog(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) != 0 {
		t.Fatalf("expected empty log, got %d entries", len(log.Entries))
	}
}

func TestRecordAndLoadAudit(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")

	entry := AuditEntry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		PatchFile: "001_init.json",
		Version:   "1",
		Ops:       3,
		Status:    "applied",
	}
	if err := RecordAudit(cfg, entry); err != nil {
		t.Fatalf("record: %v", err)
	}

	log, err := LoadAuditLog(cfg)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	got := log.Entries[0]
	if got.PatchFile != entry.PatchFile || got.Status != entry.Status {
		t.Fatalf("entry mismatch: %+v", got)
	}
}

func TestRecordAudit_Accumulates(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")

	for i := 0; i < 3; i++ {
		if err := RecordAudit(cfg, AuditEntry{
			Timestamp: time.Now().UTC(),
			PatchFile: "patch.json",
			Version:   "1",
			Ops:       1,
			Status:    "applied",
		}); err != nil {
			t.Fatalf("record %d: %v", i, err)
		}
	}

	log, err := LoadAuditLog(cfg)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(log.Entries))
	}
}

func TestFormatAuditLog_Empty(t *testing.T) {
	out := FormatAuditLog(&AuditLog{})
	if out != "No audit entries found.\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestFormatAuditLog_WithEntries(t *testing.T) {
	log := &AuditLog{
		Entries: []AuditEntry{
			{Timestamp: time.Now().UTC(), PatchFile: "001.json", Version: "1", Ops: 2, Status: "applied"},
			{Timestamp: time.Now().UTC(), PatchFile: "002.json", Version: "2", Ops: 1, Status: "failed", Message: "bad op"},
		},
	}
	out := FormatAuditLog(log)
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	_ = os.Stdout // suppress unused import
}
