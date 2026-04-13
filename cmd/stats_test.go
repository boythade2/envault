package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

func writeStatsVault(t *testing.T, dir string) string {
	t.Helper()
	v := &vault.Vault{
		Entries: []vault.Entry{
			{Key: "DB_HOST", Value: "localhost", UpdatedAt: time.Now()},
			{Key: "DB_PORT", Value: "5432", UpdatedAt: time.Now()},
			{Key: "DB_NAME", Value: "mydb", UpdatedAt: time.Now()},
			{Key: "APP_ENV", Value: "production", UpdatedAt: time.Now()},
			{Key: "APP_DEBUG", Value: "", UpdatedAt: time.Now()},
			{Key: "APP_SECRET", Value: "s3cr3t", UpdatedAt: time.Now()},
		},
	}
	path := filepath.Join(dir, "test.vault")
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal vault: %v", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write vault: %v", err)
	}
	return path
}

func TestStatsCommandRegistered(t *testing.T) {
	var found bool
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "stats" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'stats' command to be registered")
	}
}

func TestStatsCommandRequiresArg(t *testing.T) {
	cmd := &cobra.Command{Use: "stats"}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runStats(cmd, args)
	}
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when no vault file provided")
	}
}

func TestStatsOutputContainsTotalKeys(t *testing.T) {
	dir := t.TempDir()
	path := writeStatsVault(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"stats", path})
	if err := rootCmd.Execute(); err != nil {
		w.Close()
		os.Stdout = old
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "6") {
		t.Errorf("expected output to contain total key count 6, got:\n%s", out)
	}
}

func TestStatsOutputShowsEmptyCount(t *testing.T) {
	dir := t.TempDir()
	path := writeStatsVault(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"stats", path})
	if err := rootCmd.Execute(); err != nil {
		w.Close()
		os.Stdout = old
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	// APP_DEBUG has an empty value
	if !strings.Contains(out, "Empty") && !strings.Contains(out, "empty") {
		t.Errorf("expected output to mention empty values, got:\n%s", out)
	}
}

func TestStatsOutputShowsTopPrefixes(t *testing.T) {
	dir := t.TempDir()
	path := writeStatsVault(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"stats", path})
	if err := rootCmd.Execute(); err != nil {
		w.Close()
		os.Stdout = old
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	// Both APP_ and DB_ prefixes appear in the vault
	if !strings.Contains(out, "APP") || !strings.Contains(out, "DB") {
		t.Errorf("expected output to show top prefixes APP and DB, got:\n%s", out)
	}
}
