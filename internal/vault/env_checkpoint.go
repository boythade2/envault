package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Checkpoint represents a named point-in-time capture of vault entries.
type Checkpoint struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Entries   map[string]string `json:"entries"`
}

// CheckpointStore holds all checkpoints for a vault.
type CheckpointStore struct {
	Checkpoints []Checkpoint `json:"checkpoints"`
}

func checkpointPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".checkpoints.json")
}

// LoadCheckpoints reads the checkpoint store from disk.
// Returns an empty store if the file does not exist.
func LoadCheckpoints(vaultPath string) (*CheckpointStore, error) {
	p := checkpointPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &CheckpointStore{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read checkpoint file: %w", err)
	}
	var store CheckpointStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse checkpoint file: %w", err)
	}
	return &store, nil
}

func saveCheckpoints(vaultPath string, store *CheckpointStore) error {
	p := checkpointPath(vaultPath)
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal checkpoints: %w", err)
	}
	return os.WriteFile(p, data, 0600)
}

// SaveCheckpoint captures the current vault entries under the given label.
// If a checkpoint with the same label already exists it is overwritten.
func SaveCheckpoint(vaultPath string, v *Vault, label string) error {
	if label == "" {
		return fmt.Errorf("checkpoint label must not be empty")
	}
	store, err := LoadCheckpoints(vaultPath)
	if err != nil {
		return err
	}

	// Snapshot current values.
	snap := make(map[string]string, len(v.Entries))
	for k, e := range v.Entries {
		snap[k] = e.Value
	}

	// Replace existing checkpoint with same label, if any.
	for i, cp := range store.Checkpoints {
		if cp.Label == label {
			store.Checkpoints[i] = Checkpoint{Label: label, CreatedAt: time.Now().UTC(), Entries: snap}
			return saveCheckpoints(vaultPath, store)
		}
	}

	store.Checkpoints = append(store.Checkpoints, Checkpoint{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Entries:   snap,
	})
	return saveCheckpoints(vaultPath, store)
}

// ListCheckpoints returns all checkpoints sorted by creation time (newest first).
func ListCheckpoints(vaultPath string) ([]Checkpoint, error) {
	store, err := LoadCheckpoints(vaultPath)
	if err != nil {
		return nil, err
	}
	sorted := make([]Checkpoint, len(store.Checkpoints))
	copy(sorted, store.Checkpoints)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})
	return sorted, nil
}

// DeleteCheckpoint removes a checkpoint by label.
// Returns an error if the label is not found.
func DeleteCheckpoint(vaultPath string, label string) error {
	store, err := LoadCheckpoints(vaultPath)
	if err != nil {
		return err
	}
	idx := -1
	for i, cp := range store.Checkpoints {
		if cp.Label == label {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("checkpoint %q not found", label)
	}
	store.Checkpoints = append(store.Checkpoints[:idx], store.Checkpoints[idx+1:]...)
	return saveCheckpoints(vaultPath, store)
}
