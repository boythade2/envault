package vault

import (
	"strings"
	"testing"
)

func buildPrefixVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["DB_HOST"] = Entry{Value: "localhost"}
	v.Entries["DB_PORT"] = Entry{Value: "5432"}
	v.Entries["APP_NAME"] = Entry{Value: "envault"}
	return v
}

func TestAddKeyPrefixAll(t *testing.T) {
	v := buildPrefixVault(t)
	results, err := AddKeyPrefix(v, nil, PrefixOptions{Prefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Skipped {
			t.Errorf("unexpected skip for key %s", r.OldKey)
		}
		if !strings.HasPrefix(r.NewKey, "PROD_") {
			t.Errorf("new key %s missing prefix", r.NewKey)
		}
		if _, ok := v.Entries[r.NewKey]; !ok {
			t.Errorf("new key %s not found in vault", r.NewKey)
		}
		if _, ok := v.Entries[r.OldKey]; ok {
			t.Errorf("old key %s should have been removed", r.OldKey)
		}
	}
}

func TestAddKeyPrefixDryRun(t *testing.T) {
	v := buildPrefixVault(t)
	_, err := AddKeyPrefix(v, nil, PrefixOptions{Prefix: "PROD_", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := v.Entries["DB_HOST"]; !ok {
		t.Error("dry-run should not modify vault: DB_HOST missing")
	}
}

func TestAddKeyPrefixEmptyPrefixReturnsError(t *testing.T) {
	v := buildPrefixVault(t)
	_, err := AddKeyPrefix(v, nil, PrefixOptions{Prefix: ""})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestAddKeyPrefixSkipsConflictWithoutOverwrite(t *testing.T) {
	v := buildPrefixVault(t)
	v.Entries["PROD_DB_HOST"] = Entry{Value: "other"}
	results, err := AddKeyPrefix(v, []string{"DB_HOST"}, PrefixOptions{Prefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Error("expected DB_HOST to be skipped due to conflict")
	}
}

func TestRemoveKeyPrefixAll(t *testing.T) {
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["PROD_DB_HOST"] = Entry{Value: "localhost"}
	v.Entries["PROD_DB_PORT"] = Entry{Value: "5432"}
	v.Entries["APP_NAME"] = Entry{Value: "envault"}
	results, err := RemoveKeyPrefix(v, nil, PrefixOptions{Prefix: "PROD_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	skipped := 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		}
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped (APP_NAME), got %d", skipped)
	}
	if _, ok := v.Entries["DB_HOST"]; !ok {
		t.Error("expected DB_HOST after prefix removal")
	}
}

func TestRemoveKeyPrefixYieldsEmptyKeySkipped(t *testing.T) {
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["PROD_"] = Entry{Value: "x"}
	results, _ := RemoveKeyPrefix(v, nil, PrefixOptions{Prefix: "PROD_"})
	if len(results) != 1 || !results[0].Skipped {
		t.Error("expected skip when stripping prefix yields empty key")
	}
}

func TestFormatPrefixResultsEmpty(t *testing.T) {
	out := FormatPrefixResults(nil)
	if out != "no keys processed" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatPrefixResultsContainsOK(t *testing.T) {
	results := []PrefixResult{{OldKey: "DB_HOST", NewKey: "PROD_DB_HOST"}}
	out := FormatPrefixResults(results)
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK in output, got: %s", out)
	}
}
