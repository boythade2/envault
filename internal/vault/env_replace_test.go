package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildReplaceVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "URL", Value: "http://localhost:8080"},
		{Key: "API_URL", Value: "http://localhost:9090"},
		{Key: "NAME", Value: "envault"},
	}
	path := filepath.Join(dir, "vault.json")
	if err := v.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	return v, path
}

func TestReplaceAllOccurrences(t *testing.T) {
	v, path := buildReplaceVault(t)
	results, err := ReplaceValues(v, path, ReplaceOptions{Old: "localhost", New: "example.com", All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}
	if changed != 2 {
		t.Errorf("expected 2 changed, got %d", changed)
	}
}

func TestReplaceSelectedKeys(t *testing.T) {
	v, path := buildReplaceVault(t)
	results, err := ReplaceValues(v, path, ReplaceOptions{Keys: []string{"URL"}, Old: "localhost", New: "prod.io", All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "API_URL" && r.Changed {
			t.Error("API_URL should not have been changed")
		}
	}
}

func TestReplaceDryRunDoesNotWrite(t *testing.T) {
	v, path := buildReplaceVault(t)
	origStat, _ := os.Stat(path)
	_, err := ReplaceValues(v, path, ReplaceOptions{Old: "localhost", New: "remote", All: true, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	newStat, _ := os.Stat(path)
	if !origStat.ModTime().Equal(newStat.ModTime()) {
		t.Error("file should not have been modified during dry run")
	}
}

func TestReplaceEmptyOldReturnsError(t *testing.T) {
	v, path := buildReplaceVault(t)
	_, err := ReplaceValues(v, path, ReplaceOptions{Old: "", New: "x"})
	if err == nil {
		t.Error("expected error for empty old string")
	}
}

func TestReplaceNoMatch(t *testing.T) {
	v, path := buildReplaceVault(t)
	results, err := ReplaceValues(v, path, ReplaceOptions{Old: "notfound", New: "x", All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no changes, got change on %s", r.Key)
		}
	}
}
