package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadAuditLogNoFile(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	log, err := LoadAuditLog(vaultPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Events) != 0 {
		t.Errorf("expected empty log, got %d events", len(log.Events))
	}
}

func TestRecordAndReloadAudit(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")

	if err := RecordEvent(vaultPath, "add", "API_KEY", "added via CLI"); err != nil {
		t.Fatalf("RecordEvent: %v", err)
	}
	if err := RecordEvent(vaultPath, "remove", "DB_PASS", ""); err != nil {
		t.Fatalf("RecordEvent: %v", err)
	}

	log, err := LoadAuditLog(vaultPath)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(log.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(log.Events))
	}
	if log.Events[0].Action != "add" || log.Events[0].Key != "API_KEY" {
		t.Errorf("unexpected first event: %+v", log.Events[0])
	}
	if log.Events[1].Action != "remove" || log.Events[1].Key != "DB_PASS" {
		t.Errorf("unexpected second event: %+v", log.Events[1])
	}
}

func TestAuditEventTimestampSet(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	before := time.Now().UTC().Add(-time.Second)

	if err := RecordEvent(vaultPath, "rotate", "", "passphrase rotated"); err != nil {
		t.Fatalf("RecordEvent: %v", err)
	}

	log, _ := LoadAuditLog(vaultPath)
	if log.Events[0].Timestamp.Before(before) {
		t.Errorf("timestamp not set correctly: %v", log.Events[0].Timestamp)
	}
}

func TestAuditFilePermissions(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")

	if err := RecordEvent(vaultPath, "add", "SECRET", ""); err != nil {
		t.Fatalf("RecordEvent: %v", err)
	}

	path := auditPath(vaultPath)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
