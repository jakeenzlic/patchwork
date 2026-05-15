package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin [patch-file]",
		Short: "Pin a patch to its current checksum in the lock file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgDir, _ := cmd.Flags().GetString("dir")
			p, err := patch.LoadFromFile(args[0])
			if err != nil {
				return fmt.Errorf("loading patch: %w", err)
			}
			if err := patch.Pin(cfgDir, *p); err != nil {
				return fmt.Errorf("pin: %w", err)
			}
			fmt.Fprintf(os.Stdout, "pinned patch %q (version %s)\n", p.ID, p.Version)
			return nil
		},
	}
	pinCmd.Flags().String("dir", ".", "config directory containing the lock file")

	checkLockCmd := &cobra.Command{
		Use:   "check-lock [patch-file]",
		Short: "Verify a patch matches its pinned checksum",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgDir, _ := cmd.Flags().GetString("dir")
			p, err := patch.LoadFromFile(args[0])
			if err != nil {
				return fmt.Errorf("loading patch: %w", err)
			}
			if err := patch.CheckLock(cfgDir, *p); err != nil {
				fmt.Fprintf(os.Stderr, "lock check failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "patch %q lock OK\n", p.ID)
			return nil
		},
	}
	checkLockCmd.Flags().String("dir", ".", "config directory containing the lock file")

	showLockCmd := &cobra.Command{
		Use:   "show-lock",
		Short: "Display all pinned patches in the lock file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgDir, _ := cmd.Flags().GetString("dir")
			lf, err := patch.LoadLockFile(cfgDir)
			if err != nil {
				return err
			}
			if len(lf.Entries) == 0 {
				fmt.Println("no pinned patches")
				return nil
			}
			fmt.Printf("%-30s %-12s %s\n", "PATCH ID", "VERSION", "CHECKSUM")
			for _, e := range lf.Entries {
				fmt.Printf("%-30s %-12s %s\n", e.PatchID, e.Version, e.Checksum)
			}
			return nil
		},
	}
	showLockCmd.Flags().String("dir", ".", "config directory containing the lock file")

	rootCmd.AddCommand(pinCmd, checkLockCmd, showLockCmd)
}
