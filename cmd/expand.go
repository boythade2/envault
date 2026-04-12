package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var expandUseOS bool

var expandCmd = &cobra.Command{
	Use:   "expand <vault-file>",
	Short: "Expand ${KEY} references within vault values",
	Long: `Resolves ${KEY} placeholders in vault entry values by substituting
the value of the referenced key from the same vault.

Use --use-os to fall back to OS environment variables for unresolved references.`,
	Args: cobra.ExactArgs(1),
	RunE: runExpand,
}

func init() {
	expandCmd.Flags().BoolVar(&expandUseOS, "use-os", false, "fall back to OS env vars for unresolved references")
	rootCmd.AddCommand(expandCmd)
}

func runExpand(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return fmt.Errorf("loading vault: %w", err)
	}

	results, err := vault.ExpandVaultRefs(v, expandUseOS)
	if err != nil {
		return fmt.Errorf("expanding references: %w", err)
	}

	changed := 0
	for _, r := range results {
		if r.Changed {
			v.Entries[r.Key] = vault.Entry{
				Value:     r.Expanded,
				UpdatedAt: v.Entries[r.Key].UpdatedAt,
				Tags:      v.Entries[r.Key].Tags,
			}
			changed++
		}
	}

	if changed == 0 {
		fmt.Fprintln(os.Stdout, "no references to expand")
		return nil
	}

	if err := v.Save(vaultPath); err != nil {
		return fmt.Errorf("saving vault: %w", err)
	}

	fmt.Fprintf(os.Stdout, "expanded %d reference(s):\n", changed)
	fmt.Fprint(os.Stdout, vault.FormatExpandResults(results))
	return nil
}
