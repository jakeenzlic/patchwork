package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var outputPath string

	renameCmd := &cobra.Command{
		Use:   "rename <patch-file> <old-prefix> <new-prefix>",
		Short: "Rewrite op paths in a patch from one prefix to another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			patchFile := args[0]
			oldPrefix := args[1]
			newPrefix := args[2]

			p, err := patch.LoadFromFile(patchFile)
			if err != nil {
				return fmt.Errorf("load patch: %w", err)
			}

			updated, result, err := patch.Rename(p, oldPrefix, newPrefix)
			if err != nil {
				return fmt.Errorf("rename: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), patch.FormatRenameResult(result))

			dest := patchFile
			if outputPath != "" {
				dest = outputPath
			}

			if err := patch.Export(updated, dest); err != nil {
				return fmt.Errorf("export: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "written to %s\n", dest)
			return nil
		},
	}

	renameCmd.Flags().StringVarP(&outputPath, "output", "o", "", "write result to a different file instead of overwriting")

	if err := rootCmd.RegisterFlagCompletionFunc("rename", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	}); err != nil {
		fmt.Fprintln(os.Stderr, "warn: could not register rename completion")
	}

	rootCmd.AddCommand(renameCmd)
}
