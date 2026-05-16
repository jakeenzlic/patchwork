package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	var patchFile string
	var annotationsFile string

	cmd := &cobra.Command{
		Use:   "annotate",
		Short: "Attach notes to patch operation paths",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := patch.LoadFromFile(patchFile)
			if err != nil {
				return fmt.Errorf("load patch: %w", err)
			}

			data, err := os.ReadFile(annotationsFile)
			if err != nil {
				return fmt.Errorf("read annotations file: %w", err)
			}

			var anns []patch.Annotation
			if err := json.Unmarshal(data, &anns); err != nil {
				return fmt.Errorf("parse annotations: %w", err)
			}

			result, err := patch.Annotate(*p, anns)
			if err != nil {
				return fmt.Errorf("annotate: %w", err)
			}

			fmt.Println(patch.FormatAnnotations(result))
			return nil
		},
	}

	cmd.Flags().StringVarP(&patchFile, "patch", "p", "", "Path to patch file (required)")
	cmd.Flags().StringVarP(&annotationsFile, "annotations", "a", "", "Path to JSON annotations file (required)")
	_ = cmd.MarkFlagRequired("patch")
	_ = cmd.MarkFlagRequired("annotations")

	rootCmd.AddCommand(cmd)
}
