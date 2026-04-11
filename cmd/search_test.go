package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"envault/internal/vault"
)

func writeSearchVault(t *testing.T, entries map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envault")

	v := &vault.Vault{
		Entries:   make(map[string]vault.Entry),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for k, val := range entries {
		v.Entries[k] = vault.Entry{Key: k, Value: val, UpdatedAt: time.Now()}
	}

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSearchCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "search <query>" {
			return
		}
	}
	t.Fatal("search command not registered")
}

func TestSearchCommandRequiresArg(t *testing.T) {
	rootCmd.SetArgs([]string{"search"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no query arg provided")
	}
}

func TestSearchByKeyPartialMatch(t *testing.T) {
	path := writeSearchVault(t, map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"REDIS_URL":    "redis://localhost",
		"APP_SECRET":   "supersecret",
	})

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"search", "URL", "--file", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected DATABASE_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "REDIS_URL") {
		t.Errorf("expected REDIS_URL in output, got: %s", out)
	}
	if strings.Contains(out, "APP_SECRET") {
		t.Errorf("did not expect APP_SECRET in output, got: %s", out)
	}
}

func TestSearchNoResults(t *testing.T) {
	path := writeSearchVault(t, map[string]string{
		"APP_SECRET": "supersecret",
	})

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"search", "NONEXISTENT", "--file", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No matching entries found") {
		t.Errorf("expected no-results message, got: %s", buf.String())
	}
}
