package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records a change made to a vault entry.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"` // "add", "update", "remove"
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
}

// History holds a log of changes for a vault file.
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// historyPath returns the path to the history file for a given vault file.
func historyPath(vaultPath string) string {
	ext := filepath.Ext(vaultPath)
	base := vaultPath[:len(vaultPath)-len(ext)]
	return base + ".history.json"
}

// LoadHistory loads the history file associated with the given vault path.
// If no history file exists, an empty History is returned.
func LoadHistory(vaultPath string) (*History, error) {
	path := historyPath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &History{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading history file: %w", err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("parsing history file: %w", err)
	}
	return &h, nil
}

// Record appends a new entry to the history and persists it to disk.
func (h *History) Record(vaultPath, action, key, oldValue, newValue string) error {
	h.Entries = append(h.Entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		OldValue:  oldValue,
		NewValue:  newValue,
	})
	return h.save(vaultPath)
}

// EntriesForKey returns all history entries for the given key in chronological order.
func (h *History) EntriesForKey(key string) []HistoryEntry {
	var result []HistoryEntry
	for _, e := range h.Entries {
		if e.Key == key {
			result = append(result, e)
		}
	}
	return result
}

// save writes the history to disk.
func (h *History) save(vaultPath string) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling history: %w", err)
	}
	path := historyPath(vaultPath)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing history file: %w", err)
	}
	return nil
}
