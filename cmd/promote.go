package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var (
	promotePassphrase string
	promoteKeys       []string
	promoteOverwrite  bool
	promoteDryRun     bool
)

func init() {
	promoteCmd := &cobra.Command{
		Use:   "promote <src-vault> <dst-vault>",
		Short: "Promote entries from one vault to another (e.g. staging → production)",
		Args:  cobra.ExactArgs(2),
		RunE:  runPromote,
	}

	promoteCmd.Flags().StringVarP(&promotePassphrase, "passphrase", "p", "", "Passphrase for both vaults")
	promoteCmd.Flags().StringSliceVarP(&promoteKeys, "keys", "k", nil, "Comma-separated list of keys to promote (default: all)")
	promoteCmd.Flags().BoolVar(&promoteOverwrite, "overwrite", false, "Overwrite existing keys in destination")
	promoteCmd.Flags().BoolVar(&promoteDryRun, "dry-run", false, "Preview changes without writing")
	_ = promoteCmd.MarkFlagRequired("passphrase")

	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("source vault not found: %s", srcPath)
	}

	// Normalise keys to uppercase
	normKeys := make([]string, len(promoteKeys))
	for i, k := range promoteKeys {
		normKeys[i] = strings.ToUpper(strings.TrimSpace(k))
	}

	opts := vault.PromoteOptions{
		Keys:      normKeys,
		Overwrite: promoteOverwrite,
		DryRun:    promoteDryRun,
	}

	results, err := vault.PromoteEntries(srcPath, dstPath, promotePassphrase, opts)
	if err != nil {
		return fmt.Errorf("promote failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No entries to promote.")
		return nil
	}

	vault.FormatPromoteResults(results, promoteDryRun, os.Stdout)

	promoted := 0
	for _, r := range results {
		if !r.Skipped {
			promoted++
		}
	}

	if promoteDryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "\n[dry-run] %d key(s) would be promoted.\n", promoted)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "\n%d key(s) promoted from %s → %s\n", promoted, srcPath, dstPath)
	}
	return nil
}
