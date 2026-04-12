package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func buildCopyVault(t *testing.T, dir, name, pass string, entries []Entry) string {
	t.Helper()
	v := &Vault{}
	for _, e := range entries {
		e.CreatedAt = time.Now().UTC()
		e.UpdatedAt = time.Now().UTC()
		v.Entries = append(v.Entries, e)
	}
	path := filepath.Join(dir, name)
	if err := v.Save(path, pass); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path
}

func TestCopyEntrySuccess(t *testing.T) {
	dir := t.TempDir()
	src := buildCopyVault(t, dir, "src.vault", "pass", []Entry{{Key: "API_KEY", Value: "secret"}})
	dst := filepath.Join(dir, "dst.vault")

	res, err := CopyEntry(src, dst, "API_KEY", "", "pass", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.DestKey != "API_KEY" {
		t.Errorf("expected dest key API_KEY, got %s", res.DestKey)
	}

	v, _ := LoadOrCreate(dst, "pass")
	if len(v.Entries) != 1 || v.Entries[0].Value != "secret" {
		t.Error("destination vault missing expected entry")
	}
}

func TestCopyEntryRename(t *testing.T) {
	dir := t.TempDir()
	src := buildCopyVault(t, dir, "src.vault", "pass", []Entry{{Key: "FOO", Value: "bar"}})
	dst := buildCopyVault(t, dir, "dst.vault", "pass", nil)

	res, err := CopyEntry(src, dst, "FOO", "BAR", "pass", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.DestKey != "BAR" {
		t.Errorf("expected dest key BAR, got %s", res.DestKey)
	}
}

func TestCopyEntryOverwriteRejected(t *testing.T) {
	dir := t.TempDir()
	src := buildCopyVault(t, dir, "src.vault", "pass", []Entry{{Key: "X", Value: "new"}})
	dst := buildCopyVault(t, dir, "dst.vault", "pass", []Entry{{Key: "X", Value: "old"}})

	_, err := CopyEntry(src, dst, "X", "", "pass", false)
	if err == nil {
		t.Fatal("expected error when overwrite is false")
	}
}

func TestCopyEntryOverwriteAllowed(t *testing.T) {
	dir := t.TempDir()
	src := buildCopyVault(t, dir, "src.vault", "pass", []Entry{{Key: "X", Value: "new"}})
	dst := buildCopyVault(t, dir, "dst.vault", "pass", []Entry{{Key: "X", Value: "old"}})

	res, err := CopyEntry(src, dst, "X", "", "pass", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Overwrote {
		t.Error("expected Overwrote to be true")
	}
	v, _ := LoadOrCreate(dst, "pass")
	if v.Entries[0].Value != "new" {
		t.Errorf("expected value 'new', got %s", v.Entries[0].Value)
	}
}

func TestCopyEntryMissingKey(t *testing.T) {
	dir := t.TempDir()
	src := buildCopyVault(t, dir, "src.vault", "pass", []Entry{{Key: "A", Value: "1"}})
	dst := filepath.Join(dir, "dst.vault")

	_, err := CopyEntry(src, dst, "MISSING", "", "pass", false)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestCopyEntryFilePermissions(t *testing.T) {
	dir := t.TempDir()
	src := buildCopyVault(t, dir, "src.vault", "pass", []Entry{{Key: "K", Value: "v"}})
	dst := filepath.Join(dir, "dst.vault")

	if _, err := CopyEntry(src, dst, "K", "", "pass", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("stat dst: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
