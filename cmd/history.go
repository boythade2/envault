package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	historyCmd := &cobra.Command{
		Use:   "history <vault-file>",
		Short: "Show the change history of a vault file",
		Long:  "Displays a chronological log of additions, updates, and removals recorded for the given vault file.",
		Args:  cobra.ExactArgs(1),
		RunE:  runHistory,
	}
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	h, err := vault.LoadHistory(vaultPath)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}

	if len(h.Entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No history recorded for this vault.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTION\tKEY\tOLD VALUE\tNEW VALUE")
	fmt.Fprintln(w, "---------\t------\t---\t---------\t---------")
	for _, e := range h.Entries {
		oldVal := e.OldValue
		if oldVal == "" {
			oldVal = "-"
		}
		newVal := e.NewValue
		if newVal == "" {
			newVal = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Action,
			e.Key,
			oldVal,
			newVal,
		)
	}
	w.Flush()

	_ = os.Stderr // satisfy import if needed
	return nil
}
