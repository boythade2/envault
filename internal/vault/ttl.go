package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TTLEntry represents a TTL record for a vault key.
type TTLEntry struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TTLStore holds TTL records keyed by variable name.
type TTLStore struct {
	Entries map[string]TTLEntry `json:"entries"`
}

func ttlPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".ttl.json")
}

// LoadTTLStore loads the TTL store for the given vault file.
func LoadTTLStore(vaultPath string) (*TTLStore, error) {
	p := ttlPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &TTLStore{Entries: make(map[string]TTLEntry)}, nil
	}
	if err != nil {
		return nil, err
	}
	var store TTLStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	if store.Entries == nil {
		store.Entries = make(map[string]TTLEntry)
	}
	return &store, nil
}

// Save persists the TTL store to disk.
func (s *TTLStore) Save(vaultPath string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ttlPath(vaultPath), data, 0600)
}

// SetTTL sets a TTL duration for the given key.
func (s *TTLStore) SetTTL(key string, d time.Duration) {
	s.Entries[key] = TTLEntry{
		Key:       key,
		ExpiresAt: time.Now().Add(d),
	}
}

// RemoveTTL removes the TTL record for the given key.
func (s *TTLStore) RemoveTTL(key string) {
	delete(s.Entries, key)
}

// ExpiredKeys returns keys whose TTL has elapsed.
func (s *TTLStore) ExpiredKeys() []string {
	now := time.Now()
	var expired []string
	for k, e := range s.Entries {
		if now.After(e.ExpiresAt) {
			expired = append(expired, k)
		}
	}
	return expired
}

// TTLStatus returns a human-readable status for the given key.
func (s *TTLStore) TTLStatus(key string) string {
	e, ok := s.Entries[key]
	if !ok {
		return "no TTL set"
	}
	remaining := time.Until(e.ExpiresAt)
	if remaining <= 0 {
		return "expired"
	}
	return fmt.Sprintf("expires in %s", remaining.Round(time.Second))
}
