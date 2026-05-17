package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	checkpointCmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "Manage named checkpoints of applied patch state",
	}

	createCmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a checkpoint with the current applied patch state",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			cp, err := patch.CreateCheckpoint(dir, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("checkpoint %q created at %s with %d applied patch(es)\n",
				cp.Name, cp.CreatedAt.Format(time.RFC3339), len(cp.Applied))
			return nil
		},
	}
	createCmd.Flags().String("dir", ".", "config directory")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved checkpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			cps, err := patch.LoadCheckpoints(dir)
			if err != nil {
				return err
			}
			if len(cps) == 0 {
				fmt.Println("no checkpoints found")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tCREATED\tAPPLIED")
			for _, c := range cps {
				fmt.Fprintf(w, "%s\t%s\t%d\n",
					c.Name, c.CreatedAt.Format(time.RFC3339), len(c.Applied))
			}
			return w.Flush()
		},
	}
	listCmd.Flags().String("dir", ".", "config directory")

	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show patches recorded in a checkpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			cp, ok, err := patch.FindCheckpoint(dir, args[0])
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("checkpoint %q not found", args[0])
			}
			fmt.Printf("checkpoint: %s\ncreated:    %s\napplied (%d):\n",
				cp.Name, cp.CreatedAt.Format(time.RFC3339), len(cp.Applied))
			for _, id := range cp.Applied {
				fmt.Printf("  - %s\n", id)
			}
			return nil
		},
	}
	showCmd.Flags().String("dir", ".", "config directory")

	checkpointCmd.AddCommand(createCmd, listCmd, showCmd)
	rootCmd.AddCommand(checkpointCmd)
}
