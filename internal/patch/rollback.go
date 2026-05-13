package patch

import (
	"fmt"
	"os"
	"path/filepath"
)

// RollbackOptions configures rollback behaviour.
type RollbackOptions struct {
	// PatchDir is the directory containing patch files.
	PatchDir string
	// HistoryFile is the path to the history JSON file.
	HistoryFile string
	// Target is the config file to roll back.
	Target string
	// Steps is the number of applied patches to undo (default 1).
	Steps int
}

// Rollback reverses the last N applied patches against target.
// It loads history, identifies the patches to undo in reverse order,
// generates inverse operations, and writes the result back to target.
func Rollback(opts RollbackOptions) error {
	if opts.Steps <= 0 {
		opts.Steps = 1
	}

	h, err := LoadHistory(opts.HistoryFile)
	if err != nil {
		return fmt.Errorf("rollback: load history: %w", err)
	}

	applied := h.Applied
	if len(applied) == 0 {
		return fmt.Errorf("rollback: no applied patches in history")
	}
	if opts.Steps > len(applied) {
		return fmt.Errorf("rollback: requested %d steps but only %d patches applied", opts.Steps, len(applied))
	}

	// Patches to undo, most-recent first.
	toUndo := applied[len(applied)-opts.Steps:]

	// Load current config.
	raw, err := os.ReadFile(opts.Target)
	if err != nil {
		return fmt.Errorf("rollback: read target: %w", err)
	}
	current, err := parseConfig(raw, opts.Target)
	if err != nil {
		return fmt.Errorf("rollback: parse target: %w", err)
	}

	// Apply inverse patches in reverse order.
	for i := len(toUndo) - 1; i >= 0; i-- {
		patchFile := filepath.Join(opts.PatchDir, toUndo[i])
		p, err := LoadFromFile(patchFile)
		if err != nil {
			return fmt.Errorf("rollback: load patch %s: %w", toUndo[i], err)
		}
		inverse := invertPatch(p)
		current, err = Apply(current, inverse)
		if err != nil {
			return fmt.Errorf("rollback: apply inverse of %s: %w", toUndo[i], err)
		}
	}

	// Write result back.
	out, err := Export(current, inferFormat(opts.Target))
	if err != nil {
		return fmt.Errorf("rollback: export: %w", err)
	}
	if err := os.WriteFile(opts.Target, out, 0644); err != nil {
		return fmt.Errorf("rollback: write target: %w", err)
	}

	// Trim history.
	h.Applied = applied[:len(applied)-opts.Steps]
	return h.Save(opts.HistoryFile)
}
