package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var addPassphrase string

var addCmd = &cobra.Command{
	Use:   "add <vault-file> <KEY> <VALUE>",
	Short: "Add or update an environment variable in a vault",
	Args:  cobra.ExactArgs(3),
	RunE:  runAdd,
}

func init() {
	addCmd.Flags().StringVarP(&addPassphrase, "passphrase", "p", "", "Passphrase to encrypt the vault (required)")
	_ = addCmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	key := args[1]
	value := args[2]

	v, err := vault.LoadOrCreate(vaultFile, addPassphrase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
		return err
	}

	v.Set(key, value)

	if err := v.Save(vaultFile, addPassphrase); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving vault: %v\n", err)
		return err
	}

	fmt.Printf("Key %q added/updated in %s\n", key, vaultFile)
	return nil
}
