package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	statsCmd := &cobra.Command{
		Use:   "stats <vault-file>",
		Short: "Display aggregate statistics for a vault file",
		Long: `Compute and display statistics for a vault file, including total
key count, empty values, unique values, oldest/newest update
timestamps, and the most common key prefixes.`,
		Args: cobra.ExactArgs(1),
		RunE: runStats,
	}
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	path := args[0]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("vault file not found: %s", path)
	}

	v, err := vault.LoadOrCreate(path)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	stats := vault.ComputeStats(v)
	fmt.Print(vault.FormatStats(stats))
	return nil
}
