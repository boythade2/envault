package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildRenameEnvVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vault.env")
	v, err := LoadOrCreate(p)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["APP_HOST"] = Entry{Value: "localhost"}
	v.Entries["APP_PORT"] = Entry{Value: "8080"}
	v.Entries["DB_HOST"] = Entry{Value: "db.local"}
	if err := v.Save(p, os.FileMode(0600)); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return p
}

func TestBulkRenameStripPrefix(t *testing.T) {
	p := buildRenameEnvVault(t)
	results, err := BulkRenameKeys(p, BulkRenameOptions{StripPrefix: "APP_"})
	if err != nil {
		t.Fatalf("BulkRenameKeys: %v", err)
	}
	renamed := 0
	for _, r := range results {
		if r.Renamed {
			renamed++
		}
	}
	if renamed != 2 {
		t.Errorf("expected 2 renamed, got %d", renamed)
	}
	v, _ := LoadOrCreate(p)
	if _, ok := v.Entries["HOST"]; !ok {
		t.Error("expected key HOST after strip")
	}
	if _, ok := v.Entries["PORT"]; !ok {
		t.Error("expected key PORT after strip")
	}
}

func TestBulkRenameAddPrefix(t *testing.T) {
	p := buildRenameEnvVault(t)
	_, err := BulkRenameKeys(p, BulkRenameOptions{AddPrefix: "NEW_"})
	if err != nil {
		t.Fatalf("BulkRenameKeys: %v", err)
	}
	v, _ := LoadOrCreate(p)
	if _, ok := v.Entries["NEW_APP_HOST"]; !ok {
		t.Error("expected key NEW_APP_HOST")
	}
}

func TestBulkRenameDryRun(t *testing.T) {
	p := buildRenameEnvVault(t)
	_, err := BulkRenameKeys(p, BulkRenameOptions{StripPrefix: "APP_", DryRun: true})
	if err != nil {
		t.Fatalf("BulkRenameKeys: %v", err)
	}
	v, _ := LoadOrCreate(p)
	if _, ok := v.Entries["APP_HOST"]; !ok {
		t.Error("dry run should not modify vault; APP_HOST missing")
	}
}

func TestBulkRenameConflictSkipped(t *testing.T) {
	p := buildRenameEnvVault(t)
	// Add a key that would conflict after strip
	v, _ := LoadOrCreate(p)
	v.Entries["HOST"] = Entry{Value: "conflict"}
	v.Save(p, os.FileMode(0600))

	results, err := BulkRenameKeys(p, BulkRenameOptions{StripPrefix: "APP_", Overwrite: false})
	if err != nil {
		t.Fatalf("BulkRenameKeys: %v", err)
	}
	for _, r := range results {
		if r.OldKey == "APP_HOST" && r.Renamed {
			t.Error("APP_HOST should have been skipped due to conflict")
		}
	}
}

func TestFormatBulkRenameResults(t *testing.T) {
	results := []EnvRenameResult{
		{OldKey: "APP_HOST", NewKey: "HOST", Renamed: true},
		{OldKey: "DB_HOST", NewKey: "DB_HOST", Renamed: false, Reason: "no change"},
	}
	out := FormatBulkRenameResults(results)
	if out == "" {
		t.Error("expected non-empty output")
	}
	if len(out) < 10 {
		t.Error("output too short")
	}
}
