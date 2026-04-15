package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var placeholderDryRun bool
var placeholderKeys []string

func init() {
	placeholderCmd := &cobra.Command{
		Use:   "placeholder <vault-file>",
		Short: "Resolve {{KEY}} placeholders in vault values",
		Long: `Scans vault entry values for {{KEY}} tokens and replaces them with
the value of the referenced key from the same vault or from OS environment
variables. Use --dry-run to preview changes without modifying the file.`,
		Args: cobra.ExactArgs(1),
		RunE: runPlaceholder,
	}

	placeholderCmd.Flags().BoolVar(&placeholderDryRun, "dry-run", false, "preview changes without writing")
	placeholderCmd.Flags().StringSliceVar(&placeholderKeys, "keys", nil, "comma-separated list of keys to process (default: all)")

	rootCmd.AddCommand(placeholderCmd)
}

func runPlaceholder(cmd *cobra.Command, args []string) error {
	path := args[0]

	v, err := vault.LoadOrCreate(path)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	results, err := vault.ResolvePlaceholders(v, placeholderKeys, placeholderDryRun)
	if err != nil {
		return fmt.Errorf("placeholder resolution failed: %w", err)
	}

	fmt.Fprint(cmd.OutOrStdout(), vault.FormatPlaceholderResults(results))

	if placeholderDryRun {
		fmt.Fprintln(cmd.OutOrStdout(), "(dry-run: no changes written)")
		return nil
	}

	if err := v.Save(path); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}
	if changed > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%d placeholder(s) resolved and saved.\n", changed)
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("vault file not found after save: %w", err)
	}
	return nil
}
