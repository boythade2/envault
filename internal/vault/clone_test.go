package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildCloneVault(t *testing.T, dir, passphrase string) string {
	t.Helper()
	path := filepath.Join(dir, "src.vault")
	v, err := LoadOrCreate(path, passphrase)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	if err := v.AddEncrypted("API_KEY", "secret123", passphrase); err != nil {
		t.Fatalf("AddEncrypted: %v", err)
	}
	if err := v.AddEncrypted("DB_URL", "postgres://localhost/db", passphrase); err != nil {
		t.Fatalf("AddEncrypted: %v", err)
	}
	if err := v.Save(path, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

func TestCloneVaultCreatesDestination(t *testing.T) {
	dir := t.TempDir()
	src := buildCloneVault(t, dir, "old-pass")
	dst := filepath.Join(dir, "dst.vault")

	result, err := CloneVault(src, dst, "old-pass", "new-pass")
	if err != nil {
		t.Fatalf("CloneVault: %v", err)
	}
	if result.EntriesCount != 2 {
		t.Errorf("expected 2 entries, got %d", result.EntriesCount)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("destination file not created: %v", err)
	}
}

func TestCloneVaultPreservesValues(t *testing.T) {
	dir := t.TempDir()
	src := buildCloneVault(t, dir, "pass1")
	dst := filepath.Join(dir, "dst.vault")

	if _, err := CloneVault(src, dst, "pass1", "pass2"); err != nil {
		t.Fatalf("CloneVault: %v", err)
	}

	v, err := LoadOrCreate(dst, "pass2")
	if err != nil {
		t.Fatalf("LoadOrCreate dst: %v", err)
	}
	val, err := v.GetDecrypted("API_KEY", "pass2")
	if err != nil {
		t.Fatalf("GetDecrypted: %v", err)
	}
	if val != "secret123" {
		t.Errorf("expected 'secret123', got %q", val)
	}
}

func TestCloneVaultRejectsExistingDest(t *testing.T) {
	dir := t.TempDir()
	src := buildCloneVault(t, dir, "pass")
	dst := filepath.Join(dir, "dst.vault")
	// Pre-create destination.
	if err := os.WriteFile(dst, []byte{}, 0600); err != nil {
		t.Fatal(err)
	}
	_, err := CloneVault(src, dst, "pass", "pass")
	if err == nil {
		t.Error("expected error when destination already exists")
	}
}

func TestWriteCloneMeta(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "dst.vault")
	result := CloneResult{Source: "src.vault", Destination: dst, EntriesCount: 3}
	if err := WriteCloneMeta(result); err != nil {
		t.Fatalf("WriteCloneMeta: %v", err)
	}
	info, err := os.Stat(cloneMetaPath(dst))
	if err != nil {
		t.Fatalf("meta file not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
