package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit <config>",
		Short: "Display the audit log for a config file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := args[0]

			log, err := patch.LoadAuditLog(configPath)
			if err != nil {
				return fmt.Errorf("load audit log: %w", err)
			}

			jsonFlag, _ := cmd.Flags().GetBool("json")
			if jsonFlag {
				return patch.Export(log, os.Stdout, "json")
			}

			fmt.Print(patch.FormatAuditLog(log))
			return nil
		},
	}

	auditCmd.Flags().Bool("json", false, "Output audit log as JSON")
	rootCmd.AddCommand(auditCmd)

	auditClearCmd := &cobra.Command{
		Use:   "audit-clear <config>",
		Short: "Clear the audit log for a config file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := args[0]
			p := patch.AuditPath(configPath)
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("clear audit log: %w", err)
			}
			fmt.Println("Audit log cleared.")
			return nil
		},
	}
	rootCmd.AddCommand(auditClearCmd)
}
