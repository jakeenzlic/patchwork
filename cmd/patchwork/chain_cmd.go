package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var patchDir string
	var historyDir string
	var configFile string

	cmd := &cobra.Command{
		Use:   "chain",
		Short: "Apply a directory of patches sequentially, skipping already-applied ones",
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := patch.LoadDir(patchDir)
			if err != nil {
				return fmt.Errorf("load patches: %w", err)
			}

			raw, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("read config: %w", err)
			}
			var cfg map[string]any
			if err := json.Unmarshal(raw, &cfg); err != nil {
				return fmt.Errorf("parse config: %w", err)
			}

			res, err := patch.Chain(patches, cfg, historyDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "chain failed:\n%s", patch.FormatChainResult(res))
				return err
			}

			fmt.Print(patch.FormatChainResult(res))

			out, err := json.MarshalIndent(res.Config, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal result: %w", err)
			}
			return os.WriteFile(configFile, append(out, '\n'), 0o644)
		},
	}

	cmd.Flags().StringVarP(&patchDir, "patches", "p", "patches", "Directory containing patch files")
	cmd.Flags().StringVarP(&historyDir, "history", "H", ".patchwork", "Directory for history state")
	cmd.Flags().StringVarP(&configFile, "config", "c", "config.json", "Config file to update in place")

	rootCmd.AddCommand(cmd)
}
