package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	replaceCmd := &cobra.Command{
		Use:   "replace <vault-file>",
		Short: "Find and replace substrings in vault entry values",
		Args:  cobra.ExactArgs(1),
		RunE:  runReplace,
	}
	replaceCmd.Flags().String("old", "", "Substring to search for (required)")
	replaceCmd.Flags().String("new", "", "Replacement string")
	replaceCmd.Flags().StringSlice("keys", nil, "Limit replacement to specific keys")
	replaceCmd.Flags().Bool("all", true, "Replace all occurrences (default true)")
	replaceCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	_ = replaceCmd.MarkFlagRequired("old")
	rootCmd.AddCommand(replaceCmd)
}

func runReplace(cmd *cobra.Command, args []string) error {
	path := args[0]
	old, _ := cmd.Flags().GetString("old")
	newVal, _ := cmd.Flags().GetString("new")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	all, _ := cmd.Flags().GetBool("all")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	v, err := vault.LoadOrCreate(path)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	results, err := vault.ReplaceValues(v, path, vault.ReplaceOptions{
		Keys:   keys,
		Old:    old,
		New:    newVal,
		All:    all,
		DryRun: dryRun,
	})
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintln(os.Stdout, "[dry-run] "+vault.FormatReplaceResults(results))
	} else {
		fmt.Fprint(os.Stdout, vault.FormatReplaceResults(results))
	}
	return nil
}
