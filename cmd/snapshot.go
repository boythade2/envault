package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var snapshotLabel string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot <vault-file>",
	Short: "Save a named snapshot of a vault file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSnapshot,
}

var snapshotListCmd = &cobra.Command{
	Use:   "list <vault-file>",
	Short: "List all snapshots for a vault file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSnapshotList,
}

func init() {
	snapshotCmd.Flags().StringVarP(&snapshotLabel, "label", "l", "", "Optional label for the snapshot")
	snapshotCmd.AddCommand(snapshotListCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	name, err := vault.SaveSnapshot(vaultPath, v, snapshotLabel)
	if err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved: %s\n", name)
	return nil
}

func runSnapshotList(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	snapshots, err := vault.ListSnapshots(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}

	if len(snapshots) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No snapshots found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tCREATED\tENTRIES")
	for _, s := range snapshots {
		created := s.CreatedAt.Format(time.RFC3339)
		fmt.Fprintf(w, "%s\t%s\t%d\n", s.Name, created, s.EntryCount)
	}
	w.Flush()

	_ = os.Stderr
	return nil
}
