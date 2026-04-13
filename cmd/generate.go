package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var generateCmd = &cobra.Command{
	Use:   "generate <vault-file> <key>",
	Short: "Generate a random secret and store it in the vault",
	Args:  cobra.ExactArgs(2),
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().IntP("length", "l", 32, "Length of the generated secret")
	generateCmd.Flags().Bool("upper", true, "Include uppercase letters")
	generateCmd.Flags().Bool("digits", true, "Include digits")
	generateCmd.Flags().Bool("symbols", false, "Include symbols")
	generateCmd.Flags().Bool("dry-run", false, "Preview the generated value without saving")
	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	key := args[1]

	length, err := cmd.Flags().GetInt("length")
	if err != nil {
		return err
	}
	useUpper, _ := cmd.Flags().GetBool("upper")
	useDigits, _ := cmd.Flags().GetBool("digits")
	useSymbols, _ := cmd.Flags().GetBool("symbols")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	opts := vault.GenerateOptions{
		Length:     length,
		UseUpper:   useUpper,
		UseDigits:  useDigits,
		UseSymbols: useSymbols,
		DryRun:     dryRun,
	}

	result, err := vault.GenerateAndStore(vaultFile, key, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	fmt.Println(vault.FormatGenerateResult(result, dryRun))
	return nil
}
