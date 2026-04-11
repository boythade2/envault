package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var restoreCmd = &cobra.Command{
	Use:   "restore <vault-file> <snapshot-label>",
	Short: "Restore a vault from a previously saved snapshot",
	Args:  cobra.ExactArgs(2),
	RunE:  runRestore,
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	label := args[1]

	passphrase, err := promptPassword("Enter passphrase: ")
	if err != nil {
		return fmt.Errorf("reading passphrase: %w", err)
	}

	restored, err := vault.RestoreSnapshot(vaultPath, label, passphrase)
	if err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Restored snapshot %q → %s (%d entries)\n",
		label, vaultPath, len(restored.Entries))
	return nil
}
