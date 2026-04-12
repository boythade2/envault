package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var comparePassA string
var comparePassB string

var compareCmd = &cobra.Command{
	Use:   "compare <vaultA> <vaultB>",
	Short: "Compare two encrypted vault files and show key differences",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompare,
}

func init() {
	compareCmd.Flags().StringVar(&comparePassA, "pass-a", "", "Passphrase for vault A")
	compareCmd.Flags().StringVar(&comparePassB, "pass-b", "", "Passphrase for vault B (defaults to --pass-a)")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	pathA, pathB := args[0], args[1]

	if comparePassA == "" {
		return fmt.Errorf("--pass-a is required")
	}
	passB := comparePassB
	if passB == "" {
		passB = comparePassA
	}

	vA, err := vault.LoadOrCreate(pathA, comparePassA)
	if err != nil {
		return fmt.Errorf("loading vault A: %w", err)
	}

	vB, err := vault.LoadOrCreate(pathB, passB)
	if err != nil {
		return fmt.Errorf("loading vault B: %w", err)
	}

	result := vault.CompareVaults(vA, vB)

	if len(result.OnlyInA) == 0 && len(result.OnlyInB) == 0 && len(result.Changed) == 0 {
		fmt.Fprintln(os.Stdout, "Vaults are identical.")
		return nil
	}

	fmt.Fprintln(os.Stdout, result.Summary())
	return nil
}
