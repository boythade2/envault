package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var mergeStrategy string

func init() {
	mergeCmd := &cobra.Command{
		Use:   "merge <base-vault> <source-vault>",
		Short: "Merge entries from a source vault into a base vault",
		Long: `Merge entries from source-vault into base-vault.

Conflict resolution strategies:
  ours    - keep the base vault value (default)
  theirs  - overwrite with the source vault value
  error   - abort if any key conflicts`,
		Args: cobra.ExactArgs(2),
		RunE: runMerge,
	}

	mergeCmd.Flags().StringVarP(&mergeStrategy, "strategy", "s", "ours",
		"Conflict resolution strategy: ours, theirs, error")

	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	basePath := args[0]
	srcPath := args[1]

	passphrase, err := promptPassword("Passphrase: ")
	if err != nil {
		return fmt.Errorf("reading passphrase: %w", err)
	}

	dst, err := vault.LoadOrCreate(basePath, passphrase)
	if err != nil {
		return fmt.Errorf("loading base vault %q: %w", basePath, err)
	}

	src, err := vault.LoadOrCreate(srcPath, passphrase)
	if err != nil {
		return fmt.Errorf("loading source vault %q: %w", srcPath, err)
	}

	strategy := vault.MergeStrategy(mergeStrategy)
	result, err := vault.MergeVaults(dst, src, strategy)
	if err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	if err := dst.Save(basePath, passphrase); err != nil {
		return fmt.Errorf("saving merged vault: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Merge complete:\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  Added:   %d\n", len(result.Added))
	fmt.Fprintf(cmd.OutOrStdout(), "  Updated: %d\n", len(result.Updated))
	fmt.Fprintf(cmd.OutOrStdout(), "  Skipped: %d\n", len(result.Skipped))
	return nil
}
