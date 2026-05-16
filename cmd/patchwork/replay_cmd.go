package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var configPath string
	var patchDir string
	var stopAt string

	replayCmd := &cobra.Command{
		Use:   "replay",
		Short: "Replay patches against a config from scratch",
		Long:  "Re-applies all patches in order against the base config, optionally stopping at a given patch ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := patch.LoadSnapshot("", configPath)
			if err != nil {
				// Fall back to parsing the config file directly.
				base, err = parseConfig(configPath)
				if err != nil {
					return fmt.Errorf("loading config: %w", err)
				}
			}

			patches, err := patch.LoadDir(patchDir)
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}

			result, err := patch.Replay(base, patches, stopAt)
			if err != nil {
				return fmt.Errorf("replay failed: %w", err)
			}

			fmt.Print(patch.FormatReplayResult(result))

			if stopAt != "" {
				fmt.Fprintf(os.Stderr, "Stopped at patch %q\n", stopAt)
			}

			return nil
		},
	}

	replayCmd.Flags().StringVarP(&configPath, "config", "c", "config.json", "Path to the base config file")
	replayCmd.Flags().StringVarP(&patchDir, "patches", "p", "patches", "Directory containing patch files")
	replayCmd.Flags().StringVar(&stopAt, "stop-at", "", "Stop replay at this patch ID (inclusive)")

	rootCmd.AddCommand(replayCmd)
}
