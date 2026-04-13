package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func buildPromoteVault(t *testing.T, dir string, name string, keys map[string]string) string {
	t.Helper()
	v := &Vault{Entries: make(map[string]Entry)}
	for k, val := range keys {
		v.Entries[k] = Entry{Value: val, UpdatedAt: time.Now()}
	}
	path := filepath.Join(dir, name)
	if err := v.Save(path, "pass"); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return path
}

func TestPromoteAllKeys(t *testing.T) {
	dir := t.TempDir()
	src := buildPromoteVault(t, dir, "src.vault", map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"})
	dst := buildPromoteVault(t, dir, "dst.vault", map[string]string{})

	results, err := PromoteEntries(src, dst, "pass", PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Skipped {
			t.Errorf("key %s should not be skipped", r.Key)
		}
	}

	v, _ := LoadOrCreate(dst, "pass")
	if v.Entries["DB_HOST"].Value != "localhost" {
		t.Errorf("expected DB_HOST=localhost in destination")
	}
}

func TestPromoteSkipsExistingWithoutOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := buildPromoteVault(t, dir, "src.vault", map[string]string{"API_KEY": "new"})
	dst := buildPromoteVault(t, dir, "dst.vault", map[string]string{"API_KEY": "old"})

	results, err := PromoteEntries(src, dst, "pass", PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected API_KEY to be skipped")
	}

	v, _ := LoadOrCreate(dst, "pass")
	if v.Entries["API_KEY"].Value != "old" {
		t.Errorf("expected destination value to remain 'old'")
	}
}

func TestPromoteOverwriteReplaces(t *testing.T) {
	dir := t.TempDir()
	src := buildPromoteVault(t, dir, "src.vault", map[string]string{"TOKEN": "newtoken"})
	dst := buildPromoteVault(t, dir, "dst.vault", map[string]string{"TOKEN": "oldtoken"})

	_, err := PromoteEntries(src, dst, "pass", PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, _ := LoadOrCreate(dst, "pass")
	if v.Entries["TOKEN"].Value != "newtoken" {
		t.Errorf("expected TOKEN to be overwritten with 'newtoken'")
	}
}

func TestPromoteDryRunDoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	src := buildPromoteVault(t, dir, "src.vault", map[string]string{"SECRET": "val"})
	dst := buildPromoteVault(t, dir, "dst.vault", map[string]string{})

	_, err := PromoteEntries(src, dst, "pass", PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, _ := LoadOrCreate(dst, "pass")
	if _, ok := v.Entries["SECRET"]; ok {
		t.Errorf("dry-run should not write SECRET to destination")
	}
}

func TestPromoteMissingKeyReported(t *testing.T) {
	dir := t.TempDir()
	src := buildPromoteVault(t, dir, "src.vault", map[string]string{"EXISTING": "val"})
	dst := buildPromoteVault(t, dir, "dst.vault", map[string]string{})

	results, err := PromoteEntries(src, dst, "pass", PromoteOptions{Keys: []string{"MISSING_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected MISSING_KEY to be reported as skipped")
	}
}

func TestFormatPromoteResultsOutput(t *testing.T) {
	f, _ := os.CreateTemp("", "promote-out")
	defer os.Remove(f.Name())

	results := []PromoteResult{
		{Key: "DB_HOST", Skipped: false},
		{Key: "OLD_KEY", Skipped: true, Reason: "key already exists"},
	}
	FormatPromoteResults(results, false, f)
	f.Seek(0, 0)
	buf := make([]byte, 256)
	n, _ := f.Read(buf)
	out := string(buf[:n])
	if out == "" {
		t.Error("expected non-empty output")
	}
}
