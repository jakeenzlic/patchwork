package patch

import (
	"fmt"
	"strings"
)

// BlameEntry records which patch last touched a given config path.
type BlameEntry struct {
	ConfigPath string
	PatchID    string
	Op         string
	AppliedAt  string
}

// Blame inspects an audit log and returns, for each config path touched by
// any patch in patches, the most-recent audit entry that modified it.
func Blame(patches []Patch, log []AuditEntry) ([]BlameEntry, error) {
	// Build a set of patch IDs we care about.
	wanted := make(map[string]struct{}, len(patches))
	for _, p := range patches {
		wanted[p.ID] = struct{}{}
	}

	// Map config-path -> latest entry (log is assumed chronological).
	latest := make(map[string]BlameEntry)

	for _, entry := range log {
		if _, ok := wanted[entry.PatchID]; !ok {
			continue
		}
		// Find the patch to inspect its ops.
		for _, p := range patches {
			if p.ID != entry.PatchID {
				continue
			}
			for _, op := range p.Ops {
				latest[op.Path] = BlameEntry{
					ConfigPath: op.Path,
					PatchID:    p.ID,
					Op:         op.Op,
					AppliedAt:  entry.AppliedAt,
				}
			}
		}
	}

	result := make([]BlameEntry, 0, len(latest))
	for _, e := range latest {
		result = append(result, e)
	}
	return result, nil
}

// FormatBlame returns a human-readable blame report.
func FormatBlame(entries []BlameEntry) string {
	if len(entries) == 0 {
		return "no blame information available\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-40s %-20s %-10s %s\n", "PATH", "PATCH", "OP", "APPLIED AT"))
	sb.WriteString(strings.Repeat("-", 90) + "\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-40s %-20s %-10s %s\n",
			e.ConfigPath, e.PatchID, e.Op, e.AppliedAt))
	}
	return sb.String()
}
