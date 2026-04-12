package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func buildBackupVault(t *testing.T) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	v := &Vault{
		Entries: map[string]Entry{
			"API_KEY": {Value: "secret", UpdatedAt: time.Now()},
		},
	}
	if err := v.Save(path); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path, v
}

func TestCreateBackupCreatesFile(t *testing.T) {
	path, _ := buildBackupVault(t)
	meta, err := CreateBackup(path, "before-deploy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	backupFile := filepath.Join(backupDir(path), meta.Filename)
	if _, err := os.Stat(backupFile); err != nil {
		t.Errorf("backup file not found: %v", err)
	}
	if meta.Label != "before-deploy" {
		t.Errorf("expected label 'before-deploy', got %q", meta.Label)
	}
}

func TestCreateBackupMetaFile(t *testing.T) {
	path, _ := buildBackupVault(t)
	meta, err := CreateBackup(path, "test-meta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	metaFile := filepath.Join(backupDir(path), meta.Filename+".meta")
	if _, err := os.Stat(metaFile); err != nil {
		t.Errorf("meta file not found: %v", err)
	}
}

func TestListBackupsNoDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	backups, err := ListBackups(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("expected empty list, got %d", len(backups))
	}
}

func TestListBackupsReturnsSortedNewestFirst(t *testing.T) {
	path, _ := buildBackupVault(t)
	CreateBackup(path, "first")
	time.Sleep(10 * time.Millisecond)
	CreateBackup(path, "second")

	backups, err := ListBackups(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 2 {
		t.Fatalf("expected 2 backups, got %d", len(backups))
	}
	if backups[0].Label != "second" {
		t.Errorf("expected newest first, got label %q", backups[0].Label)
	}
}

func TestRestoreBackupRestoresContent(t *testing.T) {
	path, _ := buildBackupVault(t)
	meta, err := CreateBackup(path, "restore-test")
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}

	// Overwrite vault with different data
	v2 := &Vault{Entries: map[string]Entry{}}
	v2.Save(path)

	if err := RestoreBackup(path, meta.Filename); err != nil {
		t.Fatalf("restore: %v", err)
	}

	v3, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if _, ok := v3.Entries["API_KEY"]; !ok {
		t.Error("expected API_KEY to be restored")
	}
}

func TestRestoreBackupNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.vault")
	err := RestoreBackup(path, "nonexistent.json")
	if err == nil {
		t.Error("expected error for missing backup")
	}
}
