package patch

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// WatchEvent represents a change detected in the patch directory.
type WatchEvent struct {
	Path    string
	ModTime time.Time
	Err     error
}

// WatchOptions configures the behaviour of Watch.
type WatchOptions struct {
	// PatchDir is the directory to watch for new or modified patch files.
	PatchDir string
	// ConfigPath is the target config file that patches are applied to.
	ConfigPath string
	// HistoryDir is the directory used to persist apply history.
	HistoryDir string
	// Interval controls how often the directory is polled.
	// Defaults to 5 seconds when zero.
	Interval time.Duration
	// OnApply is called after each successful patch application.
	OnApply func(event WatchEvent)
	// OnError is called when an error occurs during a poll cycle.
	OnError func(err error)
}

// Watch polls PatchDir at the configured Interval and automatically applies
// any patch files that have not yet been recorded in the history. It blocks
// until ctx is cancelled.
func Watch(ctx context.Context, opts WatchOptions) error {
	if opts.PatchDir == "" {
		return fmt.Errorf("watch: PatchDir must not be empty")
	}
	if opts.ConfigPath == "" {
		return fmt.Errorf("watch: ConfigPath must not be empty")
	}

	interval := opts.Interval
	if interval <= 0 {
		interval = 5 * time.Second
	}

	onApply := opts.OnApply
	if onApply == nil {
		onApply = func(WatchEvent) {}
	}
	onError := opts.OnError
	if onError == nil {
		onError = func(err error) { log.Println("watch error:", err) }
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := pollAndApply(opts, onApply, onError); err != nil {
				onError(err)
			}
		}
	}
}

// pollAndApply loads all patch files from PatchDir, filters out those already
// in the history, and applies the remainder in order.
func pollAndApply(opts WatchOptions, onApply func(WatchEvent), onError func(error)) error {
	patches, err := LoadDir(opts.PatchDir)
	if err != nil {
		return fmt.Errorf("watch: loading patches: %w", err)
	}

	histPath := HistoryPath(opts.HistoryDir, opts.ConfigPath)
	hist, err := LoadHistory(histPath)
	if err != nil {
		return fmt.Errorf("watch: loading history: %w", err)
	}

	for _, p := range patches {
		if hist.Applied(p.ID) {
			continue
		}

		// Read and parse the current config.
		data, err := os.ReadFile(opts.ConfigPath)
		if err != nil {
			onError(fmt.Errorf("watch: reading config %s: %w", opts.ConfigPath, err))
			continue
		}

		cfg, err := parseConfig(data, opts.ConfigPath)
		if err != nil {
			onError(fmt.Errorf("watch: parsing config: %w", err))
			continue
		}

		result, err := Apply(p, cfg)
		if err != nil {
			onError(fmt.Errorf("watch: applying patch %s: %w", p.ID, err))
			continue
		}

		ext := filepath.Ext(opts.ConfigPath)
		if err := Export(result, opts.ConfigPath, inferFormat(ext)); err != nil {
			onError(fmt.Errorf("watch: exporting config: %w", err))
			continue
		}

		hist.Record(p.ID)
		if err := hist.Save(histPath); err != nil {
			onError(fmt.Errorf("watch: saving history: %w", err))
			continue
		}

		onApply(WatchEvent{
			Path:    opts.ConfigPath,
			ModTime: time.Now(),
		})
	}

	return nil
}
