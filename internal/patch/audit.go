package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEntry records a single patch application event.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	PatchFile string    `json:"patch_file"`
	Version   string    `json:"version"`
	Ops       int       `json:"ops"`
	User      string    `json:"user,omitempty"`
	Status    string    `json:"status"` // "applied" | "failed" | "rolled_back"
	Message   string    `json:"message,omitempty"`
}

// AuditLog is a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry `json:"entries"`
}

// AuditPath returns the path to the audit log file.
func AuditPath(configPath string) string {
	dir := filepath.Dir(configPath)
	return filepath.Join(dir, ".patchwork", "audit.json")
}

// LoadAuditLog reads the audit log from disk, returning an empty log if not found.
func LoadAuditLog(configPath string) (*AuditLog, error) {
	p := AuditPath(configPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read audit log: %w", err)
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parse audit log: %w", err)
	}
	return &log, nil
}

// RecordAudit appends an entry to the audit log.
func RecordAudit(configPath string, entry AuditEntry) error {
	log, err := LoadAuditLog(configPath)
	if err != nil {
		return err
	}
	log.Entries = append(log.Entries, entry)
	p := AuditPath(configPath)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("create audit dir: %w", err)
	}
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal audit log: %w", err)
	}
	return os.WriteFile(p, data, 0o644)
}

// FormatAuditLog returns a human-readable summary of the audit log.
func FormatAuditLog(log *AuditLog) string {
	if len(log.Entries) == 0 {
		return "No audit entries found.\n"
	}
	out := ""
	for _, e := range log.Entries {
		line := fmt.Sprintf("[%s] %s  version=%s ops=%d status=%s",
			e.Timestamp.Format(time.RFC3339), e.PatchFile, e.Version, e.Ops, e.Status)
		if e.Message != "" {
			line += "  msg=" + e.Message
		}
		out += line + "\n"
	}
	return out
}
