package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"envault/internal/crypto"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [file]",
	Short: "Decrypt an encrypted .env file",
	Long:  `Decrypt an encrypted .env file using a passphrase and write the plaintext to stdout or a file.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDecrypt,
}

func init() {
	decryptCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	rootCmd.AddCommand(decryptCmd)
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	ciphertext, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", inputFile, err)
	}

	fmt.Fprint(os.Stderr, "Enter passphrase: ")
	passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %w", err)
	}

	plaintext, err := crypto.Decrypt(strings.TrimSpace(string(ciphertext)), string(passphrase))
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	outputFile, _ := cmd.Flags().GetString("output")
	if outputFile == "" {
		fmt.Print(plaintext)
		return nil
	}

	if err := os.WriteFile(outputFile, []byte(plaintext), 0600); err != nil {
		return fmt.Errorf("failed to write output file %q: %w", outputFile, err)
	}

	fmt.Fprintf(os.Stderr, "Decrypted content written to %q\n", outputFile)
	return nil
}
