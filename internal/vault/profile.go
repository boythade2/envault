package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Profile represents a named environment profile (e.g. "dev", "staging", "prod").
type Profile struct {
	Name      string    `json:"name"`
	VaultFile string    `json:"vault_file"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProfileStore holds all registered profiles for a project directory.
type ProfileStore struct {
	Profiles map[string]Profile `json:"profiles"`
}

func profilePath(dir string) string {
	return filepath.Join(dir, ".envault_profiles.json")
}

// LoadProfiles loads the profile store from dir, returning an empty store if none exists.
func LoadProfiles(dir string) (*ProfileStore, error) {
	path := profilePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &ProfileStore{Profiles: make(map[string]Profile)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read profiles: %w", err)
	}
	var store ProfileStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse profiles: %w", err)
	}
	if store.Profiles == nil {
		store.Profiles = make(map[string]Profile)
	}
	return &store, nil
}

// Save persists the profile store to dir.
func (ps *ProfileStore) Save(dir string) error {
	data, err := json.MarshalIndent(ps, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profiles: %w", err)
	}
	return os.WriteFile(profilePath(dir), data, 0600)
}

// AddProfile registers a new profile, returning an error if it already exists.
func (ps *ProfileStore) AddProfile(name, vaultFile string) error {
	if _, exists := ps.Profiles[name]; exists {
		return fmt.Errorf("profile %q already exists", name)
	}
	now := time.Now().UTC()
	ps.Profiles[name] = Profile{
		Name:      name,
		VaultFile: vaultFile,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return nil
}

// RemoveProfile deletes a profile by name.
func (ps *ProfileStore) RemoveProfile(name string) error {
	if _, exists := ps.Profiles[name]; !exists {
		return fmt.Errorf("profile %q not found", name)
	}
	delete(ps.Profiles, name)
	return nil
}

// GetProfile returns the profile with the given name.
func (ps *ProfileStore) GetProfile(name string) (Profile, error) {
	p, ok := ps.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}
