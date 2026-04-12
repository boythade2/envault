package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage key aliases for a vault",
	}

	addAliasCmd := &cobra.Command{
		Use:   "add <vault> <alias> <key>",
		Short: "Add an alias pointing to an existing key",
		Args:  cobra.ExactArgs(3),
		RunE:  runAliasAdd,
	}

	removeAliasCmd := &cobra.Command{
		Use:   "remove <vault> <alias>",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(2),
		RunE:  runAliasRemove,
	}

	listAliasCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all aliases",
		Args:  cobra.ExactArgs(1),
		RunE:  runAliasList,
	}

	aliasCmd.AddCommand(addAliasCmd, removeAliasCmd, listAliasCmd)
	rootCmd.AddCommand(aliasCmd)
}

func runAliasAdd(cmd *cobra.Command, args []string) error {
	vaultPath, alias, key := args[0], args[1], args[2]
	if err := vault.AddAlias(vaultPath, alias, key); err != nil {
		return fmt.Errorf("add alias: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Alias %q -> %q added.\n", alias, key)
	return nil
}

func runAliasRemove(cmd *cobra.Command, args []string) error {
	vaultPath, alias := args[0], args[1]
	if err := vault.RemoveAlias(vaultPath, alias); err != nil {
		return fmt.Errorf("remove alias: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Alias %q removed.\n", alias)
	return nil
}

func runAliasList(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	store, err := vault.LoadAliases(vaultPath)
	if err != nil {
		return fmt.Errorf("load aliases: %w", err)
	}
	if len(store.Aliases) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No aliases defined.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ALIAS\tKEY")
	for alias, key := range store.Aliases {
		fmt.Fprintf(w, "%s\t%s\n", alias, key)
	}
	return w.Flush()
}
