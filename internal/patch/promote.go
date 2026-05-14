package patch

import (
	"fmt"
	"os"
	"path/filepath"
)

// PromoteResult holds the outcome of a promotion between environments.
type PromoteResult struct {
	SourceEnv string
	TargetEnv string
	Copied    []string
	Skipped   []string
}

// Promote copies unapplied patches from one environment's history context to
// another by reading the source patch directory and skipping any patch IDs
// already recorded in the target history.
func Promote(patchDir, sourceEnv, targetEnv, configDir string) (*PromoteResult, error) {
	sourceDir := filepath.Join(patchDir, sourceEnv)
	targetDir := filepath.Join(patchDir, targetEnv)

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("source environment directory not found: %s", sourceDir)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create target directory: %w", err)
	}

	targetHistoryPath := HistoryPath(filepath.Join(configDir, targetEnv))
	targetHistory, err := LoadHistory(targetHistoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load target history: %w", err)
	}

	applied := make(map[string]bool, len(targetHistory.Applied))
	for _, id := range targetHistory.Applied {
		applied[id] = true
	}

	patches, err := LoadDir(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load source patches: %w", err)
	}

	result := &PromoteResult{
		SourceEnv: sourceEnv,
		TargetEnv: targetEnv,
	}

	for _, p := range patches {
		if applied[p.ID] {
			result.Skipped = append(result.Skipped, p.ID)
			continue
		}

		destPath := filepath.Join(targetDir, p.ID+".json")
		if err := Export(p, destPath, "json"); err != nil {
			return nil, fmt.Errorf("failed to export patch %s: %w", p.ID, err)
		}
		result.Copied = append(result.Copied, p.ID)
	}

	return result, nil
}

// FormatPromoteResult returns a human-readable summary of a promotion.
func FormatPromoteResult(r *PromoteResult) string {
	out := fmt.Sprintf("Promote: %s → %s\n", r.SourceEnv, r.TargetEnv)
	out += fmt.Sprintf("  Copied:  %d patch(es)\n", len(r.Copied))
	for _, id := range r.Copied {
		out += fmt.Sprintf("    + %s\n", id)
	}
	out += fmt.Sprintf("  Skipped: %d patch(es) (already applied)\n", len(r.Skipped))
	for _, id := range r.Skipped {
		out += fmt.Sprintf("    ~ %s\n", id)
	}
	return out
}
