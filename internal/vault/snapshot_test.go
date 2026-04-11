package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func buildSnapshotVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	v := &Vault{
		Entries: map[string]Entry{
			"DB_HOST": {Value: "localhost", UpdatedAt: time.Now()},
			"DB_PORT": {Value: "5432", UpdatedAt: time.Now()},
		},
	}
	return v, vaultPath
}

func TestSaveSnapshotCreatesFile(t *testing.T) {
	v, vaultPath := buildSnapshotVault(t)

	if err := SaveSnapshot(v, vaultPath, "before-migration"); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	snaps, err := ListSnapshots(vaultPath)
	if err != nil {
		t.Fatalf("ListSnapshots: %v", err)
	}
	if len(snaps) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(snaps))
	}
}

func TestSnapshotPreservesEntries(t *testing.T) {
	v, vaultPath := buildSnapshotVault(t)

	if err := SaveSnapshot(v, vaultPath, "checkpoint"); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	snaps, _ := ListSnapshots(vaultPath)
	snap := snaps[0]

	if snap.Label != "checkpoint" {
		t.Errorf("expected label 'checkpoint', got %q", snap.Label)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Entries["DB_HOST"].Value != "localhost" {
		t.Errorf("unexpected DB_HOST value: %s", snap.Entries["DB_HOST"].Value)
	}
}

func TestListSnapshotsNoDirectory(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "missing.vault")

	snaps, err := ListSnapshots(vaultPath)
	if err != nil {
		t.Fatalf("expected no error for missing dir, got %v", err)
	}
	if snaps != nil {
		t.Errorf("expected nil snapshots, got %v", snaps)
	}
}

func TestSnapshotFilePermissions(t *testing.T) {
	v, vaultPath := buildSnapshotVault(t)

	if err := SaveSnapshot(v, vaultPath, "perm-test"); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	dir := snapshotDir(vaultPath)
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		info, err := os.Stat(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0600 {
			t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
		}
	}
}

func TestSanitizeLabel(t *testing.T) {
	cases := []struct{ input, want string }{
		{"before-migration", "before-migration"},
		{"hello world", "hello_world"},
		{"v1.0.0", "v1_0_0"},
	}
	for _, c := range cases {
		got := sanitizeLabel(c.input)
		if got != c.want {
			t.Errorf("sanitizeLabel(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
