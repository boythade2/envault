package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildArchiveVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	vaultFile := filepath.Join(dir, "test.vault")
	passphrase := "archivepass"
	v, err := LoadOrCreate(vaultFile, passphrase)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["OLD_KEY"] = Entry{Value: "old_value"}
	v.Entries["KEEP_KEY"] = Entry{Value: "keep_value"}
	if err := v.Save(vaultFile, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return vaultFile, passphrase
}

func TestLoadArchiveNoFile(t *testing.T) {
	dir := t.TempDir()
	store, err := LoadArchive(filepath.Join(dir, "missing.vault"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store.Entries) != 0 {
		t.Fatalf("expected empty archive")
	}
}

func TestArchiveEntriesRemovesFromVault(t *testing.T) {
	vaultFile, pass := buildArchiveVault(t)
	archived, err := ArchiveEntries(vaultFile, pass, []string{"OLD_KEY"}, "deprecated", false)
	if err != nil {
		t.Fatalf("ArchiveEntries: %v", err)
	}
	if len(archived) != 1 || archived[0] != "OLD_KEY" {
		t.Fatalf("unexpected archived keys: %v", archived)
	}
	v, err := LoadOrCreate(vaultFile, pass)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if _, ok := v.Entries["OLD_KEY"]; ok {
		t.Fatal("OLD_KEY should have been removed from vault")
	}
	if _, ok := v.Entries["KEEP_KEY"]; !ok {
		t.Fatal("KEEP_KEY should remain in vault")
	}
}

func TestArchiveEntriesPersistedToFile(t *testing.T) {
	vaultFile, pass := buildArchiveVault(t)
	_, err := ArchiveEntries(vaultFile, pass, []string{"OLD_KEY"}, "cleanup", false)
	if err != nil {
		t.Fatalf("ArchiveEntries: %v", err)
	}
	store, err := LoadArchive(vaultFile)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(store.Entries) != 1 {
		t.Fatalf("expected 1 archived entry, got %d", len(store.Entries))
	}
	if store.Entries[0].Key != "OLD_KEY" {
		t.Fatalf("expected OLD_KEY, got %s", store.Entries[0].Key)
	}
	if store.Entries[0].Reason != "cleanup" {
		t.Fatalf("expected reason 'cleanup', got %s", store.Entries[0].Reason)
	}
	if store.Entries[0].ArchivedAt.IsZero() {
		t.Fatal("ArchivedAt should be set")
	}
}

func TestArchiveDryRunDoesNotModify(t *testing.T) {
	vaultFile, pass := buildArchiveVault(t)
	_, err := ArchiveEntries(vaultFile, pass, []string{"OLD_KEY"}, "", true)
	if err != nil {
		t.Fatalf("ArchiveEntries dry-run: %v", err)
	}
	v, _ := LoadOrCreate(vaultFile, pass)
	if _, ok := v.Entries["OLD_KEY"]; !ok {
		t.Fatal("dry-run should not remove key from vault")
	}
	p := archivePath(vaultFile)
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("dry-run should not create archive file")
	}
}

func TestArchiveFilePermissions(t *testing.T) {
	vaultFile, pass := buildArchiveVault(t)
	_, err := ArchiveEntries(vaultFile, pass, []string{"OLD_KEY"}, "", false)
	if err != nil {
		t.Fatalf("ArchiveEntries: %v", err)
	}
	info, err := os.Stat(archivePath(vaultFile))
	if err != nil {
		t.Fatalf("stat archive: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestArchiveMissingKeyReturnsError(t *testing.T) {
	vaultFile, pass := buildArchiveVault(t)
	_, err := ArchiveEntries(vaultFile, pass, []string{"NONEXISTENT"}, "", false)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
