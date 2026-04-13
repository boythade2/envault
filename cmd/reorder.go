package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	var dryRun bool

	reorderCmd := &cobra.Command{
		Use:   "reorder <vault-file> <key1> [key2 ...]",
		Short: "Reorder vault entries by moving specified keys to the front",
		Long: `Reorder moves the listed keys to the beginning of the vault entry list
in the order provided. All remaining keys retain their original relative order.

Use --dry-run to preview the changes without writing to disk.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReorder(args[0], args[1:], dryRun)
		},
	}

	reorderCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")
	rootCmd.AddCommand(reorderCmd)
}

func runReorder(vaultPath string, keys []string, dryRun bool) error {
	results, err := vault.ReorderEntries(vaultPath, keys, dryRun)
	if err != nil {
		return fmt.Errorf("reorder: %w", err)
	}
	fmt.Fprintln(os.Stdout, vault.FormatReorderResults(results, dryRun))
	return nil
}
