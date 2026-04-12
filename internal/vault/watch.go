package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WatchState records the last-known checksum of a vault file.
type WatchState struct {
	VaultPath string    `json:"vault_path"`
	Checksum  string    `json:"checksum"`
	RecordedAt time.Time `json:"recorded_at"`
}

func watchStatePath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	return filepath.Join(dir, ".envault_watch")
}

// ChecksumFile returns the SHA-256 hex digest of the file at path.
func ChecksumFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:]), nil
}

// SaveWatchState persists the current checksum of vaultPath to disk.
func SaveWatchState(vaultPath string) error {
	checksum, err := ChecksumFile(vaultPath)
	if err != nil {
		return err
	}
	state := WatchState{
		VaultPath:  vaultPath,
		Checksum:   checksum,
		RecordedAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal watch state: %w", err)
	}
	return os.WriteFile(watchStatePath(vaultPath), data, 0600)
}

// LoadWatchState reads the previously saved watch state for vaultPath.
// Returns nil, nil if no state file exists.
func LoadWatchState(vaultPath string) (*WatchState, error) {
	data, err := os.ReadFile(watchStatePath(vaultPath))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read watch state: %w", err)
	}
	var state WatchState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unmarshal watch state: %w", err)
	}
	return &state, nil
}

// HasChanged reports whether vaultPath has changed since the last saved state.
// Returns true if no state has been saved yet.
func HasChanged(vaultPath string) (bool, error) {
	state, err := LoadWatchState(vaultPath)
	if err != nil {
		return false, err
	}
	if state == nil {
		return true, nil
	}
	current, err := ChecksumFile(vaultPath)
	if err != nil {
		return false, err
	}
	return current != state.Checksum, nil
}
