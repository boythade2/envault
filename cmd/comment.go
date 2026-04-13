package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	commentCmd := &cobra.Command{
		Use:   "comment",
		Short: "Manage inline comments for vault entries",
	}

	setCmd := &cobra.Command{
		Use:   "set <vault-file> <key> <comment>",
		Short: "Set a comment on a vault entry",
		Args:  cobra.ExactArgs(3),
		RunE:  runCommentSet,
	}

	unsetCmd := &cobra.Command{
		Use:   "unset <vault-file> <key>",
		Short: "Remove a comment from a vault entry",
		Args:  cobra.ExactArgs(2),
		RunE:  runCommentUnset,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault-file>",
		Short: "List all comments in a vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runCommentList,
	}

	commentCmd.AddCommand(setCmd, unsetCmd, listCmd)
	rootCmd.AddCommand(commentCmd)
}

func runCommentSet(cmd *cobra.Command, args []string) error {
	vaultFile, key, comment := args[0], args[1], args[2]
	if err := vault.SetComment(vaultFile, key, comment); err != nil {
		return fmt.Errorf("set comment: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Comment set for %q\n", key)
	return nil
}

func runCommentUnset(cmd *cobra.Command, args []string) error {
	vaultFile, key := args[0], args[1]
	if err := vault.RemoveComment(vaultFile, key); err != nil {
		return fmt.Errorf("remove comment: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Comment removed for %q\n", key)
	return nil
}

func runCommentList(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	cs, err := vault.LoadComments(vaultFile)
	if err != nil {
		return fmt.Errorf("load comments: %w", err)
	}
	if len(cs.Comments) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No comments found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tCOMMENT")
	for k, c := range cs.Comments {
		fmt.Fprintf(w, "%s\t%s\n", k, c)
	}
	return w.Flush()
}
