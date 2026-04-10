package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "add <vault-file> <KEY> <VALUE>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("add command not registered")
	}
}

func TestRemoveCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "remove <vault-file> <KEY>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("remove command not registered")
	}
}

func TestAddAndRemoveEntry(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	passphrase := "test-passphrase"

	// Add a key
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"add", vaultPath, "MY_KEY", "my_value", "--passphrase", passphrase})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		t.Fatal("vault file was not created")
	}

	// Remove the key
	rootCmd.SetArgs([]string{"remove", vaultPath, "MY_KEY", "--passphrase", passphrase})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("remove command failed: %v", err)
	}
}

func TestRemoveNonExistentKey(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	passphrase := "test-passphrase"

	// Create vault first with a key
	rootCmd.SetArgs([]string{"add", vaultPath, "EXISTING_KEY", "value", "--passphrase", passphrase})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	// Try to remove a key that doesn't exist
	errBuf := &bytes.Buffer{}
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs([]string{"remove", vaultPath, "MISSING_KEY", "--passphrase", passphrase})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when removing non-existent key")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}
