package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var (
	maskShowChars int
	maskChar      string
	maskPatterns  []string
)

func init() {
	maskCmd := &cobra.Command{
		Use:   "mask <vault-file>",
		Short: "Display vault entries with sensitive values masked",
		Args:  cobra.ExactArgs(1),
		RunE:  runMask,
	}
	maskCmd.Flags().IntVar(&maskShowChars, "show-chars", 0, "Number of leading characters to reveal before masking")
	maskCmd.Flags().StringVar(&maskChar, "mask-char", "*", "Character used for masking")
	maskCmd.Flags().StringArrayVar(&maskPatterns, "pattern", nil, "Key regex patterns to mask (default: all keys)")
	rootCmd.AddCommand(maskCmd)
}

func runMask(cmd *cobra.Command, args []string) error {
	v, err := vault.LoadOrCreate(args[0])
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	opts := vault.MaskOptions{
		ShowChars: maskShowChars,
		MaskChar:  maskChar,
		Patterns:  maskPatterns,
	}

	masked, results := vault.MaskValues(v.Entries, opts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE\tMASKED")
	fmt.Fprintln(w, "---\t-----\t------")
	for i, e := range masked {
		maskedStr := "no"
		if results[i].Masked {
			maskedStr = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.Key, e.Value, maskedStr)
	}
	w.Flush()
	return nil
}
