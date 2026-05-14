// Package main is the entry point for the patchwork CLI tool.
// It provides subcommands for applying, rolling back, planning,
// diffing, linting, and exporting config patches.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "patchwork",
		Short: "A lightweight diff-based config migration tool",
		Long: `patchwork tracks and applies incremental changes to JSON/YAML
configs across deployments using a patch file format.`,
	}

	root.AddCommand(
		applyCmd(),
		rollbackCmd(),
		planCmd(),
		diffCmd(),
		lintCmd(),
		exportCmd(),
	)

	return root
}

func applyCmd() *cobra.Command {
	var patchDir, configFile, historyFile string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply pending patches to a config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return patch.Run(patchDir, configFile, historyFile)
		},
	}

	cmd.Flags().StringVarP(&patchDir, "patches", "p", "patches", "Directory containing patch files")
	cmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Target config file to patch")
	cmd.Flags().StringVar(&historyFile, "history", ".patchwork_history.json", "Path to the history file")

	return cmd
}

func rollbackCmd() *cobra.Command {
	var configFile, historyFile string
	var steps int

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Roll back the last N applied patches",
		RunE: func(cmd *cobra.Command, args []string) error {
			return patch.Rollback(configFile, historyFile, steps)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Target config file")
	cmd.Flags().StringVar(&historyFile, "history", ".patchwork_history.json", "Path to the history file")
	cmd.Flags().IntVarP(&steps, "steps", "n", 1, "Number of patches to roll back")

	return cmd
}

func planCmd() *cobra.Command {
	var patchDir, historyFile string

	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Preview which patches would be applied",
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := patch.LoadDir(patchDir)
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}
			history, err := patch.LoadHistory(historyFile)
			if err != nil {
				return fmt.Errorf("loading history: %w", err)
			}
			entries, err := patch.Plan(patches, history)
			if err != nil {
				return fmt.Errorf("planning: %w", err)
			}
			fmt.Println(patch.FormatPlan(entries))
			return nil
		},
	}

	cmd.Flags().StringVarP(&patchDir, "patches", "p", "patches", "Directory containing patch files")
	cmd.Flags().StringVar(&historyFile, "history", ".patchwork_history.json", "Path to the history file")

	return cmd
}

func diffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [before] [after]",
		Short: "Generate a patch by diffing two config files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			before, err := patch.LoadFromFile(args[0])
			if err != nil {
				return fmt.Errorf("loading before config: %w", err)
			}
			after, err := patch.LoadFromFile(args[1])
			if err != nil {
				return fmt.Errorf("loading after config: %w", err)
			}
			ops := patch.Diff(before, after)
			for _, op := range ops {
				fmt.Printf("%s\t%s\t%v\n", op.Op, op.Path, op.Value)
			}
			return nil
		},
	}

	return cmd
}

func lintCmd() *cobra.Command {
	var patchDir string

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint patch files for common issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := patch.LoadDir(patchDir)
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}
			hasWarnings := false
			for _, p := range patches {
				warnings := patch.Lint(p)
				for _, w := range warnings {
					fmt.Printf("WARN [%s]: %s\n", p.Version, w)
					hasWarnings = true
				}
			}
			if !hasWarnings {
				fmt.Println("No lint warnings found.")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&patchDir, "patches", "p", "patches", "Directory containing patch files")

	return cmd
}

func exportCmd() *cobra.Command {
	var patchDir, outputFile, format string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export all patches to a single JSON or YAML file",
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := patch.LoadDir(patchDir)
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}
			return patch.Export(patches, outputFile, format)
		},
	}

	cmd.Flags().StringVarP(&patchDir, "patches", "p", "patches", "Directory containing patch files")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "patches_export.json", "Output file path")
	cmd.Flags().StringVarP(&format, "format", "f", "", "Output format: json or yaml (inferred from output filename if omitted)")

	return cmd
}
