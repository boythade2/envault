package vault

import (
	"fmt"
	"time"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey string
	NewKey string
}

// RenameEntry renames an existing key in the vault to a new key name.
// It returns an error if the old key does not exist or the new key already exists.
func RenameEntry(vaultPath, passphrase, oldKey, newKey string) (*RenameResult, error) {
	if oldKey == "" {
		return nil, fmt.Errorf("old key must not be empty")
	}
	if newKey == "" {
		return nil, fmt.Errorf("new key must not be empty")
	}
	if oldKey == newKey {
		return nil, fmt.Errorf("old key and new key are the same")
	}

	v, err := LoadOrCreate(vaultPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	entry, exists := v.Entries[oldKey]
	if !exists {
		return nil, fmt.Errorf("key %q not found in vault", oldKey)
	}

	if _, conflict := v.Entries[newKey]; conflict {
		return nil, fmt.Errorf("key %q already exists in vault", newKey)
	}

	// Copy entry under new key, update timestamp.
	entry.UpdatedAt = time.Now().UTC()
	v.Entries[newKey] = entry
	delete(v.Entries, oldKey)

	if err := v.Save(vaultPath, passphrase); err != nil {
		return nil, fmt.Errorf("save vault: %w", err)
	}

	return &RenameResult{OldKey: oldKey, NewKey: newKey}, nil
}
