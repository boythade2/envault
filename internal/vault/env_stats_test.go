package vault

import (
	"strings"
	"testing"
	"time"
)

func buildStatsVault(t *testing.T) *Vault {
	t.Helper()
	now := time.Now()
	return &Vault{
		Entries: []Entry{
			{Key: "APP_HOST", Value: "localhost", UpdatedAt: now.Add(-48 * time.Hour)},
			{Key: "APP_PORT", Value: "8080", UpdatedAt: now.Add(-24 * time.Hour)},
			{Key: "APP_SECRET", Value: "s3cr3t", UpdatedAt: now},
			{Key: "DB_HOST", Value: "localhost", UpdatedAt: now.Add(-12 * time.Hour)},
			{Key: "DB_PASS", Value: "", UpdatedAt: now.Add(-6 * time.Hour)},
			{Key: "STANDALONE", Value: "only", UpdatedAt: now.Add(-1 * time.Hour)},
		},
	}
}

func TestComputeStatsTotalKeys(t *testing.T) {
	v := buildStatsVault(t)
	s := ComputeStats(v)
	if s.TotalKeys != 6 {
		t.Fatalf("expected 6 total keys, got %d", s.TotalKeys)
	}
}

func TestComputeStatsEmptyValues(t *testing.T) {
	v := buildStatsVault(t)
	s := ComputeStats(v)
	if s.EmptyValues != 1 {
		t.Fatalf("expected 1 empty value, got %d", s.EmptyValues)
	}
}

func TestComputeStatsUniqueValues(t *testing.T) {
	v := buildStatsVault(t)
	s := ComputeStats(v)
	// "localhost" appears twice, so unique count = total distinct values appearing once
	// values: localhost(x2), 8080, s3cr3t, ""(x1), only => unique = 4
	if s.UniqueValues != 4 {
		t.Fatalf("expected 4 unique values, got %d", s.UniqueValues)
	}
}

func TestComputeStatsTopPrefixes(t *testing.T) {
	v := buildStatsVault(t)
	s := ComputeStats(v)
	if len(s.TopPrefixes) == 0 {
		t.Fatal("expected at least one prefix")
	}
	// APP_ has 3 keys, DB_ has 2 — APP should be first
	if s.TopPrefixes[0].Prefix != "APP" {
		t.Fatalf("expected top prefix APP, got %s", s.TopPrefixes[0].Prefix)
	}
	if s.TopPrefixes[0].Count != 3 {
		t.Fatalf("expected APP prefix count 3, got %d", s.TopPrefixes[0].Count)
	}
}

func TestComputeStatsOldestNewest(t *testing.T) {
	v := buildStatsVault(t)
	s := ComputeStats(v)
	if s.OldestUpdated.IsZero() || s.NewestUpdated.IsZero() {
		t.Fatal("expected non-zero timestamps")
	}
	if !s.OldestUpdated.Before(s.NewestUpdated) {
		t.Fatal("oldest should be before newest")
	}
}

func TestComputeStatsEmptyVault(t *testing.T) {
	v := &Vault{}
	s := ComputeStats(v)
	if s.TotalKeys != 0 {
		t.Fatalf("expected 0 keys for empty vault, got %d", s.TotalKeys)
	}
}

func TestFormatStatsEmptyVault(t *testing.T) {
	v := &Vault{}
	out := FormatStats(ComputeStats(v))
	if !strings.Contains(out, "empty") {
		t.Fatalf("expected empty message, got: %s", out)
	}
}

func TestFormatStatsContainsFields(t *testing.T) {
	v := buildStatsVault(t)
	out := FormatStats(ComputeStats(v))
	for _, want := range []string{"Total keys", "Empty values", "Unique values", "Oldest update", "Newest update", "Top prefixes"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q", want)
		}
	}
}
