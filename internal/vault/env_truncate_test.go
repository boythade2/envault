package vault

import (
	"strings"
	"testing"
)

func buildTruncateVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "SHORT", Value: "hi"},
		{Key: "LONG_KEY", Value: "abcdefghijklmnopqrstuvwxyz"},
		{Key: "EXACT", Value: "12345"},
	}
	return v
}

func TestTruncateZeroMaxLenReturnsError(t *testing.T) {
	v := buildTruncateVault(t)
	_, err := TruncateEntries(v, TruncateOptions{MaxLen: 0})
	if err == nil {
		t.Fatal("expected error for MaxLen=0")
	}
}

func TestTruncateShortValuesSkipped(t *testing.T) {
	v := buildTruncateVault(t)
	results, err := TruncateEntries(v, TruncateOptions{MaxLen: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "SHORT" && !r.Skipped {
			t.Errorf("expected SHORT to be skipped")
		}
	}
}

func TestTruncateLongValueIsCut(t *testing.T) {
	v := buildTruncateVault(t)
	results, err := TruncateEntries(v, TruncateOptions{MaxLen: 10, Suffix: "..."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "LONG_KEY" {
			if len(r.Truncated) != 10 {
				t.Errorf("expected truncated length 10, got %d", len(r.Truncated))
			}
			if !strings.HasSuffix(r.Truncated, "...") {
				t.Errorf("expected suffix '...', got %q", r.Truncated)
			}
		}
	}
}

func TestTruncateDryRunDoesNotModify(t *testing.T) {
	v := buildTruncateVault(t)
	original := v.Entries[1].Value
	_, err := TruncateEntries(v, TruncateOptions{MaxLen: 5, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Entries[1].Value != original {
		t.Errorf("dry run should not modify vault; got %q", v.Entries[1].Value)
	}
}

func TestTruncateSelectedKeysOnly(t *testing.T) {
	v := buildTruncateVault(t)
	results, err := TruncateEntries(v, TruncateOptions{
		MaxLen: 5,
		Keys:   []string{"LONG_KEY"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestFormatTruncateResultsNoneChanged(t *testing.T) {
	v := buildTruncateVault(t)
	results, _ := TruncateEntries(v, TruncateOptions{MaxLen: 100})
	out := FormatTruncateResults(results)
	if !strings.Contains(out, "No values") {
		t.Errorf("expected no-op message, got %q", out)
	}
}

func TestFormatTruncateResultsShowsCount(t *testing.T) {
	v := buildTruncateVault(t)
	results, _ := TruncateEntries(v, TruncateOptions{MaxLen: 5, DryRun: true})
	out := FormatTruncateResults(results)
	if !strings.Contains(out, "truncated") {
		t.Errorf("expected truncated summary, got %q", out)
	}
}
