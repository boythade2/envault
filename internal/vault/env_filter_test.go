package vault

import (
	"strings"
	"testing"
)

func buildFilterVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "DB_HOST", Value: "localhost", Tags: []string{"db"}},
		{Key: "DB_PORT", Value: "5432", Tags: []string{"db"}},
		{Key: "APP_SECRET", Value: "s3cr3t", Tags: []string{"app"}},
		{Key: "APP_DEBUG", Value: "true", Tags: []string{"app"}},
		{Key: "REDIS_URL", Value: "redis://localhost", Tags: []string{}},
	}
	return v
}

func TestFilterByKeyPrefix(t *testing.T) {
	v := buildFilterVault(t)
	results := FilterEntries(v, FilterOptions{KeyPrefix: "DB_"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !strings.HasPrefix(r.Key, "DB_") {
			t.Errorf("unexpected key %q", r.Key)
		}
	}
}

func TestFilterByKeySuffix(t *testing.T) {
	v := buildFilterVault(t)
	results := FilterEntries(v, FilterOptions{KeySuffix: "_URL"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "REDIS_URL" {
		t.Errorf("expected REDIS_URL, got %q", results[0].Key)
	}
}

func TestFilterByValueContains(t *testing.T) {
	v := buildFilterVault(t)
	results := FilterEntries(v, FilterOptions{ValueContains: "localhost"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestFilterInvertMatch(t *testing.T) {
	v := buildFilterVault(t)
	results := FilterEntries(v, FilterOptions{KeyPrefix: "DB_", InvertMatch: true})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestFilterByTag(t *testing.T) {
	v := buildFilterVault(t)
	results := FilterEntries(v, FilterOptions{TagFilter: []string{"app"}})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestFilterNoMatch(t *testing.T) {
	v := buildFilterVault(t)
	results := FilterEntries(v, FilterOptions{KeyPrefix: "NONEXISTENT_"})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
	out := FormatFilterResults(results)
	if !strings.Contains(out, "no entries matched") {
		t.Errorf("expected no-match message, got %q", out)
	}
}

func TestFormatFilterResultsHasHeaders(t *testing.T) {
	v := buildFilterVault(t)
	results := FilterEntries(v, FilterOptions{KeyPrefix: "DB_"})
	out := FormatFilterResults(results)
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "VALUE") {
		t.Errorf("expected headers in output, got %q", out)
	}
}
