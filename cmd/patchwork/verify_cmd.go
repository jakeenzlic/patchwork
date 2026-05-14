package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	verifyCmd.Flags().StringP("config", "c", "", "Path to the target config file (required)")
	verifyCmd.Flags().StringP("patch", "p", "", "Path to the patch file to verify (required)")
	verifyCmd.Flags().Bool("json", false, "Output result as JSON")
	_ = verifyCmd.MarkFlagRequired("config")
	_ = verifyCmd.MarkFlagRequired("patch")
	rootCmd.AddCommand(verifyCmd)
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that a patch is applicable to the given config",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		patchPath, _ := cmd.Flags().GetString("patch")
		jsonOut, _ := cmd.Flags().GetBool("json")

		p, err := patch.LoadFromFile(patchPath)
		if err != nil {
			return fmt.Errorf("load patch: %w", err)
		}

		raw, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("read config: %w", err)
		}

		cfg, err := parseConfig(raw, configPath)
		if err != nil {
			return fmt.Errorf("parse config: %w", err)
		}

		result, err := patch.Verify(p, cfg)
		if err != nil {
			return fmt.Errorf("verify: %w", err)
		}

		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Patch : %s\n", result.PatchID)
		fmt.Printf("SHA256: %s\n", result.Checksum)
		if result.Match {
			fmt.Println("Status: ✓ patch is applicable")
		} else {
			fmt.Printf("Status: ✗ %s\n", result.Message)
			os.Exit(1)
		}
		return nil
	},
}

// parseConfig is a thin shim so the cmd package can reuse internal parsing
// without importing unexported symbols directly.
func parseConfig(data []byte, filename string) (map[string]any, error) {
	return patch.ParseConfigBytes(data, filename)
}
