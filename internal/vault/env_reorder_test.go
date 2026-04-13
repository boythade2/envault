package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildReorderVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vault.json")
	v, err := LoadOrCreate(p)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	for _, kv := range [][]string{{"ALPHA", "1"}, {"BETA", "2"}, {"GAMMA", "3"}, {"DELTA", "4"}} {
		v.Entries = append(v.Entries, Entry{Key: kv[0], Value: kv[1]})
	}
	data, _ := v.marshal()
	os.WriteFile(p, data, 0600)
	return p
}

func TestReorderMovesKeysToFront(t *testing.T) {
	p := buildReorderVault(t)
	results, err := ReorderEntries(p, []string{"GAMMA", "ALPHA"}, false)
	if err != nil {
		t.Fatalf("ReorderEntries: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	v, _ := LoadOrCreate(p)
	if v.Entries[0].Key != "GAMMA" || v.Entries[1].Key != "ALPHA" {
		t.Errorf("unexpected order: %v", v.Entries)
	}
	// Remaining keys preserve original relative order.
	if v.Entries[2].Key != "BETA" || v.Entries[3].Key != "DELTA" {
		t.Errorf("tail order wrong: %v", v.Entries[2:])
	}
}

func TestReorderSkipsMissingKey(t *testing.T) {
	p := buildReorderVault(t)
	results, err := ReorderEntries(p, []string{"MISSING", "BETA"}, false)
	if err != nil {
		t.Fatalf("ReorderEntries: %v", err)
	}
	var skipped bool
	for _, r := range results {
		if r.Key == "MISSING" && r.Skipped {
			skipped = true
		}
	}
	if !skipped {
		t.Error("expected MISSING to be skipped")
	}
}

func TestReorderDryRunDoesNotWrite(t *testing.T) {
	p := buildReorderVault(t)
	original, _ := os.ReadFile(p)
	_, err := ReorderEntries(p, []string{"DELTA", "GAMMA"}, true)
	if err != nil {
		t.Fatalf("ReorderEntries dry-run: %v", err)
	}
	after, _ := os.ReadFile(p)
	if string(original) != string(after) {
		t.Error("dry-run must not modify the vault file")
	}
}

func TestReorderSkipsDuplicateKeyInList(t *testing.T) {
	p := buildReorderVault(t)
	results, err := ReorderEntries(p, []string{"ALPHA", "ALPHA"}, false)
	if err != nil {
		t.Fatalf("ReorderEntries: %v", err)
	}
	var dupSkipped bool
	for _, r := range results {
		if r.Key == "ALPHA" && r.Skipped && r.Reason == "duplicate in key list" {
			dupSkipped = true
		}
	}
	if !dupSkipped {
		t.Error("expected second ALPHA to be reported as duplicate")
	}
}

func TestFormatReorderResultsSummary(t *testing.T) {
	results := []ReorderResult{
		{Key: "FOO", Position: 1},
		{Key: "BAR", Skipped: true, Reason: "key not found"},
	}
	out := FormatReorderResults(results, false)
	if out == "" {
		t.Error("expected non-empty format output")
	}
	if out == "no entries reordered" {
		t.Error("expected formatted results, not empty message")
	}
}
