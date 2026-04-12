package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envault/internal/vault"
)

func writeWatchVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "watch.vault")
	v, err := vault.LoadOrCreate(vaultPath, "secret")
	if err != nil {
		t.Fatalf("create vault: %v", err)
	}
	v.Entries["FOO"] = vault.Entry{Value: "bar"}
	if err := v.Save(vaultPath, "secret"); err != nil {
		t.Fatalf("save vault: %v", err)
	}
	return vaultPath
}

func TestWatchCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "watch <vault-file>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("watch command not registered")
	}
}

func TestWatchRecordRequiresArg(t *testing.T) {
	rootCmd.SetArgs([]string{"watch", "record"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no vault path given")
	}
}

func TestWatchRecordCreatesState(t *testing.T) {
	vaultPath := writeWatchVault(t)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"watch", "record", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "recorded") {
		t.Errorf("expected 'recorded' in output, got: %s", buf.String())
	}
	state, err := vault.LoadWatchState(vaultPath)
	if err != nil || state == nil {
		t.Error("expected watch state to be saved on disk")
	}
}

func TestWatchStatusUnchanged(t *testing.T) {
	vaultPath := writeWatchVault(t)
	if err := vault.SaveWatchState(vaultPath); err != nil {
		t.Fatalf("save watch state: %v", err)
	}
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"watch", "status", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "UNCHANGED") {
		t.Errorf("expected UNCHANGED, got: %s", buf.String())
	}
}

func TestWatchStatusChanged(t *testing.T) {
	vaultPath := writeWatchVault(t)
	if err := vault.SaveWatchState(vaultPath); err != nil {
		t.Fatalf("save watch state: %v", err)
	}
	f, _ := os.OpenFile(vaultPath, os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString("\n")
	f.Close()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"watch", "status", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "CHANGED") {
		t.Errorf("expected CHANGED, got: %s", buf.String())
	}
}
