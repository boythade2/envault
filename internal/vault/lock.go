package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockState represents a vault lock record.
type LockState struct {
	Locked    bool      `json:"locked"`
	LockedAt  time.Time `json:"locked_at,omitempty"`
	LockedBy  string    `json:"locked_by,omitempty"`
}

func lockPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	return filepath.Join(dir, "."+base+".lock")
}

// LoadLockState reads the lock file for the given vault. If no lock file
// exists, it returns an unlocked state without error.
func LoadLockState(vaultPath string) (LockState, error) {
	path := lockPath(vaultPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return LockState{Locked: false}, nil
	}
	if err != nil {
		return LockState{}, fmt.Errorf("read lock file: %w", err)
	}
	var state LockState
	if err := json.Unmarshal(data, &state); err != nil {
		return LockState{}, fmt.Errorf("parse lock file: %w", err)
	}
	return state, nil
}

// LockVault writes a lock file for the given vault path.
func LockVault(vaultPath, lockedBy string) error {
	existing, err := LoadLockState(vaultPath)
	if err != nil {
		return err
	}
	if existing.Locked {
		return fmt.Errorf("vault is already locked by %q", existing.LockedBy)
	}
	state := LockState{
		Locked:   true,
		LockedAt: time.Now().UTC(),
		LockedBy: lockedBy,
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lock state: %w", err)
	}
	path := lockPath(vaultPath)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write lock file: %w", err)
	}
	return nil
}

// UnlockVault removes the lock file for the given vault path.
func UnlockVault(vaultPath string) error {
	state, err := LoadLockState(vaultPath)
	if err != nil {
		return err
	}
	if !state.Locked {
		return fmt.Errorf("vault is not locked")
	}
	path := lockPath(vaultPath)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove lock file: %w", err)
	}
	return nil
}

// AssertUnlocked returns an error if the vault at vaultPath is locked.
func AssertUnlocked(vaultPath string) error {
	state, err := LoadLockState(vaultPath)
	if err != nil {
		return err
	}
	if state.Locked {
		return fmt.Errorf("vault is locked by %q since %s", state.LockedBy, state.LockedAt.Format(time.RFC3339))
	}
	return nil
}
