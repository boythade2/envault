package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// RuleLevel indicates severity of a lint rule violation.
type RuleLevel string

const (
	RuleLevelError RuleLevel = "error"
	RuleLevelWarn  RuleLevel = "warn"
)

// LintRule defines a single named lint rule applied to vault entries.
type LintRule struct {
	Name    string
	Level   RuleLevel
	Message string
	Check   func(key, value string) bool
}

// LintRuleResult is a single finding from running a lint rule.
type LintRuleResult struct {
	Key     string
	Rule    string
	Level   RuleLevel
	Message string
}

func (r LintRuleResult) String() string {
	return fmt.Sprintf("[%s] %s — %s (%s)", strings.ToUpper(string(r.Level)), r.Key, r.Message, r.Rule)
}

// DefaultLintRules returns the built-in set of lint rules.
func DefaultLintRules() []LintRule {
	return []LintRule{
		{
			Name:    "no-lowercase-key",
			Level:   RuleLevelWarn,
			Message: "key contains lowercase letters; prefer UPPER_SNAKE_CASE",
			Check: func(key, _ string) bool {
				return key != strings.ToUpper(key)
			},
		},
		{
			Name:    "no-empty-value",
			Level:   RuleLevelWarn,
			Message: "value is empty",
			Check: func(_, value string) bool {
				return strings.TrimSpace(value) == ""
			},
		},
		{
			Name:    "no-spaces-in-key",
			Level:   RuleLevelError,
			Message: "key contains spaces",
			Check: func(key, _ string) bool {
				return strings.Contains(key, " ")
			},
		},
		{
			Name:    "no-special-chars-in-key",
			Level:   RuleLevelError,
			Message: "key contains invalid characters (only A-Z, 0-9, _ allowed)",
			Check: func(key, _ string) bool {
				matched, _ := regexp.MatchString(`[^A-Za-z0-9_]`, key)
				return matched
			},
		},
		{
			Name:    "no-numeric-prefix",
			Level:   RuleLevelWarn,
			Message: "key starts with a digit",
			Check: func(key, _ string) bool {
				return len(key) > 0 && key[0] >= '0' && key[0] <= '9'
			},
		},n}

// RunLintRules applies a set of rules to a vault and returns all findings.
func RunLintRules(v *Vault, rules []LintRule) []LintRuleResult {
	var results []LintRuleResult
	for key, entry := range v.Entries {
		for _, rule := range rules {
			if rule.Check(key, entry.Value) {
				results = append(results, LintRuleResult{
					Key:     key,
					Rule:    rule.Name,
					Level:   rule.Level,
					Message: rule.Message,
				})
			}
		}
	}
	return results
}
