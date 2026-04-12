package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AliasStore struct {
	Aliases map[string]string `json:"aliases"` // alias -> canonical key
}

func aliasPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	return filepath.Join(dir, name+".aliases.json")
}

func LoadAliases(vaultPath string) (*AliasStore, error) {
	p := aliasPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &AliasStore{Aliases: make(map[string]string)}, nil
	}
	if err != nil {
		return nil, err
	}
	var store AliasStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	if store.Aliases == nil {
		store.Aliases = make(map[string]string)
	}
	return &store, nil
}

func saveAliases(vaultPath string, store *AliasStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(aliasPath(vaultPath), data, 0600)
}

func AddAlias(vaultPath, alias, key string) error {
	store, err := LoadAliases(vaultPath)
	if err != nil {
		return err
	}
	if _, exists := store.Aliases[alias]; exists {
		return fmt.Errorf("alias %q already exists", alias)
	}
	store.Aliases[alias] = key
	return saveAliases(vaultPath, store)
}

func RemoveAlias(vaultPath, alias string) error {
	store, err := LoadAliases(vaultPath)
	if err != nil {
		return err
	}
	if _, exists := store.Aliases[alias]; !exists {
		return fmt.Errorf("alias %q not found", alias)
	}
	delete(store.Aliases, alias)
	return saveAliases(vaultPath, store)
}

func ResolveAlias(vaultPath, nameOrAlias string) (string, error) {
	store, err := LoadAliases(vaultPath)
	if err != nil {
		return "", err
	}
	if canonical, ok := store.Aliases[nameOrAlias]; ok {
		return canonical, nil
	}
	return nameOrAlias, nil
}
