package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildRenameVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	pass := "rename-secret"

	v, err := LoadOrCreate(path, pass)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["OLD_KEY"] = Entry{Value: "hello"}
	v.Entries["OTHER_KEY"] = Entry{Value: "world"}
	if err := v.Save(path, pass); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path, pass
}

func TestRenameEntrySuccess(t *testing.T) {
	path, pass := buildRenameVault(t)

	res, err := RenameEntry(path, pass, "OLD_KEY", "NEW_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OldKey != "OLD_KEY" || res.NewKey != "NEW_KEY" {
		t.Errorf("unexpected result: %+v", res)
	}

	v, _ := LoadOrCreate(path, pass)
	if _, exists := v.Entries["NEW_KEY"]; !exists {
		t.Error("NEW_KEY should exist after rename")
	}
	if _, exists := v.Entries["OLD_KEY"]; exists {
		t.Error("OLD_KEY should not exist after rename")
	}
	if v.Entries["NEW_KEY"].Value != "hello" {
		t.Errorf("value mismatch: got %q", v.Entries["NEW_KEY"].Value)
	}
}

func TestRenameEntryNotFound(t *testing.T) {
	path, pass := buildRenameVault(t)
	_, err := RenameEntry(path, pass, "MISSING", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameEntryConflict(t *testing.T) {
	path, pass := buildRenameVault(t)
	_, err := RenameEntry(path, pass, "OLD_KEY", "OTHER_KEY")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestRenameEntrySameKey(t *testing.T) {
	path, pass := buildRenameVault(t)
	_, err := RenameEntry(path, pass, "OLD_KEY", "OLD_KEY")
	if err == nil {
		t.Fatal("expected error when old and new keys are identical")
	}
}

func TestRenameEntryEmptyKey(t *testing.T) {
	path, pass := buildRenameVault(t)
	_, err := RenameEntry(path, pass, "", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for empty old key")
	}
}

func TestRenamePreservesOtherEntries(t *testing.T) {
	path, pass := buildRenameVault(t)
	_, err := RenameEntry(path, pass, "OLD_KEY", "RENAMED_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := LoadOrCreate(path, pass)
	if _, ok := v.Entries["OTHER_KEY"]; !ok {
		t.Error("OTHER_KEY should be preserved after rename")
	}
}

func TestRenameVaultFilePermissions(t *testing.T) {
	path, pass := buildRenameVault(t)
	_, err := RenameEntry(path, pass, "OLD_KEY", "RENAMED_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
