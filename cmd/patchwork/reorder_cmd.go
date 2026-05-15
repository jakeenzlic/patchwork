package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	reorderCmd := &cobra.Command{
		Use:   "reorder <patch-dir>",
		Short: "Display patches sorted by priority",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]

			patches, err := patch.LoadDir(dir)
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}

			res, err := patch.Reorder(patches)
			if err != nil {
				return fmt.Errorf("reordering patches: %w", err)
			}

			apply, _ := cmd.Flags().GetBool("apply")
			if apply {
				if len(res.Moved) == 0 {
					fmt.Fprintln(os.Stdout, "patches already in correct order")
					return nil
				}
				fmt.Fprintf(os.Stdout, "reordered %d patch(es) by priority\n", len(res.Moved))
				for _, id := range res.Moved {
					fmt.Fprintf(os.Stdout, "  moved: %s\n", id)
				}
				return nil
			}

			fmt.Print(patch.FormatReorderResult(res))
			return nil
		},
	}

	reorderCmd.Flags().Bool("apply", false, "show which patches would be moved")
	rootCmd.AddCommand(reorderCmd)
}
