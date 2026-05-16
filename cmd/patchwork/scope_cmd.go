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
	var prefix string

	scopeCmd := &cobra.Command{
		Use:   "scope",
		Short: "Filter and preview patches restricted to a config subtree",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := patch.NewScope(prefix)
			if err != nil {
				return fmt.Errorf("invalid scope: %w", err)
			}

			patches, err := patch.LoadDir(patchDir)
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}

			cfgBytes, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("reading config: %w", err)
			}
			cfg, err := parseConfig(configPath, cfgBytes)
			if err != nil {
				return fmt.Errorf("parsing config: %w", err)
			}

			scoped := s.Filter(patches)

			totalOps := 0
			for _, p := range patches {
				totalOps += len(p.Operations)
			}
			keptOps := 0
			for _, p := range scoped {
				keptOps += len(p.Operations)
			}

			fmt.Println(patch.FormatScope(prefix, totalOps, keptOps))
			fmt.Printf("patches in scope: %d\n", len(scoped))

			for _, p := range scoped {
				entries, err := patch.Plan([]patch.Patch{p}, cfg, "")
				if err != nil {
					fmt.Fprintf(os.Stderr, "warn: plan error for %s: %v\n", p.ID, err)
					continue
				}
				fmt.Print(patch.FormatPlan(entries))
			}
			return nil
		},
	}

	scopeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file (required)")
	scopeCmd.Flags().StringVarP(&patchDir, "patches", "p", "patches", "Directory containing patch files")
	scopeCmd.Flags().StringVarP(&prefix, "prefix", "x", "", "Scope prefix (e.g. 'database') (required)")
	_ = scopeCmd.MarkFlagRequired("config")
	_ = scopeCmd.MarkFlagRequired("prefix")

	rootCmd.AddCommand(scopeCmd)
}
