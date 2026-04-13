package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	lintCheckCmd := &cobra.Command{
		Use:   "lintcheck <vault-file>",
		Short: "Run extended lint rules against a vault file",
		Args:  cobra.ExactArgs(1),
		RunE:  runLintCheck,
	}
	lintCheckCmd.Flags().StringP("level", "l", "", "filter by level: warn or error")
	rootCmd.AddCommand(lintCheckCmd)
}

func runLintCheck(cmd *cobra.Command, args []string) error {
	path := args[0]
	levelFilter, _ := cmd.Flags().GetString("level")

	v, err := vault.LoadOrCreate(path)
	if err != nil {
		return fmt.Errorf("loading vault: %w", err)
	}

	rules := vault.DefaultLintRules()
	results := vault.RunLintRules(v, rules)

	// Apply level filter
	var filtered []vault.LintRuleResult
	for _, r := range results {
		if levelFilter == "" || string(r.Level) == levelFilter {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) == 0 {
		fmt.Println("No lint issues found.")
		return nil
	}

	// Sort by level (errors first) then key
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Level != filtered[j].Level {
			return filtered[i].Level == vault.RuleLevelError
		}
		return filtered[i].Key < filtered[j].Key
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "LEVEL\tKEY\tRULE\tMESSAGE")
	for _, r := range filtered {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Level, r.Key, r.Rule, r.Message)
	}
	w.Flush()

	// Exit non-zero if any errors found
	for _, r := range filtered {
		if r.Level == vault.RuleLevelError {
			os.Exit(1)
		}
	}
	return nil
}
