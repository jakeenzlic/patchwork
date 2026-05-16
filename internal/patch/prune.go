package patch

import "fmt"

// PruneResult holds the outcome of a prune operation.
type PruneResult struct {
	Removed []string
	Kept    []string
}

// Prune removes patches from the provided list that have already been applied
// according to the history at historyDir, and whose IDs match the given
// candidateIDs filter. If candidateIDs is empty all applied patches are pruned.
func Prune(patches []Patch, historyDir string, candidateIDs []string) (PruneResult, error) {
	hist, err := LoadHistory(historyDir)
	if err != nil {
		return PruneResult{}, fmt.Errorf("prune: load history: %w", err)
	}

	applied := make(map[string]bool, len(hist.Entries))
	for _, e := range hist.Entries {
		applied[e.PatchID] = true
	}

	filter := make(map[string]bool, len(candidateIDs))
	for _, id := range candidateIDs {
		filter[id] = true
	}

	var result PruneResult
	for _, p := range patches {
		if !applied[p.ID] {
			result.Kept = append(result.Kept, p.ID)
			continue
		}
		if len(filter) > 0 && !filter[p.ID] {
			result.Kept = append(result.Kept, p.ID)
			continue
		}
		result.Removed = append(result.Removed, p.ID)
	}
	return result, nil
}

// FormatPruneResult returns a human-readable summary of a PruneResult.
func FormatPruneResult(r PruneResult) string {
	out := fmt.Sprintf("pruned %d patch(es)\n", len(r.Removed))
	for _, id := range r.Removed {
		out += fmt.Sprintf("  - removed: %s\n", id)
	}
	for _, id := range r.Kept {
		out += fmt.Sprintf("  - kept:    %s\n", id)
	}
	return out
}
