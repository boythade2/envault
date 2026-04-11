package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time copy of a vault's entries.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Entries   map[string]Entry  `json:"entries"`
}

// snapshotDir returns the directory used to store snapshots for a vault file.
func snapshotDir(vaultPath string) string {
	base := filepath.Dir(vaultPath)
	name := filepath.Base(vaultPath)
	return filepath.Join(base, ".envault_snapshots", name)
}

// SaveSnapshot writes a labelled snapshot of the given vault to disk.
func SaveSnapshot(v *Vault, vaultPath, label string) error {
	dir := snapshotDir(vaultPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create snapshot dir: %w", err)
	}

	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   v.Entries,
	}

	fileName := fmt.Sprintf("%d_%s.json", snap.CreatedAt.Unix(), sanitizeLabel(label))
	filePath := filepath.Join(dir, fileName)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	return os.WriteFile(filePath, data, 0600)
}

// ListSnapshots returns all snapshots stored for a vault file, oldest first.
func ListSnapshots(vaultPath string) ([]Snapshot, error) {
	dir := snapshotDir(vaultPath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return fmt.Errorf("read snapshot dir: %w", err)
	}

	var snapshots []Snapshot
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var s Snapshot
		if err := json.Unmarshal(data, &s); err != nil {
			continue
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

// sanitizeLabel replaces characters that are unsafe in file names.
func sanitizeLabel(label string) string {
	out := make([]byte, len(label))
	for i := 0; i < len(label); i++ {
		c := label[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			out[i] = c
		} else {
			out[i] = '_'
		}
	}
	return string(out)
}
