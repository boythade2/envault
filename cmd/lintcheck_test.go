package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"envault/internal/vault"
)

func writeLintCheckVault(t *testing.T, entries map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.json")
	v := &vault.Vault{Entries: make(map[string]vault.Entry)}
	for k, val := range entries {
		v.Entries[k] = vault.Entry{Value: val, UpdatedAt: time.Now()}
	}
	data, _ := json.Marshal(v)
	os.WriteFile(path, data, 0600)
	return path
}

func TestLintCheckCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "lintcheck <vault-file>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("lintcheck command not registered")
	}
}

func TestLintCheckCommandRequiresArg(t *testing.T) {
	rootCmd.SetArgs([]string{"lintcheck"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestLintCheckCleanVault(t *testing.T) {
	path := writeLintCheckVault(t, map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_SECRET":   "supersecret",
	})
	out := captureOutput(t, func() {
		rootCmd.SetArgs([]string{"lintcheck", path})
		rootCmd.Execute()
	})
	if !strings.Contains(out, "No lint issues") {
		t.Errorf("expected no issues message, got: %s", out)
	}
}

func TestLintCheckFindsLowercaseKey(t *testing.T) {
	path := writeLintCheckVault(t, map[string]string{
		"lowercase_key": "value",
	})
	out := captureOutput(t, func() {
		rootCmd.SetArgs([]string{"lintcheck", path})
		rootCmd.Execute()
	})
	if !strings.Contains(out, "no-lowercase-key") {
		t.Errorf("expected no-lowercase-key finding, got: %s", out)
	}
}

func TestLintCheckLevelFilter(t *testing.T) {
	path := writeLintCheckVault(t, map[string]string{
		"lowercase_key": "",
	})
	out := captureOutput(t, func() {
		rootCmd.SetArgs([]string{"lintcheck", "--level", "error", path})
		rootCmd.Execute()
	})
	// lowercase and empty are both warn; with error filter, none should appear
	if strings.Contains(out, "no-lowercase-key") {
		t.Errorf("expected warn rules to be filtered out, got: %s", out)
	}
}
