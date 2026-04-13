package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	var reason string
	var dryRun bool

	archiveCmd := &cobra.Command{
		Use:   "archive <vault-file> <key> [key...]",
		Short: "Archive one or more keys from a vault",
		Long:  "Move keys from a vault into an archive file, preserving their values for reference.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runArchive(cmd, args, reason, dryRun)
		},
	}

	archiveCmd.Flags().StringVar(&reason, "reason", "", "reason for archiving")
	archiveCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")

	archiveListCmd := &cobra.Command{
		Use:   "list <vault-file>",
		Short: "List archived keys",
		Args:  cobra.ExactArgs(1),
		RunE:  runArchiveList,
	}

	archiveCmd.AddCommand(archiveListCmd)
	rootCmd.AddCommand(archiveCmd)
}

func runArchive(cmd *cobra.Command, args []string, reason string, dryRun bool) error {
	vaultFile := args[0]
	keys := args[1:]
	passphrase, err := promptPassword("Passphrase: ")
	if err != nil {
		return err
	}
	archived, err := vault.ArchiveEntries(vaultFile, passphrase, keys, reason, dryRun)
	if err != nil {
		return err
	}
	fmt.Print(vault.FormatArchiveResults(archived, dryRun))
	return nil
}

func runArchiveList(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	store, err := vault.LoadArchive(vaultFile)
	if err != nil {
		return err
	}
	if len(store.Entries) == 0 {
		fmt.Println("no archived entries")
		return nil
	}
	fmt.Printf("%-30s %-30s %s\n", "KEY", "ARCHIVED AT", "REASON")
	fmt.Println(strings.Repeat("-", 72))
	for _, e := range store.Entries {
		fmt.Printf("%-30s %-30s %s\n", e.Key, e.ArchivedAt.Format("2006-01-02 15:04:05"), e.Reason)
	}
	return nil
}
