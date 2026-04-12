package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildAliasVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadAliasesNoFile(t *testing.T) {
	vp := buildAliasVault(t)
	store, err := LoadAliases(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Aliases) != 0 {
		t.Errorf("expected empty aliases, got %d", len(store.Aliases))
	}
}

func TestAddAndResolveAlias(t *testing.T) {
	vp := buildAliasVault(t)
	if err := AddAlias(vp, "db", "DATABASE_URL"); err != nil {
		t.Fatalf("AddAlias failed: %v", err)
	}
	resolved, err := ResolveAlias(vp, "db")
	if err != nil {
		t.Fatalf("ResolveAlias failed: %v", err)
	}
	if resolved != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", resolved)
	}
}

func TestResolveAliasPassthrough(t *testing.T) {
	vp := buildAliasVault(t)
	resolved, err := ResolveAlias(vp, "SOME_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "SOME_KEY" {
		t.Errorf("expected passthrough, got %s", resolved)
	}
}

func TestAddDuplicateAliasReturnsError(t *testing.T) {
	vp := buildAliasVault(t)
	_ = AddAlias(vp, "db", "DATABASE_URL")
	err := AddAlias(vp, "db", "OTHER_KEY")
	if err == nil {
		t.Fatal("expected error for duplicate alias")
	}
}

func TestRemoveAlias(t *testing.T) {
	vp := buildAliasVault(t)
	_ = AddAlias(vp, "db", "DATABASE_URL")
	if err := RemoveAlias(vp, "db"); err != nil {
		t.Fatalf("RemoveAlias failed: %v", err)
	}
	store, _ := LoadAliases(vp)
	if _, ok := store.Aliases["db"]; ok {
		t.Error("alias should have been removed")
	}
}

func TestRemoveNonExistentAlias(t *testing.T) {
	vp := buildAliasVault(t)
	err := RemoveAlias(vp, "ghost")
	if err == nil {
		t.Fatal("expected error removing non-existent alias")
	}
}

func TestAliasFilePermissions(t *testing.T) {
	vp := buildAliasVault(t)
	_ = AddAlias(vp, "s", "SECRET_KEY")
	p := aliasPath(vp)
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
