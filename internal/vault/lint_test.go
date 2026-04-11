package vault

import (
	"testing"
)

func buildLintVault(entries []Entry) *Vault {
	return &Vault{Entries: entries}
}

func TestLintCleanVault(t *testing.T) {
	v := buildLintVault([]Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "API_KEY", Value: "secret123"},
	})
	result := LintVault(v)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestLintLowercaseKeyWarn(t *testing.T) {
	v := buildLintVault([]Entry{
		{Key: "database_url", Value: "postgres://localhost/db"},
	})
	result := LintVault(v)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", result.Issues[0].Severity)
	}
}

func TestLintEmptyValueWarn(t *testing.T) {
	v := buildLintVault([]Entry{
		{Key: "MY_VAR", Value: ""},
	})
	result := LintVault(v)
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "MY_VAR" && issue.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected empty-value warning for MY_VAR")
	}
}

func TestLintKeyWithSpacesError(t *testing.T) {
	v := buildLintVault([]Entry{
		{Key: "MY VAR", Value: "value"},
	})
	result := LintVault(v)
	if !result.HasErrors() {
		t.Error("expected at least one error for key with spaces")
	}
}

func TestLintKeyStartsWithDigitWarn(t *testing.T) {
	v := buildLintVault([]Entry{
		{Key: "1_VAR", Value: "value"},
	})
	result := LintVault(v)
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "1_VAR" && issue.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn for key starting with digit")
	}
}

func TestLintIssueString(t *testing.T) {
	issue := LintIssue{Key: "FOO", Severity: "warn", Message: "some issue"}
	s := issue.String()
	if s != "[WARN] FOO: some issue" {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestHasErrorsFalseForWarnsOnly(t *testing.T) {
	result := LintResult{
		Issues: []LintIssue{
			{Key: "X", Severity: "warn", Message: "lowercase"},
		},
	}
	if result.HasErrors() {
		t.Error("expected HasErrors to be false when only warnings present")
	}
}
