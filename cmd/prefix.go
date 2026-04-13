package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	var dryRun bool
	var overwrite bool
	var keys []string

	addCmd := &cobra.Command{
		Use:   "prefix-add <vault-file> <prefix>",
		Short: "Prepend a prefix to vault keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrefixAdd(args[0], args[1], keys, dryRun, overwrite)
		},
	}
	addCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without modifying the vault")
	addCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys on conflict")
	addCmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated list of keys to process (default: all)")

	removeCmd := &cobra.Command{
		Use:   "prefix-remove <vault-file> <prefix>",
		Short: "Strip a prefix from vault keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrefixRemove(args[0], args[1], keys, dryRun, overwrite)
		},
	}
	removeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without modifying the vault")
	removeCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys on conflict")
	removeCmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated list of keys to process (default: all)")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
}

func runPrefixAdd(path, prefix string, keys []string, dryRun, overwrite bool) error {
	v, err := vault.LoadOrCreate(path)
	if err != nil {
		return fmt.Errorf("loading vault: %w", err)
	}
	opts := vault.PrefixOptions{Prefix: prefix, DryRun: dryRun, Overwrite: overwrite}
	results, err := vault.AddKeyPrefix(v, keys, opts)
	if err != nil {
		return err
	}
	if !dryRun {
		if err := v.Save(path); err != nil {
			return fmt.Errorf("saving vault: %w", err)
		}
	}
	fmt.Fprintln(os.Stdout, vault.FormatPrefixResults(results))
	if dryRun {
		fmt.Fprintln(os.Stdout, "(dry-run: no changes written)")
	}
	return nil
}

func runPrefixRemove(path, prefix string, keys []string, dryRun, overwrite bool) error {
	v, err := vault.LoadOrCreate(path)
	if err != nil {
		return fmt.Errorf("loading vault: %w", err)
	}
	opts := vault.PrefixOptions{Prefix: prefix, DryRun: dryRun, Overwrite: overwrite}
	results, err := vault.RemoveKeyPrefix(v, keys, opts)
	if err != nil {
		return err
	}
	if !dryRun {
		if err := v.Save(path); err != nil {
			return fmt.Errorf("saving vault: %w", err)
		}
	}
	out := vault.FormatPrefixResults(results)
	if strings.TrimSpace(out) == "" {
		out = "no matching keys found"
	}
	fmt.Fprintln(os.Stdout, out)
	if dryRun {
		fmt.Fprintln(os.Stdout, "(dry-run: no changes written)")
	}
	return nil
}
