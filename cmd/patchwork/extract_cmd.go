package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var configFile string
	var prefix string
	var outputFmt string

	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract a sub-tree from a config file by path prefix",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("reading config: %w", err)
			}

			cfg, err := parseConfig(raw, configFile)
			if err != nil {
				return fmt.Errorf("parsing config: %w", err)
			}

			result, err := patch.Extract(cfg, prefix)
			if err != nil {
				return err
			}

			switch outputFmt {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(result.Data)
			case "summary":
				fmt.Print(patch.FormatExtractResult(result))
			default:
				return fmt.Errorf("unknown output format %q (use json or summary)", outputFmt)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to config file (required)")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Path prefix to extract (required)")
	cmd.Flags().StringVarP(&outputFmt, "output", "o", "json", "Output format: json|summary")
	_ = cmd.MarkFlagRequired("config")
	_ = cmd.MarkFlagRequired("prefix")

	rootCmd.AddCommand(cmd)
}
