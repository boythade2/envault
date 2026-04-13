package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"envault/internal/vault"
)

func writeArchiveVault(t *testing.T, passphrase string) string {
	t.Helper()
	dir := t.TempDir()
	vaultFile := filepath.Join(dir, "test.vault")
	v, err := vault.LoadOrCreate(vaultFile, passphrase)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["DEPRECATED_KEY"] = vault.Entry{Value: "old"}
	v.Entries["ACTIVE_KEY"] = vault.Entry{Value: "active"}
	if err := v.Save(vaultFile, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return vaultFile
}

func TestArchiveCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "archive <vault-file> <key> [key...]" {
			return
		}
	}
	t.Fatal("archive command not registered")
}

func TestArchiveCommandRequiresTwoArgs(t *testing.T) {
	cmd := rootCmd
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"archive", "only-one-arg"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error with fewer than 2 args")
	}
}

func TestArchiveListSubcommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "archive <vault-file> <key> [key...]" {
			for _, child := range sub.Commands() {
				if child.Use == "list <vault-file>" {
					return
				}
			}
		}
	}
	t.Fatal("archive list subcommand not registered")
}

func TestArchiveEntriesIntegration(t *testing.T) {
	pass := "testpass"
	vaultFile := writeArchiveVault(t, pass)
	archived, err := vault.ArchiveEntries(vaultFile, pass, []string{"DEPRECATED_KEY"}, "old feature", false)
	if err != nil {
		t.Fatalf("ArchiveEntries: %v", err)
	}
	if len(archived) != 1 || archived[0] != "DEPRECATED_KEY" {
		t.Fatalf("unexpected result: %v", archived)
	}
	v, err := vault.LoadOrCreate(vaultFile, pass)
	if err != nil {
		t.Fatalf("reload vault: %v", err)
	}
	if _, ok := v.Entries["DEPRECATED_KEY"]; ok {
		t.Fatal("DEPRECATED_KEY should be gone from vault")
	}
	if _, ok := v.Entries["ACTIVE_KEY"]; !ok {
		t.Fatal("ACTIVE_KEY should remain")
	}
}

func TestArchiveStoreContainsReason(t *testing.T) {
	pass := "testpass"
	vaultFile := writeArchiveVault(t, pass)
	_, err := vault.ArchiveEntries(vaultFile, pass, []string{"DEPRECATED_KEY"}, "sunset", false)
	if err != nil {
		t.Fatalf("ArchiveEntries: %v", err)
	}
	store, err := vault.LoadArchive(vaultFile)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(store.Entries) == 0 {
		t.Fatal("expected archived entries")
	}
	if store.Entries[0].Reason != "sunset" {
		t.Fatalf("expected reason 'sunset', got %q", store.Entries[0].Reason)
	}
}

func TestArchiveFormatDryRun(t *testing.T) {
	out := vault.FormatArchiveResults([]string{"KEY_A", "KEY_B"}, true)
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	for _, line := range []string{"KEY_A", "KEY_B", "dry-run"} {
		if !bytes.Contains([]byte(out), []byte(line)) {
			t.Fatalf("output missing %q: %s", line, out)
		}
	}
}

func TestArchiveFileIsValidJSON(t *testing.T) {
	pass := "testpass"
	vaultFile := writeArchiveVault(t, pass)
	_, err := vault.ArchiveEntries(vaultFile, pass, []string{"DEPRECATED_KEY"}, "", false)
	if err != nil {
		t.Fatalf("ArchiveEntries: %v", err)
	}
	p := filepath.Join(filepath.Dir(vaultFile), "."+filepath.Base(vaultFile)+".archive.json")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read archive file: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("invalid JSON in archive file: %v", err)
	}
}
