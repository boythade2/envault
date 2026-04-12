package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var copyOverwrite bool
var copyDestKey string
var copyPassphrase string

func init() {
	copyCmd := &cobra.Command{
		Use:   "copy <src-vault> <dst-vault> <key>",
		Short: "Copy a key from one vault to another",
		Long: `Copy a single key (and its value) from a source vault into a destination vault.

The destination vault is created if it does not exist.
Use --dest-key to rename the key in the destination.
Use --overwrite to replace an existing key in the destination.`,
		Args: cobra.ExactArgs(3),
		RunE: runCopy,
	}

	copyCmd.Flags().BoolVar(&copyOverwrite, "overwrite", false, "overwrite key if it already exists in destination")
	copyCmd.Flags().StringVar(&copyDestKey, "dest-key", "", "rename the key in the destination vault")
	copyCmd.Flags().StringVar(&copyPassphrase, "passphrase", "", "passphrase for both vaults")
	_ = copyCmd.MarkFlagRequired("passphrase")

	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]
	srcKey := args[2]

	res, err := vault.CopyEntry(srcPath, dstPath, srcKey, copyDestKey, copyPassphrase, copyOverwrite)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	if res.Overwrote {
		fmt.Fprintf(cmd.OutOrStdout(), "overwritten: %s → %s:%s\n", res.SourceKey, dstPath, res.DestKey)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "copied: %s → %s:%s\n", res.SourceKey, dstPath, res.DestKey)
	}
	return nil
}
