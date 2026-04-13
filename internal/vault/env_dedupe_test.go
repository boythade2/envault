package vault

import (
	"testing"
)

func buildDedupeVault() *Vault {
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "3000"},
		{Key: "HOST", Value: "remotehost"},
		{Key: "DEBUG", Value: "true"},
		{Key: "PORT", Value: "4000"},
		{Key: "PORT", Value: "5000"},
	}
	return v
}

func TestDedupeNoDuplicates(t *testing.T) {
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	results, err := DedupeEntries(v, DedupeOptions{Strategy: "first"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
	if len(v.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(v.Entries))
	}
}

func TestDedupeStrategyFirst(t *testing.T) {
	v := buildDedupeVault()
	results, err := DedupeEntries(v, DedupeOptions{Strategy: "first"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 dedupe results, got %d", len(results))
	}
	// HOST: first value should be kept
	for _, r := range results {
		if r.Key == "HOST" && r.Kept != "localhost" {
			t.Errorf("expected kept HOST=localhost, got %q", r.Kept)
		}
		if r.Key == "PORT" && r.Kept != "3000" {
			t.Errorf("expected kept PORT=3000, got %q", r.Kept)
		}
	}
	if len(v.Entries) != 3 {
		t.Errorf("expected 3 entries after dedupe, got %d", len(v.Entries))
	}
}

func TestDedupeStrategyLast(t *testing.T) {
	v := buildDedupeVault()
	_, err := DedupeEntries(v, DedupeOptions{Strategy: "last"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range v.Entries {
		if e.Key == "HOST" && e.Value != "remotehost" {
			t.Errorf("expected HOST=remotehost, got %q", e.Value)
		}
		if e.Key == "PORT" && e.Value != "5000" {
			t.Errorf("expected PORT=5000, got %q", e.Value)
		}
	}
}

func TestDedupeDryRunDoesNotModify(t *testing.T) {
	v := buildDedupeVault()
	before := len(v.Entries)
	_, err := DedupeEntries(v, DedupeOptions{Strategy: "first", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v.Entries) != before {
		t.Errorf("dry run should not modify entries: before=%d after=%d", before, len(v.Entries))
	}
}

func TestDedupeInvalidStrategy(t *testing.T) {
	v := buildDedupeVault()
	_, err := DedupeEntries(v, DedupeOptions{Strategy: "random"})
	if err == nil {
		t.Error("expected error for invalid strategy")
	}
}

func TestFormatDedupeResultsEmpty(t *testing.T) {
	out := FormatDedupeResults(nil)
	if out != "No duplicate keys found." {
		t.Errorf("unexpected output: %q", out)
	}
}
