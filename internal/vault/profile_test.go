package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildProfileDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func TestLoadProfilesNoFile(t *testing.T) {
	dir := buildProfileDir(t)
	store, err := LoadProfiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Profiles) != 0 {
		t.Errorf("expected empty profiles, got %d", len(store.Profiles))
	}
}

func TestAddAndGetProfile(t *testing.T) {
	dir := buildProfileDir(t)
	store, _ := LoadProfiles(dir)

	if err := store.AddProfile("dev", "dev.vault"); err != nil {
		t.Fatalf("AddProfile: %v", err)
	}
	p, err := store.GetProfile("dev")
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if p.VaultFile != "dev.vault" {
		t.Errorf("expected vault file dev.vault, got %s", p.VaultFile)
	}
}

func TestAddDuplicateProfileReturnsError(t *testing.T) {
	dir := buildProfileDir(t)
	store, _ := LoadProfiles(dir)
	_ = store.AddProfile("staging", "staging.vault")
	err := store.AddProfile("staging", "other.vault")
	if err == nil {
		t.Fatal("expected error for duplicate profile")
	}
}

func TestRemoveProfile(t *testing.T) {
	dir := buildProfileDir(t)
	store, _ := LoadProfiles(dir)
	_ = store.AddProfile("prod", "prod.vault")
	if err := store.RemoveProfile("prod"); err != nil {
		t.Fatalf("RemoveProfile: %v", err)
	}
	if _, err := store.GetProfile("prod"); err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemoveNonExistentProfileReturnsError(t *testing.T) {
	dir := buildProfileDir(t)
	store, _ := LoadProfiles(dir)
	if err := store.RemoveProfile("ghost"); err == nil {
		t.Fatal("expected error removing non-existent profile")
	}
}

func TestSaveAndReloadProfiles(t *testing.T) {
	dir := buildProfileDir(t)
	store, _ := LoadProfiles(dir)
	_ = store.AddProfile("dev", "dev.vault")
	_ = store.AddProfile("prod", "prod.vault")
	if err := store.Save(dir); err != nil {
		t.Fatalf("Save: %v", err)
	}

	reloaded, err := LoadProfiles(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(reloaded.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(reloaded.Profiles))
	}
}

func TestProfileFilePermissions(t *testing.T) {
	dir := buildProfileDir(t)
	store, _ := LoadProfiles(dir)
	_ = store.AddProfile("dev", "dev.vault")
	_ = store.Save(dir)

	info, err := os.Stat(filepath.Join(dir, ".envault_profiles.json"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
