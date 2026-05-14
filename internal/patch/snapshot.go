package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of a config file.
type Snapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	PatchID   string                 `json:"patch_id"`
	Config    map[string]interface{} `json:"config"`
}

// SnapshotPath returns the path to the snapshot file for a given config.
func SnapshotPath(configPath string) string {
	dir := filepath.Dir(configPath)
	base := filepath.Base(configPath)
	return filepath.Join(dir, ".patchwork", "snapshots", base+".snap.json")
}

// SaveSnapshot writes a snapshot of the current config state to disk.
func SaveSnapshot(configPath, patchID string, config map[string]interface{}) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		PatchID:   patchID,
		Config:    config,
	}

	path := SnapshotPath(configPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// LoadSnapshot reads the latest snapshot for a given config path.
func LoadSnapshot(configPath string) (*Snapshot, error) {
	path := SnapshotPath(configPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}
