package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const backupDirName = ".envault_backups"

func backupDir(vaultPath string) string {
	return filepath.Join(filepath.Dir(vaultPath), backupDirName)
}

// BackupMeta holds metadata about a single backup file.
type BackupMeta struct {
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	VaultPath string    `json:"vault_path"`
	Filename  string    `json:"filename"`
}

// CreateBackup writes a timestamped copy of the vault file to the backup directory.
func CreateBackup(vaultPath, label string) (BackupMeta, error) {
	data, err := os.ReadFile(vaultPath)
	if err != nil {
		return BackupMeta{}, fmt.Errorf("read vault: %w", err)
	}

	dir := backupDir(vaultPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return BackupMeta{}, fmt.Errorf("create backup dir: %w", err)
	}

	now := time.Now().UTC()
	safeLabel := strings.ReplaceAll(label, " ", "_")
	filename := fmt.Sprintf("%s_%s.json", now.Format("20060102T150405"), safeLabel)
	destPath := filepath.Join(dir, filename)

	if err := os.WriteFile(destPath, data, 0600); err != nil {
		return BackupMeta{}, fmt.Errorf("write backup: %w", err)
	}

	meta := BackupMeta{
		Label:     label,
		CreatedAt: now,
		VaultPath: vaultPath,
		Filename:  filename,
	}
	metaPath := destPath + ".meta"
	metaData, _ := json.MarshalIndent(meta, "", "  ")
	_ = os.WriteFile(metaPath, metaData, 0600)

	return meta, nil
}

// ListBackups returns all backups for a given vault path, sorted newest first.
func ListBackups(vaultPath string) ([]BackupMeta, error) {
	dir := backupDir(vaultPath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read backup dir: %w", err)
	}

	var results []BackupMeta
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".meta") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var m BackupMeta
		if err := json.Unmarshal(raw, &m); err == nil {
			results = append(results, m)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})
	return results, nil
}

// RestoreBackup copies a backup file back over the vault path.
func RestoreBackup(vaultPath, filename string) error {
	src := filepath.Join(backupDir(vaultPath), filename)
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("backup not found: %w", err)
	}
	if err := os.WriteFile(vaultPath, data, 0600); err != nil {
		return fmt.Errorf("restore vault: %w", err)
	}
	return nil
}
