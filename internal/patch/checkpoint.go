package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint records a named point-in-time snapshot of the applied patch history.
type Checkpoint struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Applied   []string  `json:"applied"`
}

// CheckpointPath returns the path to the checkpoint index file.
func CheckpointPath(dir string) string {
	return filepath.Join(dir, ".patchwork", "checkpoints.json")
}

// LoadCheckpoints reads all saved checkpoints from disk.
func LoadCheckpoints(dir string) ([]Checkpoint, error) {
	p := CheckpointPath(dir)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return []Checkpoint{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read checkpoints: %w", err)
	}
	var cps []Checkpoint
	if err := json.Unmarshal(data, &cps); err != nil {
		return nil, fmt.Errorf("parse checkpoints: %w", err)
	}
	return cps, nil
}

// SaveCheckpoints writes the checkpoint list to disk.
func SaveCheckpoints(dir string, cps []Checkpoint) error {
	p := CheckpointPath(dir)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// CreateCheckpoint saves a named checkpoint containing the currently applied patch IDs.
func CreateCheckpoint(dir, name string) (Checkpoint, error) {
	if name == "" {
		return Checkpoint{}, fmt.Errorf("checkpoint name must not be empty")
	}
	cps, err := LoadCheckpoints(dir)
	if err != nil {
		return Checkpoint{}, err
	}
	for _, c := range cps {
		if c.Name == name {
			return Checkpoint{}, fmt.Errorf("checkpoint %q already exists", name)
		}
	}
	h, err := LoadHistory(dir)
	if err != nil {
		return Checkpoint{}, err
	}
	applied := make([]string, 0, len(h.Applied))
	for id := range h.Applied {
		applied = append(applied, id)
	}
	cp := Checkpoint{Name: name, CreatedAt: time.Now().UTC(), Applied: applied}
	cps = append(cps, cp)
	if err := SaveCheckpoints(dir, cps); err != nil {
		return Checkpoint{}, err
	}
	return cp, nil
}

// FindCheckpoint looks up a checkpoint by name.
func FindCheckpoint(dir, name string) (Checkpoint, bool, error) {
	cps, err := LoadCheckpoints(dir)
	if err != nil {
		return Checkpoint{}, false, err
	}
	for _, c := range cps {
		if c.Name == name {
			return c, true, nil
		}
	}
	return Checkpoint{}, false, nil
}
