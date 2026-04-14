package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	checkpointCmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "Manage named checkpoints of vault state",
	}

	saveCmd := &cobra.Command{
		Use:   "save <vault-file> <label>",
		Short: "Save current vault state as a named checkpoint",
		Args:  cobra.ExactArgs(2),
		RunE:  runCheckpointSave,
	}

	restoreCmd := &cobra.Command{
		Use:   "restore <vault-file> <label>",
		Short: "Restore vault entries from a named checkpoint",
		Args:  cobra.ExactArgs(2),
		RunE:  runCheckpointRestore,
	}

	listCmd := &cobra.Command{
		Use:   "list <vault-file>",
		Short: "List all saved checkpoints",
		Args:  cobra.ExactArgs(1),
		RunE:  runCheckpointList,
	}

	checkpointCmd.AddCommand(saveCmd, restoreCmd, listCmd)
	rootCmd.AddCommand(checkpointCmd)
}

func runCheckpointSave(cmd *cobra.Command, args []string) error {
	vaultPath, label := args[0], args[1]
	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return err
	}
	if err := vault.SaveCheckpoint(vaultPath, label, v); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Checkpoint %q saved.\n", label)
	return nil
}

func runCheckpointRestore(cmd *cobra.Command, args []string) error {
	vaultPath, label := args[0], args[1]
	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return err
	}
	if err := vault.RestoreCheckpoint(vaultPath, label, v); err != nil {
		return err
	}
	if err := v.Save(vaultPath); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Checkpoint %q restored.\n", label)
	return nil
}

func runCheckpointList(cmd *cobra.Command, args []string) error {
	checkpoints, err := vault.ListCheckpoints(args[0])
	if err != nil {
		return err
	}
	if len(checkpoints) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No checkpoints found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "LABEL\tCREATED AT\tKEYS")
	for _, cp := range checkpoints {
		fmt.Fprintf(w, "%s\t%s\t%d\n", cp.Label, cp.CreatedAt.Format("2006-01-02 15:04:05"), len(cp.Entries))
	}
	return w.Flush()
}
