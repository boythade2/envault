package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	dedupeCmd := &cobra.Command{
		Use:   "dedupe <vault-file>",
		Short: "Remove duplicate keys from a vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runDedupe,
	}
	dedupeCmd.Flags().String("strategy", "first",
		`Which occurrence to keep: "first" or "last"`)
	dedupeCmd.Flags().Bool("dry-run", false,
		"Preview duplicates without modifying the vault file")
	rootCmd.AddCommand(dedupeCmd)
}

func runDedupe(cmd *cobra.Command, args []string) error {
	path := args[0]
	strategy, _ := cmd.Flags().GetString("strategy")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	v, err := vault.LoadOrCreate(path)
	if err != nil {
		return fmt.Errorf("loading vault: %w", err)
	}

	opts := vault.DedupeOptions{
		Strategy: strategy,
		DryRun:   dryRun,
	}

	results, err := vault.DedupeEntries(v, opts)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, vault.FormatDedupeResults(results))

	if dryRun {
		fmt.Fprintln(os.Stdout, "(dry run — no changes written)")
		return nil
	}

	if len(results) == 0 {
		return nil
	}

	if err := v.Save(path); err != nil {
		return fmt.Errorf("saving vault: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Vault saved: %s\n", path)
	return nil
}
