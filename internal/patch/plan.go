package patch

import (
	"fmt"
	"strings"
)

// PlanEntry describes a single patch operation that would be applied.
type PlanEntry struct {
	Index   int    `json:"index"`
	Op      string `json:"op"`
	Path    string `json:"path"`
	Value   any    `json:"value,omitempty"`
	Applied bool   `json:"applied"`
}

// Plan returns a dry-run summary of what Apply would do for each operation
// in the patch, without actually modifying the target config.
// It also reports which operations have already been recorded in history.
func Plan(p *Patch, target map[string]any, history *History) ([]PlanEntry, error) {
	if err := Validate(p); err != nil {
		return nil, fmt.Errorf("plan: invalid patch: %w", err)
	}

	entries := make([]PlanEntry, 0, len(p.Ops))

	alreadyApplied := false
	if history != nil {
		alreadyApplied = history.Applied(p.Version)
	}

	for i, op := range p.Ops {
		entry := PlanEntry{
			Index:   i,
			Op:      op.Op,
			Path:    op.Path,
			Value:   op.Value,
			Applied: alreadyApplied,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// FormatPlan returns a human-readable summary of plan entries.
func FormatPlan(entries []PlanEntry) string {
	if len(entries) == 0 {
		return "no operations to apply\n"
	}

	var sb strings.Builder
	for _, e := range entries {
		status := "pending"
		if e.Applied {
			status = "already applied"
		}
		if e.Value != nil {
			fmt.Fprintf(&sb, "  [%d] %-8s %-30s = %v  (%s)\n", e.Index, e.Op, e.Path, e.Value, status)
		} else {
			fmt.Fprintf(&sb, "  [%d] %-8s %-30s  (%s)\n", e.Index, e.Op, e.Path, status)
		}
	}
	return sb.String()
}

// PendingEntries returns only the plan entries that have not yet been applied.
func PendingEntries(entries []PlanEntry) []PlanEntry {
	pending := make([]PlanEntry, 0, len(entries))
	for _, e := range entries {
		if !e.Applied {
			pending = append(pending, e)
		}
	}
	return pending
}
