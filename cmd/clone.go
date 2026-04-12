package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/user/envault/internal/vault"
)

var cloneWriteMeta bool

func init() {
	cloneCmd := &cobra.Command{
		Use:   "clone <src-vault> <dst-vault>",
		Short: "Clone a vault to a new file with a new passphrase",
		Long: `Clone copies every entry from the source vault into a new
vault file, re-encrypting all values with a fresh passphrase.
The destination file must not already exist.`,
		Args: cobra.ExactArgs(2),
		RunE: runClone,
	}
	cloneCmd.Flags().BoolVar(&cloneWriteMeta, "write-meta", false,
		"write a .clone.json sidecar recording the clone's origin")
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	fmt.Fprint(os.Stderr, "Source passphrase: ")
	oldRaw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return fmt.Errorf("read source passphrase: %w", err)
	}

	fmt.Fprint(os.Stderr, "New passphrase: ")
	newRaw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return fmt.Errorf("read new passphrase: %w", err)
	}

	fmt.Fprint(os.Stderr, "Confirm new passphrase: ")
	confirmRaw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return fmt.Errorf("read confirmation: %w", err)
	}

	if string(newRaw) != string(confirmRaw) {
		return fmt.Errorf("passphrases do not match")
	}

	result, err := vault.CloneVault(srcPath, dstPath, string(oldRaw), string(newRaw))
	if err != nil {
		return fmt.Errorf("clone vault: %w", err)
	}

	if cloneWriteMeta {
		if err := vault.WriteCloneMeta(result); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not write clone meta: %v\n", err)
		}
	}

	fmt.Printf("Cloned %d entr", result.EntriesCount)
	if result.EntriesCount == 1 {
		fmt.Print("y")
	} else {
		fmt.Print("ies")
	}
	fmt.Printf(" from %s → %s\n", result.Source, result.Destination)
	return nil
}
