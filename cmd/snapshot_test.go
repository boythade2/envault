package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envault/internal/vault"
)

func writeSnapshotVault(t *testing.T, dir string) string {
	t.Helper()
	vaultPath := filepath.Join(dir, "test.vault")
	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["KEY1"] = vault.Entry{Value: "val1"}
	v.Entries["KEY2"] = vault.Entry{Value: "val2"}
	data, _ := json.Marshal(v)
	if err := os.WriteFile(vaultPath, data, 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return vaultPath
}

func TestSnapshotCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "snapshot <vault-file>" {
			return
		}
	}
	t.Fatal("snapshot command not registered")
}

func TestSnapshotCommandRequiresArg(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no arg provided")
	}
}

func TestSnapshotSavesFile(t *testing.T) {
	dir := t.TempDir()
	vaultPath := writeSnapshotVault(t, dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot", "--label", "before-deploy", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("snapshot command failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Snapshot saved:") {
		t.Errorf("expected 'Snapshot saved:' in output, got: %s", out)
	}
}

func TestSnapshotListShowsEntries(t *testing.T) {
	dir := t.TempDir()
	vaultPath := writeSnapshotVault(t, dir)

	// Save a snapshot first
	rootCmd.SetArgs([]string{"snapshot", "--label", "v1", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("snapshot save failed: %v", err)
	}

	// Now list snapshots
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot", "list", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("snapshot list failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "NAME") {
		t.Errorf("expected header NAME in output, got: %s", out)
	}
}

func TestSnapshotListEmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	vaultPath := writeSnapshotVault(t, dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot", "list", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("snapshot list failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No snapshots found.") {
		t.Errorf("expected 'No snapshots found.' in output, got: %s", out)
	}
}
