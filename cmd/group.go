package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage key groups within a vault",
	}

	addCmd := &cobra.Command{
		Use:   "add <vault> <group>",
		Short: "Create a new group",
		Args:  cobra.ExactArgs(2),
		RunE:  runGroupAdd,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <vault> <group>",
		Short: "Remove a group",
		Args:  cobra.ExactArgs(2),
		RunE:  runGroupRemove,
	}

	assignCmd := &cobra.Command{
		Use:   "assign <vault> <group> <key>",
		Short: "Assign a key to a group",
		Args:  cobra.ExactArgs(3),
		RunE:  runGroupAssign,
	}

	unassignCmd := &cobra.Command{
		Use:   "unassign <vault> <group> <key>",
		Short: "Remove a key from a group",
		Args:  cobra.ExactArgs(3),
		RunE:  runGroupUnassign,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all groups and their keys",
		Args:  cobra.ExactArgs(1),
		RunE:  runGroupList,
	}

	groupCmd.AddCommand(addCmd, removeCmd, assignCmd, unassignCmd, listCmd)
	rootCmd.AddCommand(groupCmd)
}

func runGroupAdd(cmd *cobra.Command, args []string) error {
	if err := vault.AddGroup(args[0], args[1]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Group %q created.\n", args[1])
	return nil
}

func runGroupRemove(cmd *cobra.Command, args []string) error {
	if err := vault.RemoveGroup(args[0], args[1]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Group %q removed.\n", args[1])
	return nil
}

func runGroupAssign(cmd *cobra.Command, args []string) error {
	if err := vault.AssignKeyToGroup(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Key %q assigned to group %q.\n", args[2], args[1])
	return nil
}

func runGroupUnassign(cmd *cobra.Command, args []string) error {
	if err := vault.UnassignKeyFromGroup(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Key %q removed from group %q.\n", args[2], args[1])
	return nil
}

func runGroupList(cmd *cobra.Command, args []string) error {
	gs, err := vault.LoadGroups(args[0])
	if err != nil {
		return err
	}
	if len(gs.Groups) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No groups defined.")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "GROUP\tKEYS")
	for name, g := range gs.Groups {
		keys := "-"
		if len(g.Keys) > 0 {
			keys = ""
			for i, k := range g.Keys {
				if i > 0 {
					keys += ", "
				}
				keys += k
			}
		}
		fmt.Fprintf(w, "%s\t%s\n", name, keys)
	}
	_ = w
	w.Flush()
	_ = os.Stdout
	return nil
}
