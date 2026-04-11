package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var (
	searchValueFlag bool
	searchVaultFile string
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for entries in the vault by key or value",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().BoolVarP(&searchValueFlag, "value", "v", false, "Search within values instead of keys")
	searchCmd.Flags().StringVarP(&searchVaultFile, "file", "f", ".envault", "Path to the vault file")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	v, err := vault.LoadOrCreate(searchVaultFile)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	results := v.Search(query, searchValueFlag)
	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No matching entries found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE\tUPDATED AT")
	fmt.Fprintln(w, "---\t-----\t----------")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Key, r.Value, r.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	w.Flush()

	_ = os.Stderr
	return nil
}
