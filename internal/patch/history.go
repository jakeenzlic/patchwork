package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a single applied patch event.
type HistoryEntry struct {
	AppliedAt time.Time `json:"applied_at"`
	PatchFile string    `json:"patch_file"`
	Version   string    `json:"version"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// History holds a list of applied patch entries.
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// HistoryPath returns the default history file path relative to a config dir.
func HistoryPath(dir string) string {
	return filepath.Join(dir, ".patchwork_history.json")
}

// LoadHistory reads a history file from disk. Returns an empty History if the
// file does not exist yet.
func LoadHistory(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &History{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read %s: %w", path, err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("history: parse %s: %w", path, err)
	}
	return &h, nil
}

// Record appends a new entry and persists the history file.
func (h *History) Record(path string, entry HistoryEntry) error {
	h.Entries = append(h.Entries, entry)
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write %s: %w", path, err)
	}
	return nil
}

// Applied returns true if a patch version has already been recorded as
// successful in this history.
func (h *History) Applied(version string) bool {
	for _, e := range h.Entries {
		if e.Version == version && e.Success {
			return true
		}
	}
	return false
}
