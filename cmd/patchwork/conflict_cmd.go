package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var patchDir string

	conflictCmd := &cobra.Command{
		Use:   "conflict",
		Short: "Detect conflicting operations across patch files",
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := patch.LoadDir(patchDir)
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}

			conflicts := patch.DetectConflicts(patches)
			fmt.Print(patch.FormatConflicts(conflicts))

			if len(conflicts) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	conflictCmd.Flags().StringVarP(&patchDir, "dir", "d", "patches", "Directory containing patch files")
	rootCmd.AddCommand(conflictCmd)
}
