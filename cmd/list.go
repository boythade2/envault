package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/envault/internal/vault"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked environment files in the vault index",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	v, err := vault.LoadOrCreate(dir)
	if err != nil {
		return fmt.Errorf("loading vault index: %w", err)
	}

	if len(v.Entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No entries found. Use 'envault encrypt' to add files.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tFILE\tENCRYPTED\tUPDATED")
	fmt.Fprintln(w, "----\t----\t---------\t-------")
	for _, entry := range v.Entries {
		encStatus := "no"
		if entry.Encrypted {
			encStatus = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			entry.Name,
			entry.FilePath,
			encStatus,
			entry.UpdatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return w.Flush()
}
