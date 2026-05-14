package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot <config-file>",
	Short: "Show the latest snapshot for a config file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := args[0]

		snap, err := patch.LoadSnapshot(configPath)
		if err != nil {
			return fmt.Errorf("failed to load snapshot: %w", err)
		}
		if snap == nil {
			fmt.Fprintf(os.Stderr, "no snapshot found for %s\n", configPath)
			return nil
		}

		fmt.Printf("Snapshot for: %s\n", configPath)
		fmt.Printf("  Patch ID : %s\n", snap.PatchID)
		fmt.Printf("  Timestamp: %s\n", snap.Timestamp.Format("2006-01-02 15:04:05 UTC"))

		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			data, err := json.MarshalIndent(snap.Config, "  ", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}
			fmt.Printf("  Config:\n  %s\n", string(data))
		}
		return nil
	},
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save <config-file>",
	Short: "Manually save a snapshot of the current config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := args[0]

		cfg, err := patch.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		patchID, _ := cmd.Flags().GetString("patch-id")
		if patchID == "" {
			patchID = "manual"
		}

		if err := patch.SaveSnapshot(configPath, patchID, cfg); err != nil {
			return fmt.Errorf("failed to save snapshot: %w", err)
		}

		fmt.Printf("snapshot saved for %s (patch: %s)\n", configPath, patchID)
		return nil
	},
}

func init() {
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.Flags().BoolP("verbose", "v", false, "print full config contents")
	snapshotSaveCmd.Flags().String("patch-id", "", "patch ID to tag the snapshot with")
}
