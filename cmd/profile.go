package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage named environment profiles",
}

var profileAddCmd = &cobra.Command{
	Use:   "add <name> <vault-file>",
	Short: "Register a new profile pointing to a vault file",
	Args:  cobra.ExactArgs(2),
	RunE:  runProfileAdd,
}

var profileRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a registered profile",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileRemove,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered profiles",
	Args:  cobra.NoArgs,
	RunE:  runProfileList,
}

func init() {
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileRemoveCmd)
	profileCmd.AddCommand(profileListCmd)
	rootCmd.AddCommand(profileCmd)
}

func runProfileAdd(cmd *cobra.Command, args []string) error {
	dir, _ := os.Getwd()
	store, err := vault.LoadProfiles(dir)
	if err != nil {
		return err
	}
	if err := store.AddProfile(args[0], args[1]); err != nil {
		return err
	}
	if err := store.Save(dir); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Profile %q added → %s\n", args[0], args[1])
	return nil
}

func runProfileRemove(cmd *cobra.Command, args []string) error {
	dir, _ := os.Getwd()
	store, err := vault.LoadProfiles(dir)
	if err != nil {
		return err
	}
	if err := store.RemoveProfile(args[0]); err != nil {
		return err
	}
	if err := store.Save(dir); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Profile %q removed\n", args[0])
	return nil
}

func runProfileList(cmd *cobra.Command, args []string) error {
	dir, _ := os.Getwd()
	store, err := vault.LoadProfiles(dir)
	if err != nil {
		return err
	}
	if len(store.Profiles) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No profiles registered.")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVAULT FILE\tCREATED")
	for _, p := range store.Profiles {
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.VaultFile, p.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
