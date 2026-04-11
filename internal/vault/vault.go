package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// VaultEntry represents a single encrypted environment file entry.
type VaultEntry struct {
	Name      string    `json:"name"`
	FilePath  string    `json:"file_path"`
	Encrypted bool      `json:"encrypted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Vault holds the metadata for all managed environment files.
type Vault struct {
	Version int                    `json:"version"`
	Entries map[string]*VaultEntry `json:"entries"`
}

const vaultFileName = ".envault.json"

// LoadOrCreate loads an existing vault index from disk or creates a new one.
func LoadOrCreate(dir string) (*Vault, error) {
	path := filepath.Join(dir, vaultFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Vault{Version: 1, Entries: make(map[string]*VaultEntry)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading vault index: %w", err)
	}
	var v Vault
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parsing vault index: %w", err)
	}
	if v.Entries == nil {
		v.Entries = make(map[string]*VaultEntry)
	}
	return &v, nil
}

// Save persists the vault index to disk.
func (v *Vault) Save(dir string) error {
	path := filepath.Join(dir, vaultFileName)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("serialising vault index: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// AddEntry registers or updates an entry in the vault index.
func (v *Vault) AddEntry(name, filePath string, encrypted bool) {
	now := time.Now().UTC()
	if existing, ok := v.Entries[name]; ok {
		existing.FilePath = filePath
		existing.Encrypted = encrypted
		existing.UpdatedAt = now
		return
	}
	v.Entries[name] = &VaultEntry{
		Name:      name,
		FilePath:  filePath,
		Encrypted: encrypted,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// RemoveEntry deletes an entry from the vault index.
func (v *Vault) RemoveEntry(name string) bool {
	if _, ok := v.Entries[name]; !ok {
		return false
	}
	delete(v.Entries, name)
	return true
}

// GetEntry returns the VaultEntry for the given name, or an error if it does
// not exist.
func (v *Vault) GetEntry(name string) (*VaultEntry, error) {
	entry, ok := v.Entries[name]
	if !ok {
		return nil, fmt.Errorf("entry %q not found in vault", name)
	}
	return entry, nil
}
