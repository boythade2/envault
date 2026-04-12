package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// KeyPermission defines who can read or write a specific key.
type KeyPermission struct {
	Key       string    `json:"key"`
	ReadOnly  bool      `json:"read_only"`
	Owner     string    `json:"owner,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PermissionMap holds permissions indexed by key name.
type PermissionMap struct {
	Permissions map[string]KeyPermission `json:"permissions"`
}

func permPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	return filepath.Join(dir, ".envault_permissions.json")
}

// LoadPermissions reads the permission map for a vault, returning an empty map if none exists.
func LoadPermissions(vaultPath string) (PermissionMap, error) {
	p := permPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return PermissionMap{Permissions: make(map[string]KeyPermission)}, nil
	}
	if err != nil {
		return PermissionMap{}, fmt.Errorf("read permissions: %w", err)
	}
	var pm PermissionMap
	if err := json.Unmarshal(data, &pm); err != nil {
		return PermissionMap{}, fmt.Errorf("parse permissions: %w", err)
	}
	if pm.Permissions == nil {
		pm.Permissions = make(map[string]KeyPermission)
	}
	return pm, nil
}

// SavePermissions writes the permission map to disk.
func SavePermissions(vaultPath string, pm PermissionMap) error {
	data, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal permissions: %w", err)
	}
	return os.WriteFile(permPath(vaultPath), data, 0600)
}

// SetPermission sets or updates the permission for a single key.
func SetPermission(vaultPath, key, owner string, readOnly bool) error {
	pm, err := LoadPermissions(vaultPath)
	if err != nil {
		return err
	}
	pm.Permissions[key] = KeyPermission{
		Key:       key,
		ReadOnly:  readOnly,
		Owner:     owner,
		UpdatedAt: time.Now().UTC(),
	}
	return SavePermissions(vaultPath, pm)
}

// RemovePermission deletes the permission entry for a key.
func RemovePermission(vaultPath, key string) error {
	pm, err := LoadPermissions(vaultPath)
	if err != nil {
		return err
	}
	delete(pm.Permissions, key)
	return SavePermissions(vaultPath, pm)
}

// IsReadOnly returns true if the key is marked read-only.
func IsReadOnly(vaultPath, key string) (bool, error) {
	pm, err := LoadPermissions(vaultPath)
	if err != nil {
		return false, err
	}
	p, ok := pm.Permissions[key]
	if !ok {
		return false, nil
	}
	return p.ReadOnly, nil
}
