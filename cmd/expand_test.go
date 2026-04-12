package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"envault/internal/vault"
)

func writeExpandVault(t *testing.T, entries map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.vault")
	v, err := vault.LoadOrCreate(p)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	for k, val := range entries {
		v.Entries[k] = vault.Entry{Value: val}
	}
	if err := v.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return p
}

func TestExpandCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "expand <vault-file>" {
			return
		}
	}
	t.Error("expand command not registered")
}

func TestExpandCommandRequiresArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "expand")
	if err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestExpandNoRefsOutputsMessage(t *testing.T) {
	p := writeExpandVault(t, map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	})
	out, err := executeCommand(rootCmd, "expand", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected some output")
	}
}

func TestExpandResolvesRef(t *testing.T) {
	p := writeExpandVault(t, map[string]string{
		"HOST":   "db.example.com",
		"DB_URL": "postgres://${HOST}/app",
	})
	_, err := executeCommand(rootCmd, "expand", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	var v vault.Vault
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	got := v.Entries["DB_URL"].Value
	if got != "postgres://db.example.com/app" {
		t.Errorf("got %q, want %q", got, "postgres://db.example.com/app")
	}
}

func TestExpandUnresolvedRefReturnsError(t *testing.T) {
	p := writeExpandVault(t, map[string]string{
		"VAL": "${DOES_NOT_EXIST}",
	})
	_, err := executeCommand(rootCmd, "expand", p)
	if err == nil {
		t.Error("expected error for unresolved reference")
	}
}
