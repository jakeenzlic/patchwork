package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var message string

	stashCmd := &cobra.Command{
		Use:   "stash <config>",
		Short: "Save the current config state to the stash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, err := patch.Stash(args[0], message)
			if err != nil {
				return fmt.Errorf("stash: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Saved %s — %s\n", entry.ID, entry.Message)
			return nil
		},
	}
	stashCmd.Flags().StringVarP(&message, "message", "m", "", "Description of the stash entry")

	stashPopCmd := &cobra.Command{
		Use:   "stash-pop <config>",
		Short: "Restore the most recent stash entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, err := patch.StashPop(args[0])
			if err != nil {
				return fmt.Errorf("stash-pop: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Restored %s — %s\n", entry.ID, entry.Message)
			return nil
		},
	}

	stashListCmd := &cobra.Command{
		Use:   "stash-list <config>",
		Short: "List all stash entries for a config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := patch.LoadStash(args[0])
			if err != nil {
				return fmt.Errorf("stash-list: %w", err)
			}
			fmt.Fprint(cmd.OutOrStdout(), patch.FormatStash(entries))
			return nil
		},
	}

	stashDropCmd := &cobra.Command{
		Use:   "stash-drop <config>",
		Short: "Remove all stash entries for a config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := patch.SaveStash(args[0], []patch.StashEntry{}); err != nil {
				return fmt.Errorf("stash-drop: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Stash cleared.")
			return nil
		},
	}

	for _, c := range []*cobra.Command{stashCmd, stashPopCmd, stashListCmd, stashDropCmd} {
		rootCmd.AddCommand(c)
	}
	_ = os.Stderr
}
