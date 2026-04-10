package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envault/envault/internal/crypto"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [file]",
	Short: "Encrypt an environment file",
	Long:  `Encrypt a .env file using AES-GCM encryption. The encrypted output is written to [file].enc.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEncrypt,
}

var passphrase string

func init() {
	encryptCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for encryption (required)")
	_ = encryptCmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", inputPath, err)
	}

	ciphertext, err := crypto.Encrypt(plaintext, passphrase)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	outputPath := inputPath + ".enc"
	if err := os.WriteFile(outputPath, ciphertext, 0600); err != nil {
		return fmt.Errorf("failed to write encrypted file %q: %w", outputPath, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Encrypted %q -> %q\n", inputPath, outputPath)
	return nil
}
