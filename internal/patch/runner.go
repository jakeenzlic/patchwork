package patch

import (
	"fmt"
	"time"
)

// RunOptions configures the patch runner.
type RunOptions struct {
	// PatchDir is the directory containing patch files.
	PatchDir string
	// HistoryFile is where applied-patch history is persisted.
	HistoryFile string
	// DryRun skips writing changes but logs what would happen.
	DryRun bool
}

// RunResult summarises the outcome of a run.
type RunResult struct {
	Applied  []string
	Skipped  []string
	Failed   []string
}

// Run loads all patches from opts.PatchDir, skips already-applied ones
// according to the history file, and applies the rest in order.
func Run(cfg map[string]any, opts RunOptions) (RunResult, error) {
	history, err := LoadHistory(opts.HistoryFile)
	if err != nil {
		return RunResult{}, fmt.Errorf("runner: load history: %w", err)
	}

	patches, err := LoadDir(opts.PatchDir)
	if err != nil {
		return RunResult{}, fmt.Errorf("runner: load patches: %w", err)
	}

	var result RunResult

	for _, p := range patches {
		if history.Applied(p.Version) {
			result.Skipped = append(result.Skipped, p.Version)
			continue
		}

		if err := Validate(p); err != nil {
			entry := HistoryEntry{
				AppliedAt: time.Now().UTC(),
				PatchFile: opts.PatchDir,
				Version:   p.Version,
				Success:   false,
				Error:     err.Error(),
			}
			if !opts.DryRun {
				_ = history.Record(opts.HistoryFile, entry)
			}
			result.Failed = append(result.Failed, p.Version)
			continue
		}

		if !opts.DryRun {
			if err := Apply(cfg, p); err != nil {
				entry := HistoryEntry{
					AppliedAt: time.Now().UTC(),
					PatchFile: opts.PatchDir,
					Version:   p.Version,
					Success:   false,
					Error:     err.Error(),
				}
				_ = history.Record(opts.HistoryFile, entry)
				result.Failed = append(result.Failed, p.Version)
				continue
			}
			entry := HistoryEntry{
				AppliedAt: time.Now().UTC(),
				PatchFile: opts.PatchDir,
				Version:   p.Version,
				Success:   true,
			}
			if err := history.Record(opts.HistoryFile, entry); err != nil {
				return result, fmt.Errorf("runner: record history: %w", err)
			}
		}

		result.Applied = append(result.Applied, p.Version)
	}

	return result, nil
}
