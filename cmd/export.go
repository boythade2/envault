package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var exportFormat string

var exportCmd = &cobra.Command{
	Use:   "export <vault-file>",
	Short: "Export vault entries to a plaintext file",
	Long: `Export decrypted vault entries to a file in dotenv or JSON format.

Example:
  envault export myproject.vault --format dotenv --output .env
  envault export myproject.vault --format json --output env.json`,
	Args: cobra.ExactArgs(1),
	RunE: runExport,
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: dotenv or json")
	exportCmd.Flags().StringP("output", "o", "", "Output file path (required)")
	_ = exportCmd.MarkFlagRequired("output")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	outputPath, _ := cmd.Flags().GetString("output")

	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	if len(v.Entries) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: vault is empty, exporting empty file")
	}

	fmt := vault.ExportFormat(exportFormat)
	if err := v.Export(outputPath, fmt); err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	cmd.Printf("Exported %d entries to %s\n", len(v.Entries), outputPath)
	return nil
}
