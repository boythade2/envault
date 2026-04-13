package vault

import (
	"testing"
	"time"
)

func buildLintRulesVault(entries map[string]string) *Vault {
	v := &Vault{Entries: make(map[string]Entry)}
	for k, val := range entries {
		v.Entries[k] = Entry{Value: val, UpdatedAt: time.Now()}
	}
	return v
}

func TestDefaultLintRulesNotEmpty(t *testing.T) {
	rules := DefaultLintRules()
	if len(rules) == 0 {
		t.Fatal("expected at least one default lint rule")
	}
}

func TestLintRulesCleanVault(t *testing.T) {
	v := buildLintRulesVault(map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	})
	results := RunLintRules(v, DefaultLintRules())
	if len(results) != 0 {
		t.Fatalf("expected no lint findings, got %d: %v", len(results), results)
	}
}

func TestLintRulesLowercaseKey(t *testing.T) {
	v := buildLintRulesVault(map[string]string{"db_host": "localhost"})
	results := RunLintRules(v, DefaultLintRules())
	found := false
	for _, r := range results {
		if r.Rule == "no-lowercase-key" && r.Key == "db_host" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-lowercase-key finding for 'db_host'")
	}
}

func TestLintRulesEmptyValue(t *testing.T) {
	v := buildLintRulesVault(map[string]string{"EMPTY_VAR": ""})
	results := RunLintRules(v, DefaultLintRules())
	found := false
	for _, r := range results {
		if r.Rule == "no-empty-value" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-empty-value finding")
	}
}

func TestLintRulesSpaceInKey(t *testing.T) {
	v := buildLintRulesVault(map[string]string{"BAD KEY": "value"})
	results := RunLintRules(v, DefaultLintRules())
	found := false
	for _, r := range results {
		if r.Rule == "no-spaces-in-key" && r.Level == RuleLevelError {
			found = true
		}
	}
	if !found {
		t.Error("expected no-spaces-in-key error finding")
	}
}

func TestLintRulesNumericPrefix(t *testing.T) {
	v := buildLintRulesVault(map[string]string{"1INVALID": "value"})
	results := RunLintRules(v, DefaultLintRules())
	found := false
	for _, r := range results {
		if r.Rule == "no-numeric-prefix" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-numeric-prefix finding")
	}
}

func TestLintRuleResultString(t *testing.T) {
	r := LintRuleResult{
		Key:     "bad key",
		Rule:    "no-spaces-in-key",
		Level:   RuleLevelError,
		Message: "key contains spaces",
	}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
