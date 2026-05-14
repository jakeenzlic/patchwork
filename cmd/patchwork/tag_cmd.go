package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Tag-related commands",
	}

	listTagsCmd := &cobra.Command{
		Use:   "list <patch-dir>",
		Short: "List all tags and the patches that carry them",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := patch.LoadDir(args[0])
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}
			idx, err := patch.BuildTagIndex(patches)
			if err != nil {
				return err
			}
			fmt.Print(patch.FormatTagIndex(idx))
			return nil
		},
	}

	filterTagCmd := &cobra.Command{
		Use:   "filter <patch-dir> <tag>",
		Short: "Print IDs of patches that carry a given tag",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			patches, err := patch.LoadDir(args[0])
			if err != nil {
				return fmt.Errorf("loading patches: %w", err)
			}
			if err := patch.ValidateTag(args[1]); err != nil {
				return err
			}
			matched := patch.FilterByTag(patches, args[1])
			if len(matched) == 0 {
				fmt.Fprintf(os.Stderr, "no patches found with tag %q\n", args[1])
				return nil
			}
			for _, p := range matched {
				fmt.Println(p.String())
			}
			return nil
		},
	}

	tagCmd.AddCommand(listTagsCmd, filterTagCmd)
	rootCmd.AddCommand(tagCmd)
}
