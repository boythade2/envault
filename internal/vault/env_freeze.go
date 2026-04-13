package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FreezeRecord represents a frozen snapshot of key values at a point in time.
type FreezeRecord struct {
	FrozenAt time.Time         `json:"frozen_at"`
	Keys     map[string]string `json:"keys"`
}

// FreezePath returns the path to the freeze file for a given vault file.
func freezePath(vaultFile string) string {
	dir := filepath.Dir(vaultFile)
	base := filepath.Base(vaultFile)
	return filepath.Join(dir, "."+base+".freeze.json")
}

// FreezeEntries records the current values of the specified keys (or all keys
// if keys is empty) into a freeze file alongside the vault.
func FreezeEntries(vaultFile string, keys []string) (*FreezeRecord, error) {
	v, err := LoadOrCreate(vaultFile)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	record := &FreezeRecord{
		FrozenAt: time.Now().UTC(),
		Keys:     make(map[string]string),
	}

	if len(keys) == 0 {
		for _, e := range v.Entries {
			record.Keys[e.Key] = e.Value
		}
	} else {
		for _, k := range keys {
			found := false
			for _, e := range v.Entries {
				if e.Key == k {
					record.Keys[k] = e.Value
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("key not found: %s", k)
			}
		}
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal freeze: %w", err)
	}

	path := freezePath(vaultFile)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, fmt.Errorf("write freeze file: %w", err)
	}

	return record, nil
}

// LoadFreeze reads the freeze file for a given vault.
func LoadFreeze(vaultFile string) (*FreezeRecord, error) {
	path := freezePath(vaultFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read freeze file: %w", err)
	}

	var record FreezeRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("parse freeze file: %w", err)
	}
	return &record, nil
}

// ThawEntry removes the freeze record for a vault.
func ThawEntry(vaultFile string) error {
	path := freezePath(vaultFile)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove freeze file: %w", err)
	}
	return nil
}
