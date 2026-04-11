package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff <base-vault> <other-vault>",
	Short: "Compare two vault files and show differences",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	passphrase := os.Getenv("ENVAULT_PASSPHRASE")
	if passphrase == "" {
		return fmt.Errorf("ENVAULT_PASSPHRASE environment variable is not set")
	}

	baseVault, err := vault.LoadOrCreate(args[0], passphrase)
	if err != nil {
		return fmt.Errorf("failed to load base vault %q: %w", args[0], err)
	}

	otherVault, err := vault.LoadOrCreate(args[1], passphrase)
	if err != nil {
		return fmt.Errorf("failed to load other vault %q: %w", args[1], err)
	}

	result := vault.Diff(baseVault, otherVault)

	if !result.HasChanges() {
		fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
		return nil
	}

	w := cmd.OutOrStdout()
	for _, key := range result.Added {
		fmt.Fprintf(w, "+ %s\n", key)
	}
	for _, key := range result.Removed {
		fmt.Fprintf(w, "- %s\n", key)
	}
	for _, key := range result.Changed {
		fmt.Fprintf(w, "~ %s\n", key)
	}

	fmt.Fprintf(w, "\nSummary: %d added, %d removed, %d changed\n",
		len(result.Added), len(result.Removed), len(result.Changed))

	return nil
}
