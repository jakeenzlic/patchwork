package patch

import (
	"fmt"
	"strings"
)

// TrimResult holds the outcome of a trim operation.
type TrimResult struct {
	Removed []string
	Kept    []string
}

// Trim removes patch operations whose paths are no longer present in the
// provided config, returning a filtered patch and a result summary.
func Trim(p Patch, config map[string]any) (Patch, TrimResult, error) {
	if err := Validate(p); err != nil {
		return Patch{}, TrimResult{}, fmt.Errorf("trim: invalid patch: %w", err)
	}

	var kept, removed []Op
	var result TrimResult

	for _, op := range p.Ops {
		if op.Op == "add" {
			// add ops are always kept — they introduce new paths
			kept = append(kept, op)
			result.Kept = append(result.Kept, op.Path)
			continue
		}

		if pathExists(config, op.Path) {
			kept = append(kept, op)
			result.Kept = append(result.Kept, op.Path)
		} else {
			removed = append(removed, op)
			result.Removed = append(result.Removed, op.Path)
		}
	}

	trimmed := Patch{
		ID:      p.ID,
		Version: p.Version,
		Ops:     kept,
	}

	return trimmed, result, nil
}

// FormatTrimResult returns a human-readable summary of a TrimResult.
func FormatTrimResult(r TrimResult) string {
	var sb strings.Builder

	if len(r.Removed) == 0 {
		sb.WriteString("trim: no stale operations found\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("trim: removed %d stale operation(s):\n", len(r.Removed)))
	for _, p := range r.Removed {
		sb.WriteString(fmt.Sprintf("  - %s\n", p))
	}

	if len(r.Kept) > 0 {
		sb.WriteString(fmt.Sprintf("trim: kept %d operation(s)\n", len(r.Kept)))
	}

	return sb.String()
}
