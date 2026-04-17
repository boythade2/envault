package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envault/internal/vault"
)

var regexMatchValue bool

func init() {
	regexCmd := &cobra.Command{
		Use:   "regex <vault-file> <pattern>",
		Short: "Filter vault entries by a regular expression",
		Args:  cobra.ExactArgs(2),
		RunE:  runRegex,
	}
	regexCmd.Flags().BoolVar(&regexMatchValue, "match-value", false, "also match against entry values")
	rootCmd.AddCommand(regexCmd)
}

func runRegex(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	pattern := args[1]

	v, err := vault.LoadOrCreate(vaultFile)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	results, err := vault.RegexFilter(v, pattern, regexMatchValue)
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, vault.FormatRegexResults(results))
	return nil
}
