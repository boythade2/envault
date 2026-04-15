package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	namespaceCmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage key namespaces within a vault",
	}

	assignCmd := &cobra.Command{
		Use:   "assign <vault> <namespace> <key>",
		Short: "Assign a key to a namespace",
		Args:  cobra.ExactArgs(3),
		RunE:  runNamespaceAssign,
	}

	unassignCmd := &cobra.Command{
		Use:   "unassign <vault> <namespace> <key>",
		Short: "Remove a key from a namespace",
		Args:  cobra.ExactArgs(3),
		RunE:  runNamespaceUnassign,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all namespaces and their keys",
		Args:  cobra.ExactArgs(1),
		RunE:  runNamespaceList,
	}

	namespaceCmd.AddCommand(assignCmd, unassignCmd, listCmd)
	rootCmd.AddCommand(namespaceCmd)
}

func runNamespaceAssign(cmd *cobra.Command, args []string) error {
	vaultFile, namespace, key := args[0], args[1], args[2]
	if _, err := os.Stat(vaultFile); err != nil {
		return fmt.Errorf("vault file not found: %s", vaultFile)
	}
	if err := vault.AssignNamespace(vaultFile, namespace, key); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "assigned %q to namespace %q\n", key, namespace)
	return nil
}

func runNamespaceUnassign(cmd *cobra.Command, args []string) error {
	vaultFile, namespace, key := args[0], args[1], args[2]
	if _, err := os.Stat(vaultFile); err != nil {
		return fmt.Errorf("vault file not found: %s", vaultFile)
	}
	if err := vault.UnassignNamespace(vaultFile, namespace, key); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "unassigned %q from namespace %q\n", key, namespace)
	return nil
}

func runNamespaceList(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	if _, err := os.Stat(vaultFile); err != nil {
		return fmt.Errorf("vault file not found: %s", vaultFile)
	}
	store, err := vault.LoadNamespaces(vaultFile)
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), vault.FormatNamespaceList(store))
	return nil
}
