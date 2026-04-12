package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Permission struct {
	Key      string `json:"key"`
	ReadOnly bool   `json:"read_only"`
}

type PermissionMap map[string]Permission

func permPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	return filepath.Join(dir, name+".perms.json")
}

func LoadPermissions(vaultPath string) (PermissionMap, error) {
	p := permPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return make(PermissionMap), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read permissions: %w", err)
	}
	var pm PermissionMap
	if err := json.Unmarshal(data, &pm); err != nil {
		return nil, fmt.Errorf("parse permissions: %w", err)
	}
	return pm, nil
}

func SavePermissions(vaultPath string, pm PermissionMap) error {
	data, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal permissions: %w", err)
	}
	return os.WriteFile(permPath(vaultPath), data, 0600)
}

func SetPermission(vaultPath, key string, readOnly bool) error {
	pm, err := LoadPermissions(vaultPath)
	if err != nil {
		return err
	}
	pm[key] = Permission{Key: key, ReadOnly: readOnly}
	return SavePermissions(vaultPath, pm)
}

func RemovePermission(vaultPath, key string) error {
	pm, err := LoadPermissions(vaultPath)
	if err != nil {
		return err
	}
	delete(pm, key)
	return SavePermissions(vaultPath, pm)
}

func IsReadOnly(vaultPath, key string) (bool, error) {
	pm, err := LoadPermissions(vaultPath)
	if err != nil {
		return false, err
	}
	p, ok := pm[key]
	if !ok {
		return false, nil
	}
	return p.ReadOnly, nil
}

func AssertWritable(vaultPath, key string) error {
	ro, err := IsReadOnly(vaultPath, key)
	if err != nil {
		return err
	}
	if ro {
		return fmt.Errorf("key %q is read-only", key)
	}
	return nil
}
