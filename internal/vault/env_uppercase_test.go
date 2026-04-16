package vault

import (
	"path/filepath"
	"testing"
)

func buildUppercaseVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.env")
	pass := "testpass"
	v, err := LoadOrCreate(path, pass)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries = []Entry{
		{Key: "db_host", Value: "localhost"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "port", Value: "5432"},
	}
	data, err := v.serialize(pass)
	if err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if err := writeFile(path, data); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path, pass
}

func TestUppercaseAllKeys(t *testing.T) {
	path, pass := buildUppercaseVault(t)
	results, err := UppercaseKeys(path, pass, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var renamed, skipped int
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else {
			renamed++
		}
	}
	if renamed != 2 {
		t.Errorf("expected 2 renamed, got %d", renamed)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", skipped)
	}
}

func TestUppercaseSelectedKey(t *testing.T) {
	path, pass := buildUppercaseVault(t)
	results, err := UppercaseKeys(path, pass, []string{"port"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "PORT" {
		t.Errorf("expected KEY=PORT, got %s", results[0].Key)
	}
}

func TestUppercaseDryRunDoesNotWrite(t *testing.T) {
	path, pass := buildUppercaseVault(t)
	_, err := UppercaseKeys(path, pass, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, err := LoadOrCreate(path, pass)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	for _, e := range v.Entries {
		if e.Key == "DB_HOST" || e.Key == "PORT" {
			t.Errorf("dry run should not have modified vault, found %s", e.Key)
		}
	}
}

func TestFormatUppercaseResults(t *testing.T) {
	results := []UppercaseResult{
		{Key: "DB_HOST", OldKey: "db_host"},
		{Key: "API_KEY", OldKey: "API_KEY", Skipped: true, Reason: "already uppercase"},
	}
	out := FormatUppercaseResults(results)
	if out == "" {
		t.Error("expected non-empty output")
	}
	if len(out) < 10 {
		t.Errorf("output too short: %q", out)
	}
}
