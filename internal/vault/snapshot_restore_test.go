package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildRestoreVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	v, err := LoadOrCreate(vaultPath, "pass")
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["KEY1"] = Entry{Value: "val1"}
	v.Entries["KEY2"] = Entry{Value: "val2"}
	if err := v.Save("pass"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return v, vaultPath
}

func TestRestoreSnapshotNotFound(t *testing.T) {
	_, vaultPath := buildRestoreVault(t)
	_, err := RestoreSnapshot(vaultPath, "ghost", "pass")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestRestoreSnapshotPreservesEntries(t *testing.T) {
	v, vaultPath := buildRestoreVault(t)

	if err := SaveSnapshot(vaultPath, "before-change", v.Entries); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	// Mutate the vault.
	v.Entries["KEY1"] = Entry{Value: "mutated"}
	delete(v.Entries, "KEY2")
	if err := v.Save("pass"); err != nil {
		t.Fatalf("Save after mutation: %v", err)
	}

	// Restore from snapshot.
	restored, err := RestoreSnapshot(vaultPath, "before-change", "pass")
	if err != nil {
		t.Fatalf("RestoreSnapshot: %v", err)
	}

	if restored.Entries["KEY1"].Value != "val1" {
		t.Errorf("KEY1 = %q, want %q", restored.Entries["KEY1"].Value, "val1")
	}
	if _, ok := restored.Entries["KEY2"]; !ok {
		t.Error("KEY2 missing after restore")
	}
}

func TestRestoreSnapshotOverwritesVaultFile(t *testing.T) {
	v, vaultPath := buildRestoreVault(t)

	if err := SaveSnapshot(vaultPath, "snap1", v.Entries); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	statBefore, _ := os.Stat(vaultPath)

	_, err := RestoreSnapshot(vaultPath, "snap1", "pass")
	if err != nil {
		t.Fatalf("RestoreSnapshot: %v", err)
	}

	statAfter, _ := os.Stat(vaultPath)
	if !statAfter.ModTime().After(statBefore.ModTime()) {
		t.Error("vault file mod time should have updated after restore")
	}
}
