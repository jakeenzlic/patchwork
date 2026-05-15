package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockEntry records a pinned patch version for a given environment.
type LockEntry struct {
	PatchID   string    `json:"patch_id"`
	Version   string    `json:"version"`
	AppliedAt time.Time `json:"applied_at"`
	Checksum  string    `json:"checksum"`
}

// LockFile maps patch IDs to their locked entries.
type LockFile struct {
	Entries map[string]LockEntry `json:"entries"`
}

// LockPath returns the path to the lock file for the given config directory.
func LockPath(dir string) string {
	return filepath.Join(dir, ".patchwork.lock")
}

// LoadLockFile reads an existing lock file or returns an empty one.
func LoadLockFile(dir string) (*LockFile, error) {
	path := LockPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &LockFile{Entries: make(map[string]LockEntry)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("lock: read %s: %w", path, err)
	}
	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("lock: parse %s: %w", path, err)
	}
	if lf.Entries == nil {
		lf.Entries = make(map[string]LockEntry)
	}
	return &lf, nil
}

// SaveLockFile persists the lock file to disk.
func SaveLockFile(dir string, lf *LockFile) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("lock: mkdir %s: %w", dir, err)
	}
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("lock: marshal: %w", err)
	}
	return os.WriteFile(LockPath(dir), data, 0o644)
}

// Pin records a patch as locked at its current version and checksum.
func Pin(dir string, p Patch) error {
	lf, err := LoadLockFile(dir)
	if err != nil {
		return err
	}
	sum, err := Checksum(p)
	if err != nil {
		return fmt.Errorf("lock: checksum: %w", err)
	}
	lf.Entries[p.ID] = LockEntry{
		PatchID:   p.ID,
		Version:   p.Version,
		AppliedAt: time.Now().UTC(),
		Checksum:  sum,
	}
	return SaveLockFile(dir, lf)
}

// CheckLock verifies that a patch matches its locked checksum.
// Returns nil if no lock entry exists (not yet pinned).
func CheckLock(dir string, p Patch) error {
	lf, err := LoadLockFile(dir)
	if err != nil {
		return err
	}
	entry, ok := lf.Entries[p.ID]
	if !ok {
		return nil
	}
	sum, err := Checksum(p)
	if err != nil {
		return fmt.Errorf("lock: checksum: %w", err)
	}
	if sum != entry.Checksum {
		return fmt.Errorf("lock: patch %q checksum mismatch (locked %s, got %s)", p.ID, entry.Checksum, sum)
	}
	return nil
}
