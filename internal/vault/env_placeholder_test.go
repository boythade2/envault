package vault

import (
	"os"
	"testing"
)

func buildPlaceholderVault() *Vault {
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
		{Key: "DSN", Value: "postgres://{{HOST}}:{{PORT}}/db"},
		{Key: "GREETING", Value: "hello {{NAME}}"},
		{Key: "PLAIN", Value: "no placeholders here"},
	}
	return v
}

func TestResolvePlaceholdersInternalRef(t *testing.T) {
	v := buildPlaceholderVault()
	results, err := ResolvePlaceholders(v, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var dsn string
	for _, e := range v.Entries {
		if e.Key == "DSN" {
			dsn = e.Value
		}
	}
	if dsn != "postgres://localhost:5432/db" {
		t.Errorf("expected resolved DSN, got %q", dsn)
	}
	_ = results
}

func TestResolvePlaceholdersDryRunDoesNotModify(t *testing.T) {
	v := buildPlaceholderVault()
	_, err := ResolvePlaceholders(v, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range v.Entries {
		if e.Key == "DSN" && e.Value != "postgres://{{HOST}}:{{PORT}}/db" {
			t.Errorf("dry-run should not modify vault, got %q", e.Value)
		}
	}
}

func TestResolvePlaceholdersOSFallback(t *testing.T) {
	os.Setenv("NAME", "world")
	defer os.Unsetenv("NAME")
	v := buildPlaceholderVault()
	_, err := ResolvePlaceholders(v, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range v.Entries {
		if e.Key == "GREETING" && e.Value != "hello world" {
			t.Errorf("expected OS fallback resolution, got %q", e.Value)
		}
	}
}

func TestResolvePlaceholdersMissingReported(t *testing.T) {
	v := buildPlaceholderVault()
	results, err := ResolvePlaceholders(v, []string{"GREETING"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Missing) == 0 {
		t.Error("expected missing key NAME to be reported")
	}
}

func TestResolvePlaceholdersOnlyKeys(t *testing.T) {
	v := buildPlaceholderVault()
	results, err := ResolvePlaceholders(v, []string{"PLAIN"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "PLAIN" {
		t.Errorf("expected only PLAIN to be processed")
	}
}

func TestFormatPlaceholderResultsEmpty(t *testing.T) {
	out := FormatPlaceholderResults(nil)
	if out == "" {
		t.Error("expected non-empty output for nil results")
	}
}

func TestFormatPlaceholderResultsShowsChanges(t *testing.T) {
	v := buildPlaceholderVault()
	results, _ := ResolvePlaceholders(v, nil, true)
	out := FormatPlaceholderResults(results)
	if out == "" {
		t.Error("expected formatted output")
	}
}
