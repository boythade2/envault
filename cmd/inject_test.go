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

// writeInjectVault creates a temporary vault file with test entries for inject tests.
func writeInjectVault(t *testing.T, dir string) string {
	t.Helper()
	v := &vault.Vault{Entries: map[string]vault.Entry{}}
	v.Entries["INJECT_KEY"] = vault.Entry{Value: "hello"}
	v.Entries["INJECT_PORT"] = vault.Entry{Value: "9000"}
	path := filepath.Join(dir, "inject.vault")
	data, _ := json.Marshal(v)
	os.WriteFile(path, data, 0600)
	return path
}

func TestInjectCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "inject <vault-file> [-- command...]" {
			return
		}
	}
	t.Error("inject command not registered")
}

func TestInjectCommandRequiresArg(t *testing.T) {
	rootCmd.SetArgs([]string{"inject"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestInjectDryRun(t *testing.T) {
	dir := t.TempDir()
	path := writeInjectVault(t, dir)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"inject", path, "--dry-run"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "INJECT_KEY") {
		t.Error("expected INJECT_KEY in dry-run output")
	}
}

func TestInjectSetsEnvVars(t *testing.T) {
	dir := t.TempDir()
	path := writeInjectVault(t, dir)

	os.Unsetenv("INJECT_KEY")
	os.Unsetenv("INJECT_PORT")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"inject", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("INJECT_KEY"); got != "hello" {
		t.Errorf("expected INJECT_KEY=hello, got %q", got)
	}
}

func TestInjectWithPrefix(t *testing.T) {
	dir := t.TempDir()
	path := writeInjectVault(t, dir)

	os.Unsetenv("INJECT_KEY")
	os.Unsetenv("INJECT_PORT")

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"inject", path, "--prefix", "INJECT_KEY"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Injected: 1") {
		t.Errorf("expected 1 injected, got: %s", out)
	}
}

// TestInjectNonExistentVaultFile verifies that inject returns an error when
// the specified vault file does not exist.
func TestInjectNonExistentVaultFile(t *testing.T) {
	rootCmd.SetArgs([]string{"inject", "/nonexistent/path/to.vault"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent vault file")
	}
}
