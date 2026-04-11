package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// RestoreSnapshot loads a snapshot by label and overwrites the target vault file.
// It returns the restored Vault so callers can inspect or re-encrypt it.
func RestoreSnapshot(vaultPath, label, passphrase string) (*Vault, error) {
	dir := snapshotDir(vaultPath)
	safe := sanitizeLabel(label)
	snapshotFile := filepath.Join(dir, safe+".json")

	data, err := os.ReadFile(snapshotFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot %q not found", label)
		}
		return nil, fmt.Errorf("reading snapshot: %w", err)
	}

	var snap snapshotFile_
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("parsing snapshot: %w", err)
	}

	if len(snap.Entries) == 0 {
		return nil, fmt.Errorf("snapshot %q contains no entries", label)
	}

	// Build a fresh vault from the snapshot entries.
	v := &Vault{
		Entries:  snap.Entries,
		FilePath: vaultPath,
	}

	if err := v.Save(passphrase); err != nil {
		return nil, fmt.Errorf("saving restored vault: %w", err)
	}

	return v, nil
}

// snapshotFile_ mirrors the JSON structure written by SaveSnapshot.
type snapshotFile_ struct {
	Label     string           `json:"label"`
	CreatedAt string           `json:"created_at"`
	Entries   map[string]Entry `json:"entries"`
}
