package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"envault/internal/vault"
)

func TestRotateCommandRegistered(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "rotate <vault-file>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("rotate command not registered on root command")
	}
}

func TestRotateCommandRequiresArg(t *testing.T) {
	rootCmd.SetArgs([]string{"rotate"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no vault file argument provided")
	}
}

func TestRotatePreservesEntries(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	oldPassphrase := "old-secret"
	newPassphrase := "new-secret"

	v, err := vault.LoadOrCreate(vaultPath, oldPassphrase)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}
	v.Add("KEY1", "value1")
	v.Add("KEY2", "value2")
	if err := v.Save(vaultPath, oldPassphrase); err != nil {
		t.Fatalf("failed to save vault: %v", err)
	}

	if err := v.Rotate(vaultPath, newPassphrase); err != nil {
		t.Fatalf("rotate failed: %v", err)
	}

	v2, err := vault.LoadOrCreate(vaultPath, newPassphrase)
	if err != nil {
		t.Fatalf("failed to load rotated vault: %v", err)
	}

	if val, ok := v2.Get("KEY1"); !ok || val != "value1" {
		t.Errorf("expected KEY1=value1 after rotation, got %q", val)
	}
	if val, ok := v2.Get("KEY2"); !ok || val != "value2" {
		t.Errorf("expected KEY2=value2 after rotation, got %q", val)
	}
}

func TestRotateUpdatesFile(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	oldPassphrase := "old-secret"
	newPassphrase := "new-secret"

	v, err := vault.LoadOrCreate(vaultPath, oldPassphrase)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}
	v.Add("FOO", "bar")
	if err := v.Save(vaultPath, oldPassphrase); err != nil {
		t.Fatalf("failed to save vault: %v", err)
	}

	info1, _ := os.Stat(vaultPath)
	originalModTime := info1.ModTime()

	time.Sleep(10 * time.Millisecond)

	if err := v.Rotate(vaultPath, newPassphrase); err != nil {
		t.Fatalf("rotate failed: %v", err)
	}

	data, err := os.ReadFile(vaultPath)
	if err != nil {
		t.Fatalf("failed to read rotated vault file: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err == nil {
		t.Error("rotated vault file should not be plain JSON")
	}

	info2, _ := os.Stat(vaultPath)
	if !info2.ModTime().After(originalModTime) {
		t.Error("expected vault file modification time to be updated after rotation")
	}
}
