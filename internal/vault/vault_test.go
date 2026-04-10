package vault_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envault/internal/vault"
)

func TestLoadOrCreateNewVault(t *testing.T) {
	dir := t.TempDir()
	v, err := vault.LoadOrCreate(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Version != 1 {
		t.Errorf("expected version 1, got %d", v.Version)
	}
	if len(v.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(v.Entries))
	}
}

func TestSaveAndReload(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.LoadOrCreate(dir)
	v.AddEntry("dev", ".env.dev.enc", true)

	if err := v.Save(dir); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	v2, err := vault.LoadOrCreate(dir)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	entry, ok := v2.Entries["dev"]
	if !ok {
		t.Fatal("entry 'dev' not found after reload")
	}
	if entry.FilePath != ".env.dev.enc" || !entry.Encrypted {
		t.Errorf("unexpected entry values: %+v", entry)
	}
}

func TestAddEntryUpdatesTimestamp(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.LoadOrCreate(dir)
	v.AddEntry("prod", ".env.prod.enc", true)

	before := v.Entries["prod"].UpdatedAt
	time.Sleep(2 * time.Millisecond)
	v.AddEntry("prod", ".env.prod.enc", false)

	if !v.Entries["prod"].UpdatedAt.After(before) {
		t.Error("UpdatedAt should be refreshed on second AddEntry")
	}
}

func TestRemoveEntry(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.LoadOrCreate(dir)
	v.AddEntry("staging", ".env.staging.enc", true)

	if removed := v.RemoveEntry("staging"); !removed {
		t.Error("expected RemoveEntry to return true")
	}
	if _, ok := v.Entries["staging"]; ok {
		t.Error("entry should have been deleted")
	}
	if removed := v.RemoveEntry("nonexistent"); removed {
		t.Error("expected RemoveEntry to return false for missing entry")
	}
}

func TestVaultFilePermissions(t *testing.T) {
	dir := t.TempDir()
	v, _ := vault.LoadOrCreate(dir)
	v.AddEntry("ci", ".env.ci.enc", true)
	_ = v.Save(dir)

	info, err := os.Stat(filepath.Join(dir, ".envault.json"))
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %o", info.Mode().Perm())
	}
}
