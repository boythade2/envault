package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func buildTTLVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadTTLStoreNoFile(t *testing.T) {
	vaultPath := buildTTLVault(t)
	store, err := LoadTTLStore(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store, got %d entries", len(store.Entries))
	}
}

func TestSetAndLoadTTL(t *testing.T) {
	vaultPath := buildTTLVault(t)
	store, _ := LoadTTLStore(vaultPath)
	store.SetTTL("API_KEY", 10*time.Minute)
	if err := store.Save(vaultPath); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	reloaded, err := LoadTTLStore(vaultPath)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	e, ok := reloaded.Entries["API_KEY"]
	if !ok {
		t.Fatal("expected API_KEY in reloaded store")
	}
	if time.Until(e.ExpiresAt) <= 0 {
		t.Error("expected future expiry")
	}
}

func TestExpiredKeys(t *testing.T) {
	vaultPath := buildTTLVault(t)
	store, _ := LoadTTLStore(vaultPath)
	store.SetTTL("FRESH", 10*time.Minute)
	store.Entries["STALE"] = TTLEntry{
		Key:       "STALE",
		ExpiresAt: time.Now().Add(-1 * time.Second),
	}

	expired := store.ExpiredKeys()
	if len(expired) != 1 || expired[0] != "STALE" {
		t.Errorf("expected [STALE], got %v", expired)
	}
}

func TestRemoveTTL(t *testing.T) {
	vaultPath := buildTTLVault(t)
	store, _ := LoadTTLStore(vaultPath)
	store.SetTTL("DB_PASS", 5*time.Minute)
	store.RemoveTTL("DB_PASS")
	if _, ok := store.Entries["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be removed")
	}
}

func TestTTLStatus(t *testing.T) {
	vaultPath := buildTTLVault(t)
	store, _ := LoadTTLStore(vaultPath)

	if s := store.TTLStatus("MISSING"); s != "no TTL set" {
		t.Errorf("unexpected status: %s", s)
	}

	store.Entries["OLD"] = TTLEntry{Key: "OLD", ExpiresAt: time.Now().Add(-time.Second)}
	if s := store.TTLStatus("OLD"); s != "expired" {
		t.Errorf("expected 'expired', got %s", s)
	}

	store.SetTTL("NEW", 2*time.Minute)
	if s := store.TTLStatus("NEW"); s == "no TTL set" || s == "expired" {
		t.Errorf("unexpected status for live key: %s", s)
	}
}

func TestTTLFilePermissions(t *testing.T) {
	vaultPath := buildTTLVault(t)
	store, _ := LoadTTLStore(vaultPath)
	store.SetTTL("X", time.Minute)
	if err := store.Save(vaultPath); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	info, err := os.Stat(ttlPath(vaultPath))
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
