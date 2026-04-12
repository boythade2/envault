package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	chmodCmd := &cobra.Command{
		Use:   "chmod",
		Short: "Manage key permissions (read-only flags)",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault> <key>",
		Short: "Mark a key as read-only",
		Args:  cobra.ExactArgs(2),
		RunE:  runChmodSet,
	}

	unsetCmd := &cobra.Command{
		Use:   "unset <vault> <key>",
		Short: "Remove read-only flag from a key",
		Args:  cobra.ExactArgs(2),
		RunE:  runChmodUnset,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all key permissions",
		Args:  cobra.ExactArgs(1),
		RunE:  runChmodList,
	}

	chmodCmd.AddCommand(setCmd, unsetCmd, listCmd)
	rootCmd.AddCommand(chmodCmd)
}

func runChmodSet(cmd *cobra.Command, args []string) error {
	vaultPath, key := args[0], args[1]
	if err := vault.SetPermission(vaultPath, key, true); err != nil {
		return fmt.Errorf("chmod set: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "key %q marked as read-only\n", key)
	return nil
}

func runChmodUnset(cmd *cobra.Command, args []string) error {
	vaultPath, key := args[0], args[1]
	if err := vault.RemovePermission(vaultPath, key); err != nil {
		return fmt.Errorf("chmod unset: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "read-only flag removed from key %q\n", key)
	return nil
}

func runChmodList(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	pm, err := vault.LoadPermissions(vaultPath)
	if err != nil {
		return fmt.Errorf("chmod list: %w", err)
	}
	if len(pm) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no permissions set")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tREAD-ONLY")
	for _, p := range pm {
		ro := "false"
		if p.ReadOnly {
			ro = "true"
		}
		fmt.Fprintf(w, "%s\t%s\n", p.Key, ro)
	}
	return w.Flush()
}
