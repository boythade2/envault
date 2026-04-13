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

func writeFilterVault(t *testing.T, dir string) string {
	t.Helper()
	v := &vault.Vault{
		Entries: []vault.Entry{
			{Key: "DB_HOST", Value: "localhost", Tags: []string{"db"}, UpdatedAt: time.Now()},
			{Key: "DB_PORT", Value: "5432", Tags: []string{"db"}, UpdatedAt: time.Now()},
			{Key: "APP_SECRET", Value: "s3cr3t", Tags: []string{"app"}, UpdatedAt: time.Now()},
			{Key: "REDIS_URL", Value: "redis://localhost", Tags: []string{}, UpdatedAt: time.Now()},
		},
	}
	path := filepath.Join(dir, "vault.json")
	data, _ := json.Marshal(v)
	_ = os.WriteFile(path, data, 0600)
	return path
}

func TestFilterCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "filter <vault-file>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("filter command not registered")
	}
}

func TestFilterCommandRequiresArg(t *testing.T) {
	rootCmd.SetArgs([]string{"filter"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestFilterByKeyPrefix(t *testing.T) {
	dir := t.TempDir()
	path := writeFilterVault(t, dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"filter", path, "--key-prefix", "DB_"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "DB_PORT") {
		t.Errorf("expected DB entries in output, got: %q", out)
	}
	if strings.Contains(out, "APP_SECRET") {
		t.Errorf("APP_SECRET should not appear in output")
	}
}

func TestFilterByValueContains(t *testing.T) {
	dir := t.TempDir()
	path := writeFilterVault(t, dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"filter", path, "--value-contains", "localhost"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "REDIS_URL") {
		t.Errorf("expected localhost entries in output, got: %q", out)
	}
}

func TestFilterInvertFlag(t *testing.T) {
	dir := t.TempDir()
	path := writeFilterVault(t, dir)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"filter", path, "--key-prefix", "DB_", "--invert"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "DB_HOST") {
		t.Errorf("DB_HOST should not appear when inverted")
	}
	if !strings.Contains(out, "APP_SECRET") {
		t.Errorf("APP_SECRET should appear when DB_ prefix is inverted")
	}
}
