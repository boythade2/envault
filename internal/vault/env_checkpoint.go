package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Checkpoint struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Entries   map[string]string `json:"entries"`
}

type CheckpointStore struct {
	Checkpoints []Checkpoint `json:"checkpoints"`
}

func checkpointPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	return filepath.Join(dir, ".envault_checkpoints.json")
}

func LoadCheckpoints(vaultPath string) (*CheckpointStore, error) {
	p := checkpointPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &CheckpointStore{}, nil
	}
	if err != nil {
		return nil, err
	}
	var store CheckpointStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

func saveCheckpoints(vaultPath string, store *CheckpointStore) error {
	p := checkpointPath(vaultPath)
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func SaveCheckpoint(vaultPath, label string, v *Vault) error {
	store, err := LoadCheckpoints(vaultPath)
	if err != nil {
		return err
	}
	for _, cp := range store.Checkpoints {
		if cp.Label == label {
			return fmt.Errorf("checkpoint %q already exists", label)
		}
	}
	entries := make(map[string]string, len(v.Entries))
	for k, e := range v.Entries {
		entries[k] = e.Value
	}
	store.Checkpoints = append(store.Checkpoints, Checkpoint{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	})
	return saveCheckpoints(vaultPath, store)
}

func ListCheckpoints(vaultPath string) ([]Checkpoint, error) {
	store, err := LoadCheckpoints(vaultPath)
	if err != nil {
		return nil, err
	}
	sort.Slice(store.Checkpoints, func(i, j int) bool {
		return store.Checkpoints[i].CreatedAt.After(store.Checkpoints[j].CreatedAt)
	})
	return store.Checkpoints, nil
}

func RestoreCheckpoint(vaultPath, label string, v *Vault) error {
	store, err := LoadCheckpoints(vaultPath)
	if err != nil {
		return err
	}
	for _, cp := range store.Checkpoints {
		if cp.Label == label {
			for k, val := range cp.Entries {
				e := v.Entries[k]
				e.Value = val
				e.UpdatedAt = time.Now().UTC()
				v.Entries[k] = e
			}
			return nil
		}
	}
	return fmt.Errorf("checkpoint %q not found", label)
}
