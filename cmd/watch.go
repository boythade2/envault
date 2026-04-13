package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <vault-file>",
		Short: "Track changes to a vault file",
		Long: `Watch records the current checksum of a vault file so you can
detect whether it has been modified since the last check.

Subcommands:
  record   Save the current state of the vault
  status   Report whether the vault has changed since the last record`,
	}

	watchCmd.AddCommand(&cobra.Command{
		Use:   "record <vault-file>",
		Short: "Record the current checksum of the vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatchRecord,
	})

	watchCmd.AddCommand(&cobra.Command{
		Use:   "status <vault-file>",
		Short: "Check whether the vault has changed since the last record",
		Args:  cobra.ExactArgs(1),
		RunE:  runWatchStatus,
	})

	rootCmd.AddCommand(watchCmd)
}

func runWatchRecord(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	if err := requireVaultFile(vaultPath); err != nil {
		return err
	}
	if err := vault.SaveWatchState(vaultPath); err != nil {
		return fmt.Errorf("record watch state: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Watch state recorded for %s\n", vaultPath)
	return nil
}

func runWatchStatus(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	if err := requireVaultFile(vaultPath); err != nil {
		return err
	}
	changed, err := vault.HasChanged(vaultPath)
	if err != nil {
		return fmt.Errorf("check watch state: %w", err)
	}
	if changed {
		fmt.Fprintf(cmd.OutOrStdout(), "CHANGED: %s has been modified since the last record\n", vaultPath)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "UNCHANGED: %s matches the recorded state\n", vaultPath)
	}
	return nil
}

// requireVaultFile returns an error if the given path does not exist or is not
// accessible, providing a consistent error message across watch subcommands.
func requireVaultFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("vault file not found: %s", path)
		}
		return fmt.Errorf("vault file not accessible: %w", err)
	}
	return nil
}
