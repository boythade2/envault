package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ArchivedEntry struct {
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	ArchivedAt time.Time `json:"archived_at"`
	Reason     string    `json:"reason,omitempty"`
}

type ArchiveStore struct {
	Entries []ArchivedEntry `json:"entries"`
}

func archivePath(vaultFile string) string {
	dir := filepath.Dir(vaultFile)
	base := filepath.Base(vaultFile)
	return filepath.Join(dir, "."+base+".archive.json")
}

func LoadArchive(vaultFile string) (*ArchiveStore, error) {
	p := archivePath(vaultFile)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &ArchiveStore{}, nil
	}
	if err != nil {
		return nil, err
	}
	var store ArchiveStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

func saveArchive(vaultFile string, store *ArchiveStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(archivePath(vaultFile), data, 0600)
}

func ArchiveEntries(vaultFile, passphrase string, keys []string, reason string, dryRun bool) ([]string, error) {
	v, err := LoadOrCreate(vaultFile, passphrase)
	if err != nil {
		return nil, err
	}
	store, err := LoadArchive(vaultFile)
	if err != nil {
		return nil, err
	}
	var archived []string
	for _, key := range keys {
		entry, ok := v.Entries[key]
		if !ok {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		store.Entries = append(store.Entries, ArchivedEntry{
			Key:        key,
			Value:      entry.Value,
			ArchivedAt: time.Now().UTC(),
			Reason:     reason,
		})
		if !dryRun {
			delete(v.Entries, key)
		}
		archived = append(archived, key)
	}
	if !dryRun {
		if err := v.Save(vaultFile, passphrase); err != nil {
			return nil, err
		}
		if err := saveArchive(vaultFile, store); err != nil {
			return nil, err
		}
	}
	return archived, nil
}

func FormatArchiveResults(archived []string, dryRun bool) string {
	if len(archived) == 0 {
		return "no keys archived\n"
	}
	prefix := "archived"
	if dryRun {
		prefix = "[dry-run] would archive"
	}
	out := ""
	for _, k := range archived {
		out += fmt.Sprintf("%s: %s\n", prefix, k)
	}
	return out
}
