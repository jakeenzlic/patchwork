package patch

import "fmt"

// ReplayResult holds the outcome of replaying a sequence of patches.
type ReplayResult struct {
	Applied []string
	Skipped []string
	Final   map[string]any
}

// Replay re-applies a list of patches in order against a starting config,
// optionally stopping at a target patch ID (inclusive). If stopAt is empty
// all patches are replayed.
func Replay(base map[string]any, patches []Patch, stopAt string) (ReplayResult, error) {
	result := ReplayResult{
		Final: deepCopy(base),
	}

	for _, p := range patches {
		if err := Validate(p); err != nil {
			return result, fmt.Errorf("invalid patch %q: %w", p.ID, err)
		}

		updated, err := Apply(result.Final, p)
		if err != nil {
			result.Skipped = append(result.Skipped, p.ID)
		} else {
			result.Final = updated
			result.Applied = append(result.Applied, p.ID)
		}

		if stopAt != "" && p.ID == stopAt {
			break
		}
	}

	return result, nil
}

// FormatReplayResult returns a human-readable summary of a ReplayResult.
func FormatReplayResult(r ReplayResult) string {
	out := fmt.Sprintf("Replay complete: %d applied, %d skipped\n", len(r.Applied), len(r.Skipped))
	for _, id := range r.Applied {
		out += fmt.Sprintf("  ✔ %s\n", id)
	}
	for _, id := range r.Skipped {
		out += fmt.Sprintf("  ✘ %s (skipped)\n", id)
	}
	return out
}
