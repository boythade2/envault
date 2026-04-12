package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	injectCmd := &cobra.Command{
		Use:   "inject <vault-file> [-- command...]",
		Short: "Inject vault entries into the environment, optionally running a command",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runInject,
	}
	injectCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing environment variables")
	injectCmd.Flags().StringP("prefix", "p", "", "Only inject keys with this prefix")
	injectCmd.Flags().BoolP("dry-run", "n", false, "Print what would be injected without modifying the environment")
	rootCmd.AddCommand(injectCmd)
}

func runInject(cmd *cobra.Command, args []string) error {
	vaultFile := args[0]
	subArgs := args[1:]

	overwrite, _ := cmd.Flags().GetBool("overwrite")
	prefix, _ := cmd.Flags().GetString("prefix")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	v, err := vault.LoadOrCreate(vaultFile)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	opts := vault.InjectOptions{Overwrite: overwrite, Prefix: prefix}

	if dryRun {
		for k, e := range v.Entries {
			if prefix != "" && !strings.HasPrefix(k, prefix) {
				continue
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, e.Value)
		}
		return nil
	}

	result, err := vault.InjectEnv(v, opts)
	if err != nil {
		return fmt.Errorf("inject: %w", err)
	}

	fmt.Fprint(cmd.OutOrStdout(), vault.FormatInjectResult(result))

	if len(subArgs) > 0 {
		c := exec.Command(subArgs[0], subArgs[1:]...)
		c.Env = os.Environ()
		c.Stdout = cmd.OutOrStdout()
		c.Stderr = cmd.ErrOrStderr()
		c.Stdin = os.Stdin
		if err := c.Run(); err != nil {
			return fmt.Errorf("command failed: %w", err)
		}
	}

	return nil
}
