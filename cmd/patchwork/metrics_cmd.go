package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var jsonOut bool
	var configPath string

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show patch application metrics from the audit log",
		RunE: func(cmd *cobra.Command, args []string) error {
			auditPath := patch.AuditPath(configPath)
			entries, err := patch.LoadAuditLog(auditPath)
			if err != nil {
				return fmt.Errorf("loading audit log: %w", err)
			}

			// Convert audit entries to metrics entries.
			var metrics []patch.MetricsEntry
			for _, a := range entries {
				status := "success"
				errMsg := ""
				if a.Action == "rollback" {
					status = "rollback"
				}
				metrics = append(metrics, patch.MetricsEntry{
					PatchID:   a.PatchID,
					AppliedAt: a.Timestamp,
					Duration:  0, // audit log does not store duration
					OpCount:   len(a.Ops),
					Status:    status,
					Error:     errMsg,
				})
			}

			if jsonOut {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(metrics)
			}

			fmt.Print(patch.FormatMetrics(metrics))
			return nil
		},
	}

	metricsCmd.Flags().BoolVar(&jsonOut, "json", false, "Output metrics as JSON")
	metricsCmd.Flags().StringVarP(&configPath, "config", "c", "config.json", "Path to config file")
	rootCmd.AddCommand(metricsCmd)
}
