package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var pinNote string

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin <vault-file> <key>",
		Short: "Pin a key to protect it from accidental modification",
		Args:  cobra.ExactArgs(2),
		RunE:  runPin,
	}
	pinCmd.Flags().StringVar(&pinNote, "note", "", "Optional note describing why the key is pinned")

	unpinCmd := &cobra.Command{
		Use:   "unpin <vault-file> <key>",
		Short: "Unpin a previously pinned key",
		Args:  cobra.ExactArgs(2),
		RunE:  runUnpin,
	}

	pinsCmd := &cobra.Command{
		Use:   "pins <vault-file>",
		Short: "List all pinned keys in a vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runListPins,
	}

	rootCmd.AddCommand(pinCmd, unpinCmd, pinsCmd)
}

func runPin(cmd *cobra.Command, args []string) error {
	vaultFile, key := args[0], args[1]
	if err := vault.PinKey(vaultFile, key, pinNote); err != nil {
		return fmt.Errorf("pin: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Pinned key %q in %s\n", key, vaultFile)
	return nil
}

func runUnpin(cmd *cobra.Command, args []string) error {
	vaultFile, key := args[0], args[1]
	if err := vault.UnpinKey(vaultFile, key); err != nil {
		return fmt.Errorf("unpin: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Unpinned key %q in %s\n", key, vaultFile)
	return nil
}

func runListPins(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	pl, err := vault.LoadPins(vaultFile)
	if err != nil {
		return fmt.Errorf("list pins: %w", err)
	}
	if len(pl.Pins) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No pinned keys.")
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%-30s  %-20s  %s\n", "KEY", "PINNED AT", "NOTE")
	for _, p := range pl.Pins {
		fmt.Fprintf(cmd.OutOrStdout(), "%-30s  %-20s  %s\n", p.Key, p.PinnedAt.Format("2006-01-02 15:04:05"), p.Note)
	}
	return nil
}
