package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StashEntry represents a single stashed config state.
type StashEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message"`
	Config    map[string]interface{} `json:"config"`
}

// StashPath returns the path to the stash file for a given config.
func StashPath(configPath string) string {
	dir := filepath.Dir(configPath)
	base := filepath.Base(configPath)
	return filepath.Join(dir, ".patchwork", base+".stash.json")
}

// LoadStash reads all stash entries for a config file.
func LoadStash(configPath string) ([]StashEntry, error) {
	p := StashPath(configPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return []StashEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read stash: %w", err)
	}
	var entries []StashEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse stash: %w", err)
	}
	return entries, nil
}

// SaveStash writes stash entries to disk.
func SaveStash(configPath string, entries []StashEntry) error {
	p := StashPath(configPath)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("create stash dir: %w", err)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal stash: %w", err)
	}
	return os.WriteFile(p, data, 0o644)
}

// Stash saves the current config state with an optional message.
func Stash(configPath, message string) (StashEntry, error) {
	cfg, err := parseConfig(configPath)
	if err != nil {
		return StashEntry{}, fmt.Errorf("load config: %w", err)
	}
	entries, err := LoadStash(configPath)
	if err != nil {
		return StashEntry{}, err
	}
	entry := StashEntry{
		ID:        fmt.Sprintf("stash@{%d}", len(entries)),
		Timestamp: time.Now().UTC(),
		Message:   message,
		Config:    cfg,
	}
	entries = append(entries, entry)
	if err := SaveStash(configPath, entries); err != nil {
		return StashEntry{}, err
	}
	return entry, nil
}

// StashPop restores the most recent stash entry and removes it.
func StashPop(configPath string) (StashEntry, error) {
	entries, err := LoadStash(configPath)
	if err != nil {
		return StashEntry{}, err
	}
	if len(entries) == 0 {
		return StashEntry{}, fmt.Errorf("no stash entries found")
	}
	last := entries[len(entries)-1]
	if err := Export(last.Config, configPath, ""); err != nil {
		return StashEntry{}, fmt.Errorf("restore config: %w", err)
	}
	if err := SaveStash(configPath, entries[:len(entries)-1]); err != nil {
		return StashEntry{}, err
	}
	return last, nil
}

// FormatStash returns a human-readable summary of stash entries.
func FormatStash(entries []StashEntry) string {
	if len(entries) == 0 {
		return "No stash entries.\n"
	}
	out := ""
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		out += fmt.Sprintf("%s  %s  %s\n", e.ID, e.Timestamp.Format(time.RFC3339), e.Message)
	}
	return out
}
