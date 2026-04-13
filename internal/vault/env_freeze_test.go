package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildFreezeVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vf := filepath.Join(dir, "test.vault")
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "API_KEY", Value: "secret"},
	}
	if err := v.Save(vf); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return vf
}

func TestFreezeAllKeys(t *testing.T) {
	vf := buildFreezeVault(t)
	rec, err := FreezeEntries(vf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Keys) != 3 {
		t.Errorf("expected 3 frozen keys, got %d", len(rec.Keys))
	}
	if rec.Keys["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", rec.Keys["DB_HOST"])
	}
}

func TestFreezeSelectedKeys(t *testing.T) {
	vf := buildFreezeVault(t)
	rec, err := FreezeEntries(vf, []string{"DB_HOST", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.Keys) != 2 {
		t.Errorf("expected 2 frozen keys, got %d", len(rec.Keys))
	}
	if _, ok := rec.Keys["DB_PORT"]; ok {
		t.Error("DB_PORT should not be in frozen keys")
	}
}

func TestFreezeNonExistentKeyReturnsError(t *testing.T) {
	vf := buildFreezeVault(t)
	_, err := FreezeEntries(vf, []string{"MISSING_KEY"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestLoadFreezeNoFile(t *testing.T) {
	dir := t.TempDir()
	vf := filepath.Join(dir, "test.vault")
	rec, err := LoadFreeze(vf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Error("expected nil record when no freeze file exists")
	}
}

func TestFreezeAndReload(t *testing.T) {
	vf := buildFreezeVault(t)
	_, err := FreezeEntries(vf, nil)
	if err != nil {
		t.Fatalf("freeze: %v", err)
	}
	rec, err := LoadFreeze(vf)
	if err != nil {
		t.Fatalf("load freeze: %v", err)
	}
	if rec == nil {
		t.Fatal("expected non-nil freeze record")
	}
	if rec.Keys["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %s", rec.Keys["DB_PORT"])
	}
	if rec.FrozenAt.IsZero() {
		t.Error("expected non-zero FrozenAt timestamp")
	}
}

func TestThawRemovesFreezeFile(t *testing.T) {
	vf := buildFreezeVault(t)
	_, err := FreezeEntries(vf, nil)
	if err != nil {
		t.Fatalf("freeze: %v", err)
	}
	if err := ThawEntry(vf); err != nil {
		t.Fatalf("thaw: %v", err)
	}
	path := freezePath(vf)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected freeze file to be removed after thaw")
	}
}

func TestFreezeFilePermissions(t *testing.T) {
	vf := buildFreezeVault(t)
	_, err := FreezeEntries(vf, nil)
	if err != nil {
		t.Fatalf("freeze: %v", err)
	}
	path := freezePath(vf)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
