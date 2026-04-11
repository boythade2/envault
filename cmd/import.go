package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import environment variables from a .env or JSON file into the vault",
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

func init() {
	importCmd.Flags().StringP("vault", "v", "vault.json", "Path to the vault file")
	importCmd.Flags().StringP("format", "f", "", "File format: dotenv or json (auto-detected if omitted)")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	src := args[0]
	vaultPath, _ := cmd.Flags().GetString("vault")
	format, _ := cmd.Flags().GetString("format")

	if format == "" {
		switch strings.ToLower(filepath.Ext(src)) {
		case ".json":
			format = "json"
		default:
			format = "dotenv"
		}
	}

	var entries map[string]string
	var err error

	switch format {
	case "json":
		entries, err = vault.ImportJSON(src)
	case "dotenv":
		entries, err = vault.ImportDotenv(src)
	default:
		return fmt.Errorf("unsupported format %q; use dotenv or json", format)
	}
	if err != nil {
		return fmt.Errorf("importing file: %w", err)
	}

	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return fmt.Errorf("loading vault: %w", err)
	}

	for k, val := range entries {
		v.Add(k, val)
	}

	if err := v.Save(vaultPath); err != nil {
		return fmt.Errorf("saving vault: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Imported %d variable(s) from %s\n", len(entries), src)
	return nil
}
