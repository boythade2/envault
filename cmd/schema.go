package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Manage and validate vault schema rules",
	}

	validateCmd := &cobra.Command{
		Use:   "validate <vault>",
		Short: "Validate a vault against its schema",
		Args:  cobra.ExactArgs(1),
		RunE:  runSchemaValidate,
	}

	addRuleCmd := &cobra.Command{
		Use:   "add-rule <vault> <key>",
		Short: "Add a rule to the vault schema",
		Args:  cobra.ExactArgs(2),
		RunE:  runSchemaAddRule,
	}
	addRuleCmd.Flags().Bool("required", false, "Mark key as required")
	addRuleCmd.Flags().String("pattern", "", "Regex pattern the value must match")
	addRuleCmd.Flags().String("desc", "", "Human-readable description")

	listRulesCmd := &cobra.Command{
		Use:   "list <vault>",
		Short: "List all schema rules for a vault",
		Args:  cobra.ExactArgs(1),
		RunE:  runSchemaList,
	}

	schemaCmd.AddCommand(validateCmd, addRuleCmd, listRulesCmd)
	rootCmd.AddCommand(schemaCmd)
}

func runSchemaValidate(cmd *cobra.Command, args []string) error {
	v, err := vault.LoadOrCreate(args[0])
	if err != nil {
		return fmt.Errorf("load vault: %w", err)
	}
	s, err := vault.LoadSchema(args[0])
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}
	violations := vault.ValidateSchema(v, s)
	if len(violations) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ vault is valid")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVIOLATION")
	for _, viol := range violations {
		fmt.Fprintf(w, "%s\t%s\n", viol.Key, viol.Message)
	}
	w.Flush()
	return fmt.Errorf("%d schema violation(s) found", len(violations))
}

func runSchemaAddRule(cmd *cobra.Command, args []string) error {
	vaultPath, key := args[0], args[1]
	s, err := vault.LoadSchema(vaultPath)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}
	required, _ := cmd.Flags().GetBool("required")
	pattern, _ := cmd.Flags().GetString("pattern")
	desc, _ := cmd.Flags().GetString("desc")
	s.Rules = append(s.Rules, vault.SchemaRule{
		Key:      key,
		Required: required,
		Pattern:  pattern,
		Desc:     desc,
	})
	if err := vault.SaveSchema(vaultPath, s); err != nil {
		return fmt.Errorf("save schema: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "rule added for key %q\n", key)
	return nil
}

func runSchemaList(cmd *cobra.Command, args []string) error {
	s, err := vault.LoadSchema(args[0])
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}
	if len(s.Rules) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no schema rules defined")
		return nil
	}
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tREQUIRED\tPATTERN\tDESCRIPTION")
	for _, r := range s.Rules {
		fmt.Fprintf(w, "%s\t%v\t%s\t%s\n", r.Key, r.Required, r.Pattern, r.Desc)
	}
	w.Flush()
	_ = os.Stdout
	return nil
}
