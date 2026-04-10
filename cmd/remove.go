package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var removePassphrase string

var removeCmd = &cobra.Command{
	Use:   "remove <vault-file> <KEY>",
	Short: "Remove an environment variable from a vault",
	Args:  cobra.ExactArgs(2),
	RunE:  runRemove,
}

func init() {
	removeCmd.Flags().StringVarP(&removePassphrase, "passphrase", "p", "", "Passphrase to decrypt the vault (required)")
	_ = removeCmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	key := args[1]

	v, err := vault.LoadOrCreate(vaultFile, removePassphrase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
		return err
	}

	if !v.Remove(key) {
		fmt.Fprintf(os.Stderr, "Key %q not found in vault\n", key)
		return fmt.Errorf("key %q not found", key)
	}

	if err := v.Save(vaultFile, removePassphrase); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving vault: %v\n", err)
		return err
	}

	fmt.Printf("Key %q removed from %s\n", key, vaultFile)
	return nil
}
