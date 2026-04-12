package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PinnedEntry represents a key that has been pinned to prevent accidental modification.
type PinnedEntry struct {
	Key       string    `json:"key"`
	PinnedAt  time.Time `json:"pinned_at"`
	Note      string    `json:"note,omitempty"`
}

// PinList holds all pinned keys for a vault.
type PinList struct {
	Pins []PinnedEntry `json:"pins"`
}

func pinPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".pins.json")
}

// LoadPins loads the pin list for the given vault file.
func LoadPins(vaultPath string) (*PinList, error) {
	p := pinPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &PinList{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read pin file: %w", err)
	}
	var pl PinList
	if err := json.Unmarshal(data, &pl); err != nil {
		return nil, fmt.Errorf("parse pin file: %w", err)
	}
	return &pl, nil
}

// PinKey adds a key to the pin list.
func PinKey(vaultPath, key, note string) error {
	pl, err := LoadPins(vaultPath)
	if err != nil {
		return err
	}
	for _, p := range pl.Pins {
		if p.Key == key {
			return fmt.Errorf("key %q is already pinned", key)
		}
	}
	pl.Pins = append(pl.Pins, PinnedEntry{Key: key, PinnedAt: time.Now().UTC(), Note: note})
	return savePins(vaultPath, pl)
}

// UnpinKey removes a key from the pin list.
func UnpinKey(vaultPath, key string) error {
	pl, err := LoadPins(vaultPath)
	if err != nil {
		return err
	}
	filtered := pl.Pins[:0]
	found := false
	for _, p := range pl.Pins {
		if p.Key == key {
			found = true
			continue
		}
		filtered = append(filtered, p)
	}
	if !found {
		return fmt.Errorf("key %q is not pinned", key)
	}
	pl.Pins = filtered
	return savePins(vaultPath, pl)
}

// IsPinned returns true if the given key is pinned.
func IsPinned(vaultPath, key string) (bool, error) {
	pl, err := LoadPins(vaultPath)
	if err != nil {
		return false, err
	}
	for _, p := range pl.Pins {
		if p.Key == key {
			return true, nil
		}
	}
	return false, nil
}

func savePins(vaultPath string, pl *PinList) error {
	data, err := json.MarshalIndent(pl, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal pins: %w", err)
	}
	return os.WriteFile(pinPath(vaultPath), data, 0600)
}
