package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema <config> <schema>",
		Short: "Validate a config file against a schema",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := args[0]
			schemaPath := args[1]

			cfg, err := parseConfig(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			schema, err := patch.LoadSchema(schemaPath)
			if err != nil {
				return fmt.Errorf("load schema: %w", err)
			}

			violations := patch.ValidateSchema(cfg, schema)
			if len(violations) == 0 {
				fmt.Println("✓ Config satisfies schema")
				return nil
			}

			fmt.Fprintf(os.Stderr, "Schema violations (%d):\n", len(violations))
			for _, v := range violations {
				fmt.Fprintf(os.Stderr, "  - %s\n", v)
			}
			os.Exit(1)
			return nil
		},
	}

	rootCmd.AddCommand(schemaCmd)
}
