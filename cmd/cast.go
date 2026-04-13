package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var castCmd = &cobra.Command{
	Use:   "cast <vault-file> <key>[,key...]",
	Short: "Coerce entry values to a target type (int, float, bool, string)",
	Args:  cobra.ExactArgs(2),
	RunE:  runCast,
}

func init() {
	castCmd.Flags().StringP("type", "t", "string", "target type: int, float, bool, string")
	castCmd.Flags().Bool("dry-run", false, "preview changes without writing")
	castCmd.Flags().StringP("passphrase", "p", "", "vault passphrase")
	rootCmd.AddCommand(castCmd)
}

func runCast(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	keys := strings.Split(args[1], ",")

	typeName, _ := cmd.Flags().GetString("type")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	passphrase, _ := cmd.Flags().GetString("passphrase")

	if passphrase == "" {
		passphrase = os.Getenv("ENVAULT_PASSPHRASE")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase required: use --passphrase or ENVAULT_PASSPHRASE")
	}

	castTo := vault.CastType(typeName)
	switch castTo {
	case vault.CastString, vault.CastInt, vault.CastFloat, vault.CastBool:
	default:
		return fmt.Errorf("unsupported type %q; choose from: int, float, bool, string", typeName)
	}

	results, err := vault.CastEntries(vaultPath, passphrase, keys, castTo, dryRun)
	if err != nil {
		return fmt.Errorf("cast failed: %w", err)
	}

	fmt.Print(vault.FormatCastResults(results))

	if dryRun {
		fmt.Println("(dry-run: no changes written)")
	}
	return nil
}
