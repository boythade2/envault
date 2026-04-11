package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadHistoryNoFile(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "test.vault")

	h, err := LoadHistory(vaultPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestRecordAndReload(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "test.vault")

	h, _ := LoadHistory(vaultPath)
	if err := h.Record(vaultPath, "add", "API_KEY", "", "secret123"); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	if err := h.Record(vaultPath, "update", "API_KEY", "secret123", "newsecret"); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	h2, err := LoadHistory(vaultPath)
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(h2.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h2.Entries))
	}
	if h2.Entries[0].Action != "add" || h2.Entries[0].Key != "API_KEY" {
		t.Errorf("unexpected first entry: %+v", h2.Entries[0])
	}
	if h2.Entries[1].Action != "update" || h2.Entries[1].OldValue != "secret123" {
		t.Errorf("unexpected second entry: %+v", h2.Entries[1])
	}
}

func TestHistoryFilePermissions(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "test.vault")

	h, _ := LoadHistory(vaultPath)
	h.Record(vaultPath, "remove", "DB_PASS", "oldpass", "")

	histPath := historyPath(vaultPath)
	info, err := os.Stat(histPath)
	if err != nil {
		t.Fatalf("stat history file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}

func TestHistoryTimestampSet(t *testing.T) {
	tmp := t.TempDir()
	vaultPath := filepath.Join(tmp, "test.vault")

	h, _ := LoadHistory(vaultPath)
	h.Record(vaultPath, "add", "TOKEN", "", "abc")

	if h.Entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
