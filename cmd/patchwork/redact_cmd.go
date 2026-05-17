package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var paths []string
	var mask string
	var output string
	var format string

	cmd := &cobra.Command{
		Use:   "redact <config>",
		Short: "Mask sensitive values in a config file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := args[0]

			cfg, err := parseConfig(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			var rules []patch.RedactRule
			for _, p := range paths {
				rules = append(rules, patch.RedactRule{
					Path:     strings.TrimSpace(p),
					MaskWith: mask,
				})
			}

			result := patch.Redact(cfg, rules)
			fmt.Print(patch.FormatRedactResult(result))

			if output != "" {
				if err := patch.Export(result.Config, output, format); err != nil {
					return fmt.Errorf("export: %w", err)
				}
				fmt.Fprintf(os.Stdout, "redacted config written to %s\n", output)
			}

			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&paths, "path", "p", nil, "config path to redact (repeatable)")
	cmd.Flags().StringVar(&mask, "mask", "***", "mask string to substitute")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write redacted config to file")
	cmd.Flags().StringVar(&format, "format", "", "output format: json or yaml (inferred from --output if omitted)")
	_ = cmd.MarkFlagRequired("path")

	rootCmd.AddCommand(cmd)
}
