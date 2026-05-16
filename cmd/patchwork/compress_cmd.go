package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	bundleCmd := &cobra.Command{
		Use:   "bundle [patches-dir] [output]",
		Short: "Compress a directory of patches into a .patch.gz bundle",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, dest := args[0], args[1]

			patches, err := patch.LoadDir(dir)
			if err != nil {
				return fmt.Errorf("load patches: %w", err)
			}

			if len(patches) == 0 {
				return fmt.Errorf("bundle: no patches found in %q", dir)
			}

			if err := patch.Bundle(patches, dest); err != nil {
				return fmt.Errorf("bundle: %w", err)
			}

			fmt.Fprintf(os.Stdout, "bundled %d patch(es) → %s\n", len(patches), dest)
			return nil
		},
	}

	unbundleCmd := &cobra.Command{
		Use:   "unbundle [bundle.patch.gz] [config]",
		Short: "Decompress and apply a .patch.gz bundle to a config file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, configPath := args[0], args[1]

			if !patch.IsBundle(src) {
				return fmt.Errorf("unbundle: %q does not look like a .patch.gz bundle", src)
			}

			patches, err := patch.Unbundle(src)
			if err != nil {
				return fmt.Errorf("unbundle: %w", err)
			}

			for _, p := range patches {
				if err := patch.Run(p, configPath, configPath, false, false); err != nil {
					return fmt.Errorf("apply %s: %w", p.ID, err)
				}
			}

			fmt.Fprintf(os.Stdout, "applied %d patch(es) from bundle\n", len(patches))
			return nil
		},
	}

	rootCmd.AddCommand(bundleCmd)
	rootCmd.AddCommand(unbundleCmd)
}
