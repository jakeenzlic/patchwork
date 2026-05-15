package patch

import (
	"fmt"
	"sort"
)

// ReorderResult holds the outcome of a reorder operation.
type ReorderResult struct {
	Original []Patch
	Reordered []Patch
	Moved []string // IDs of patches that changed position
}

// Reorder sorts patches by priority (ascending), then by ID for stability.
// Patches with no priority set (zero value) are placed after those with explicit priority.
func Reorder(patches []Patch) (ReorderResult, error) {
	if len(patches) == 0 {
		return ReorderResult{}, nil
	}

	for _, p := range patches {
		if err := Validate(p); err != nil {
			return ReorderResult{}, fmt.Errorf("invalid patch %q: %w", p.ID, err)
		}
	}

	original := make([]Patch, len(patches))
	copy(original, patches)

	sorted := make([]Patch, len(patches))
	copy(sorted, patches)

	sort.SliceStable(sorted, func(i, j int) bool {
		pi, pj := sorted[i].Priority, sorted[j].Priority
		if pi == 0 && pj != 0 {
			return false
		}
		if pi != 0 && pj == 0 {
			return true
		}
		if pi != pj {
			return pi < pj
		}
		return sorted[i].ID < sorted[j].ID
	})

	moved := []string{}
	for i, p := range sorted {
		if original[i].ID != p.ID {
			moved = append(moved, p.ID)
		}
	}

	return ReorderResult{
		Original:  original,
		Reordered: sorted,
		Moved:     moved,
	}, nil
}

// FormatReorderResult returns a human-readable summary of the reorder.
func FormatReorderResult(r ReorderResult) string {
	if len(r.Moved) == 0 {
		return "no reordering needed\n"
	}
	out := fmt.Sprintf("%d patch(es) reordered:\n", len(r.Moved))
	for i, p := range r.Reordered {
		pri := ""
		if p.Priority != 0 {
			pri = fmt.Sprintf(" (priority %d)", p.Priority)
		}
		out += fmt.Sprintf("  %d. %s%s\n", i+1, p.ID, pri)
	}
	return out
}
