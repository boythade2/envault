package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func buildCheckpointVault(t *testing.T) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.vault")
	v := &Vault{
		Entries: map[string]Entry{
			"KEY_A": {Value: "alpha", UpdatedAt: time.Now().UTC()},
			"KEY_B": {Value: "beta", UpdatedAt: time.Now().UTC()},
		},
	}
	_ = os.WriteFile(p, []byte(`{}`), 0600)
	return p, v
}

func TestLoadCheckpointsNoFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "missing.vault")
	store, err := LoadCheckpoints(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Checkpoints) != 0 {
		t.Errorf("expected empty store, got %d entries", len(store.Checkpoints))
	}
}

func TestSaveAndListCheckpoints(t *testing.T) {
	p, v := buildCheckpointVault(t)
	if err := SaveCheckpoint(p, "v1", v); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}
	list, err := ListCheckpoints(p)
	if err != nil {
		t.Fatalf("ListCheckpoints: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 checkpoint, got %d", len(list))
	}
	if list[0].Label != "v1" {
		t.Errorf("expected label v1, got %s", list[0].Label)
	}
}

func TestSaveCheckpointDuplicateLabelReturnsError(t *testing.T) {
	p, v := buildCheckpointVault(t)
	if err := SaveCheckpoint(p, "dup", v); err != nil {
		t.Fatalf("first save: %v", err)
	}
	if err := SaveCheckpoint(p, "dup", v); err == nil {
		t.Error("expected error for duplicate label, got nil")
	}
}

func TestRestoreCheckpointPreservesValues(t *testing.T) {
	p, v := buildCheckpointVault(t)
	if err := SaveCheckpoint(p, "snap1", v); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}
	v.Entries["KEY_A"] = Entry{Value: "changed"}
	if err := RestoreCheckpoint(p, "snap1", v); err != nil {
		t.Fatalf("RestoreCheckpoint: %v", err)
	}
	if v.Entries["KEY_A"].Value != "alpha" {
		t.Errorf("expected alpha after restore, got %s", v.Entries["KEY_A"].Value)
	}
}

func TestRestoreCheckpointNotFound(t *testing.T) {
	p, v := buildCheckpointVault(t)
	if err := RestoreCheckpoint(p, "ghost", v); err == nil {
		t.Error("expected error for missing checkpoint, got nil")
	}
}

func TestCheckpointFilePermissions(t *testing.T) {
	p, v := buildCheckpointVault(t)
	if err := SaveCheckpoint(p, "perm-test", v); err != nil {
		t.Fatalf("SaveCheckpoint: %v", err)
	}
	cp := checkpointPath(p)
	info, err := os.Stat(cp)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
