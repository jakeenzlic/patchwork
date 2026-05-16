package patch

import "fmt"

// ChainResult holds the outcome of applying a chain of patches sequentially.
type ChainResult struct {
	Applied []string
	Skipped []string
	Failed  string
	Config  map[string]any
}

// Chain applies a sequence of patches in order to the provided config,
// stopping on the first failure. Already-applied patches (per history) are
// skipped. Returns a ChainResult describing what happened.
func Chain(patches []Patch, cfg map[string]any, historyDir string) (ChainResult, error) {
	hist, err := LoadHistory(historyDir)
	if err != nil {
		return ChainResult{}, fmt.Errorf("chain: load history: %w", err)
	}

	current := deepCopy(cfg)
	result := ChainResult{Config: current}

	for _, p := range patches {
		if hist.Applied(p.ID) {
			result.Skipped = append(result.Skipped, p.ID)
			continue
		}
		if err := Validate(p); err != nil {
			result.Failed = p.ID
			return result, fmt.Errorf("chain: validate %q: %w", p.ID, err)
		}
		updated, err := Apply(p, current)
		if err != nil {
			result.Failed = p.ID
			return result, fmt.Errorf("chain: apply %q: %w", p.ID, err)
		}
		current = updated
		result.Applied = append(result.Applied, p.ID)
	}

	result.Config = current
	return result, nil
}

// FormatChainResult returns a human-readable summary of a ChainResult.
func FormatChainResult(r ChainResult) string {
	out := ""
	for _, id := range r.Applied {
		out += fmt.Sprintf("  applied  %s\n", id)
	}
	for _, id := range r.Skipped {
		out += fmt.Sprintf("  skipped  %s\n", id)
	}
	if r.Failed != "" {
		out += fmt.Sprintf("  failed   %s\n", r.Failed)
	}
	if out == "" {
		return "  (nothing to do)\n"
	}
	return out
}
