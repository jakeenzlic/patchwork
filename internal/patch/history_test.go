package patch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadHistory_NotExist(t *testing.T) {
	h, err := LoadHistory("/tmp/patchwork_no_such_file.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestHistory_RecordAndApplied(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := &History{}
	entry := HistoryEntry{
		AppliedAt: time.Now().UTC(),
		PatchFile: "001_add_feature.json",
		Version:   "1.0.0",
		Success:   true,
	}

	if err := h.Record(path, entry); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	loaded, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", loaded.Entries[0].Version)
	}
}

func TestHistory_Applied(t *testing.T) {
	h := &History{
		Entries: []HistoryEntry{
			{Version: "1.0.0", Success: true},
			{Version: "1.1.0", Success: false},
		},
	}

	if !h.Applied("1.0.0") {
		t.Error("expected 1.0.0 to be applied")
	}
	if h.Applied("1.1.0") {
		t.Error("expected 1.1.0 not to be applied (failed)")
	}
	if h.Applied("2.0.0") {
		t.Error("expected 2.0.0 not to be applied")
	}
}

func TestHistoryPath(t *testing.T) {
	got := HistoryPath("/etc/myapp")
	want := "/etc/myapp/.patchwork_history.json"
	if got != want {
		t.Errorf("HistoryPath = %q, want %q", got, want)
	}
	_ = os.Getenv // keep os import used
}
