package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildLockVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadLockStateNoFile(t *testing.T) {
	vaultPath := buildLockVault(t)
	state, err := LoadLockState(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Locked {
		t.Error("expected unlocked state when no lock file exists")
	}
}

func TestLockAndUnlock(t *testing.T) {
	vaultPath := buildLockVault(t)

	if err := LockVault(vaultPath, "alice"); err != nil {
		t.Fatalf("LockVault: %v", err)
	}

	state, err := LoadLockState(vaultPath)
	if err != nil {
		t.Fatalf("LoadLockState: %v", err)
	}
	if !state.Locked {
		t.Error("expected vault to be locked")
	}
	if state.LockedBy != "alice" {
		t.Errorf("expected LockedBy=alice, got %q", state.LockedBy)
	}
	if state.LockedAt.IsZero() {
		t.Error("expected LockedAt to be set")
	}

	if err := UnlockVault(vaultPath); err != nil {
		t.Fatalf("UnlockVault: %v", err)
	}

	state, err = LoadLockState(vaultPath)
	if err != nil {
		t.Fatalf("LoadLockState after unlock: %v", err)
	}
	if state.Locked {
		t.Error("expected vault to be unlocked after UnlockVault")
	}
}

func TestLockAlreadyLocked(t *testing.T) {
	vaultPath := buildLockVault(t)
	if err := LockVault(vaultPath, "alice"); err != nil {
		t.Fatalf("first lock: %v", err)
	}
	if err := LockVault(vaultPath, "bob"); err == nil {
		t.Error("expected error when locking already-locked vault")
	}
}

func TestUnlockNotLocked(t *testing.T) {
	vaultPath := buildLockVault(t)
	if err := UnlockVault(vaultPath); err == nil {
		t.Error("expected error when unlocking a vault that is not locked")
	}
}

func TestAssertUnlocked(t *testing.T) {
	vaultPath := buildLockVault(t)

	if err := AssertUnlocked(vaultPath); err != nil {
		t.Errorf("expected no error for unlocked vault, got: %v", err)
	}

	if err := LockVault(vaultPath, "ci-bot"); err != nil {
		t.Fatalf("LockVault: %v", err)
	}
	if err := AssertUnlocked(vaultPath); err == nil {
		t.Error("expected error from AssertUnlocked on locked vault")
	}
}

func TestLockFilePermissions(t *testing.T) {
	vaultPath := buildLockVault(t)
	if err := LockVault(vaultPath, "tester"); err != nil {
		t.Fatalf("LockVault: %v", err)
	}
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	lp := filepath.Join(dir, "."+base+".lock")
	info, err := os.Stat(lp)
	if err != nil {
		t.Fatalf("stat lock file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected lock file permissions 0600, got %04o", perm)
	}
}
