package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate <vault-file>",
	Short: "Re-encrypt a vault with a new passphrase",
	Long: `Rotate the encryption passphrase for a vault file.
You will be prompted for the current passphrase and a new passphrase.
All entries are decrypted with the old passphrase and re-encrypted with the new one.`,
	Args: cobra.ExactArgs(1),
	RunE: runRotate,
}

func init() {
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	currentPassphrase, err := promptPassword("Current passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read current passphrase: %w", err)
	}

	v, err := vault.LoadOrCreate(vaultPath, currentPassphrase)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	newPassphrase, err := promptPassword("New passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read new passphrase: %w", err)
	}

	confirmPassphrase, err := promptPassword("Confirm new passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase confirmation: %w", err)
	}

	if newPassphrase != confirmPassphrase {
		return fmt.Errorf("passphrases do not match")
	}

	if err := v.Rotate(vaultPath, newPassphrase); err != nil {
		return fmt.Errorf("failed to rotate passphrase: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Passphrase rotated successfully for %s\n", vaultPath)
	return nil
}

func promptPassword(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	var passphrase string
	if _, err := fmt.Fscanln(os.Stdin, &passphrase); err != nil {
		return "", err
	}
	return passphrase, nil
}
