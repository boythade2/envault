package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Group struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

type GroupStore struct {
	Groups map[string]Group `json:"groups"`
}

func groupPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".groups.json")
}

func LoadGroups(vaultPath string) (*GroupStore, error) {
	p := groupPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &GroupStore{Groups: make(map[string]Group)}, nil
	}
	if err != nil {
		return nil, err
	}
	var gs GroupStore
	if err := json.Unmarshal(data, &gs); err != nil {
		return nil, err
	}
	if gs.Groups == nil {
		gs.Groups = make(map[string]Group)
	}
	return &gs, nil
}

func saveGroups(vaultPath string, gs *GroupStore) error {
	data, err := json.MarshalIndent(gs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(groupPath(vaultPath), data, 0600)
}

func AddGroup(vaultPath, name string) error {
	gs, err := LoadGroups(vaultPath)
	if err != nil {
		return err
	}
	if _, exists := gs.Groups[name]; exists {
		return fmt.Errorf("group %q already exists", name)
	}
	gs.Groups[name] = Group{Name: name, Keys: []string{}}
	return saveGroups(vaultPath, gs)
}

func RemoveGroup(vaultPath, name string) error {
	gs, err := LoadGroups(vaultPath)
	if err != nil {
		return err
	}
	if _, exists := gs.Groups[name]; !exists {
		return fmt.Errorf("group %q not found", name)
	}
	delete(gs.Groups, name)
	return saveGroups(vaultPath, gs)
}

func AssignKeyToGroup(vaultPath, name, key string) error {
	gs, err := LoadGroups(vaultPath)
	if err != nil {
		return err
	}
	g, exists := gs.Groups[name]
	if !exists {
		return fmt.Errorf("group %q not found", name)
	}
	for _, k := range g.Keys {
		if k == key {
			return fmt.Errorf("key %q already in group %q", key, name)
		}
	}
	g.Keys = append(g.Keys, key)
	sort.Strings(g.Keys)
	gs.Groups[name] = g
	return saveGroups(vaultPath, gs)
}

func UnassignKeyFromGroup(vaultPath, name, key string) error {
	gs, err := LoadGroups(vaultPath)
	if err != nil {
		return err
	}
	g, exists := gs.Groups[name]
	if !exists {
		return fmt.Errorf("group %q not found", name)
	}
	newKeys := make([]string, 0, len(g.Keys))
	for _, k := range g.Keys {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}
	if len(newKeys) == len(g.Keys) {
		return fmt.Errorf("key %q not found in group %q", key, name)
	}
	g.Keys = newKeys
	gs.Groups[name] = g
	return saveGroups(vaultPath, gs)
}
