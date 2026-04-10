package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envault/internal/vault"
)

func TestExportCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "export <vault-file>" {
			return
		}
	}
	t.Error("export command not registered on root")
}

func TestExportRequiresOutputFlag(t *testing.T) {
	_, err := executeCommand(rootCmd, "export", "somefile.vault")
	if err == nil {
		t.Error("Expected error when --output flag is missing")
	}
}

func TestExportDotenvFormat(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	outputPath := filepath.Join(dir, ".env")

	v, _ := vault.LoadOrCreate(vaultPath)
	v.Entries = append(v.Entries, vault.Entry{Key: "APP_ENV", Value: "production"})
	_ = v.Save(vaultPath)

	_, err := executeCommand(rootCmd, "export", vaultPath, "--output", outputPath, "--format", "dotenv")
	if err != nil {
		t.Fatalf("export command failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "APP_ENV=production") {
		t.Errorf("Expected APP_ENV=production in output, got:\n%s", string(data))
	}
}

func TestExportJSONFormat(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	outputPath := filepath.Join(dir, "env.json")

	v, _ := vault.LoadOrCreate(vaultPath)
	v.Entries = append(v.Entries, vault.Entry{Key: "SECRET", Value: "abc123"})
	_ = v.Save(vaultPath)

	_, err := executeCommand(rootCmd, "export", vaultPath, "--output", outputPath, "--format", "json")
	if err != nil {
		t.Fatalf("export command failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["SECRET"] != "abc123" {
		t.Errorf("Expected SECRET=abc123, got %s", m["SECRET"])
	}
}
