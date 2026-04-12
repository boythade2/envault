package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildWatchVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	v := &Vault{}
	v.Entries = map[string]Entry{
		"KEY": {Value: "value", CreatedAt: testTime(), UpdatedAt: testTime()},
	}
	if err := v.Save(vaultPath, "pass"); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return vaultPath
}

func TestChecksumFileProducesHex(t *testing.T) {
	vaultPath := buildWatchVault(t)
	sum, err := ChecksumFile(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sum) != 64 {
		t.Errorf("expected 64-char hex, got %d chars", len(sum))
	}
}

func TestChecksumFileMissingFile(t *testing.T) {
	_, err := ChecksumFile("/nonexistent/path.vault")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveAndLoadWatchState(t *testing.T) {
	vaultPath := buildWatchVault(t)
	if err := SaveWatchState(vaultPath); err != nil {
		t.Fatalf("SaveWatchState: %v", err)
	}
	state, err := LoadWatchState(vaultPath)
	if err != nil {
		t.Fatalf("LoadWatchState: %v", err)
	}
	if state == nil {
		t.Fatal("expected non-nil state")
	}
	if state.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if state.RecordedAt.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadWatchStateNoFile(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "missing.vault")
	state, err := LoadWatchState(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != nil {
		t.Error("expected nil state when no watch file exists")
	}
}

func TestHasChangedNoState(t *testing.T) {
	vaultPath := buildWatchVault(t)
	changed, err := HasChanged(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true when no state saved")
	}
}

func TestHasChangedAfterSave(t *testing.T) {
	vaultPath := buildWatchVault(t)
	if err := SaveWatchState(vaultPath); err != nil {
		t.Fatalf("SaveWatchState: %v", err)
	}
	changed, err := HasChanged(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected changed=false immediately after saving state")
	}
	// Modify the vault file
	f, _ := os.OpenFile(vaultPath, os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString("\n")
	f.Close()
	changed, err = HasChanged(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error after modification: %v", err)
	}
	if !changed {
		t.Error("expected changed=true after modifying vault file")
	}
}
