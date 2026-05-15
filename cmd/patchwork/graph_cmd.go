package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"patchwork/internal/patch"
)

func init() {
	graphCmd := &cobra.Command{
		Use:   "graph <patch-dir>",
		Short: "Show patch dependency graph and application order",
		Args:  cobra.ExactArgs(1),
		RunE:  runGraph,
	}
	graphCmd.Flags().Bool("order-only", false, "Print only the resolved application order")
	rootCmd.AddCommand(graphCmd)
}

func runGraph(cmd *cobra.Command, args []string) error {
	dir := args[0]
	orderOnly, _ := cmd.Flags().GetBool("order-only")

	patches, err := patch.LoadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading patches: %v\n", err)
		return err
	}

	if len(patches) == 0 {
		fmt.Println("No patches found.")
		return nil
	}

	g, err := patch.BuildDependencyGraph(patches)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dependency error: %v\n", err)
		return err
	}

	order, err := patch.TopologicalSort(g)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sort error: %v\n", err)
		return err
	}

	if orderOnly {
		fmt.Println("Application order:")
		for i, v := range order {
			fmt.Printf("  %d. %s\n", i+1, v)
		}
		return nil
	}

	fmt.Println("Dependency graph:")
	for _, node := range g.Nodes {
		deps := g.Edges[node]
		if len(deps) == 0 {
			fmt.Printf("  %s (no dependencies)\n", node)
		} else {
			fmt.Printf("  %s -> depends on: [%s]\n", node, strings.Join(deps, ", "))
		}
	}

	fmt.Println("\nResolved application order:")
	for i, v := range order {
		fmt.Printf("  %d. %s\n", i+1, v)
	}

	return nil
}
