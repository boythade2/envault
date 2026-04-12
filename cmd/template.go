package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var templateOutput string

func init() {
	templateCmd := &cobra.Command{
		Use:   "template <vault-file> <template-file>",
		Short: "Render a template file using vault entries",
		Long: `Reads a template file containing {{KEY}} placeholders and replaces
them with values from the specified vault. Writes the result to stdout
or to a file when --output is provided.`,
		Args: cobra.ExactArgs(2),
		RunE: runTemplate,
	}
	templateCmd.Flags().StringVarP(&templateOutput, "output", "o", "", "Write rendered output to file instead of stdout")
	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	templatePath := args[1]

	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}

	result, err := vault.RenderTemplate(templatePath, v)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	if len(result.Missing) > 0 {
		fmt.Fprintf(cmd.ErrOrStderr(), "warning: missing keys: %s\n",
			strings.Join(result.Missing, ", "))
	}

	if templateOutput != "" {
		if err := os.WriteFile(templateOutput, []byte(result.Rendered), 0600); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "rendered template written to %s\n", templateOutput)
		return nil
	}

	fmt.Fprint(cmd.OutOrStdout(), result.Rendered)
	return nil
}
