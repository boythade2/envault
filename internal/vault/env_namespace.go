package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type NamespaceStore struct {
	Namespaces map[string][]string `json:"namespaces"` // namespace -> list of keys
}

func namespacePath(vaultFile string) string {
	return filepath.Join(filepath.Dir(vaultFile), ".envault_namespaces.json")
}

func LoadNamespaces(vaultFile string) (*NamespaceStore, error) {
	path := namespacePath(vaultFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &NamespaceStore{Namespaces: make(map[string][]string)}, nil
	}
	if err != nil {
		return nil, err
	}
	var store NamespaceStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	if store.Namespaces == nil {
		store.Namespaces = make(map[string][]string)
	}
	return &store, nil
}

func saveNamespaces(vaultFile string, store *NamespaceStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(namespacePath(vaultFile), data, 0600)
}

func AssignNamespace(vaultFile, namespace, key string) error {
	store, err := LoadNamespaces(vaultFile)
	if err != nil {
		return err
	}
	for _, k := range store.Namespaces[namespace] {
		if k == key {
			return fmt.Errorf("key %q already in namespace %q", key, namespace)
		}
	}
	store.Namespaces[namespace] = append(store.Namespaces[namespace], key)
	return saveNamespaces(vaultFile, store)
}

func UnassignNamespace(vaultFile, namespace, key string) error {
	store, err := LoadNamespaces(vaultFile)
	if err != nil {
		return err
	}
	keys := store.Namespaces[namespace]
	updated := keys[:0]
	for _, k := range keys {
		if k != key {
			updated = append(updated, k)
		}
	}
	if len(updated) == len(keys) {
		return fmt.Errorf("key %q not found in namespace %q", key, namespace)
	}
	store.Namespaces[namespace] = updated
	return saveNamespaces(vaultFile, store)
}

func GetNamespaceKeys(vaultFile, namespace string) ([]string, error) {
	store, err := LoadNamespaces(vaultFile)
	if err != nil {
		return nil, err
	}
	keys, ok := store.Namespaces[namespace]
	if !ok {
		return nil, fmt.Errorf("namespace %q not found", namespace)
	}
	return keys, nil
}

func FormatNamespaceList(store *NamespaceStore) string {
	if len(store.Namespaces) == 0 {
		return "no namespaces defined"
	}
	var sb strings.Builder
	for ns, keys := range store.Namespaces {
		sb.WriteString(fmt.Sprintf("[%s]\n", ns))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
