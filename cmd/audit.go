package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit <vault-file>",
		Short: "Show the audit log for a vault file",
		Args:  cobra.ExactArgs(1),
		RunE:  runAudit,
	}
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	log, err := vault.LoadAuditLog(vaultPath)
	if err != nil {
		return fmt.Errorf("loading audit log: %w", err)
	}

	if len(log.Events) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No audit events recorded.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTION\tKEY\tDETAILS")
	fmt.Fprintln(w, "---------\t------\t---\t-------")
	for _, e := range log.Events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Action,
			e.Key,
			e.Details,
		)
	}
	_ = w.Flush()

	fmt.Fprintf(os.Stderr, "\nTotal events: %d\n", len(log.Events))
	return nil
}
