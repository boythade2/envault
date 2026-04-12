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

func writeSchemaVault(t *testing.T, entries []vault.Entry) string {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries = entries
	data, _ := json.Marshal(v)
	if err := os.WriteFile(vaultPath, data, 0600); err != nil {
		t.Fatalf("write vault: %v", err)
	}
	return vaultPath
}

func TestSchemaCommandRegistered(t *testing.T) {
	for _, use := range []string{"schema validate", "schema add-rule", "schema list"} {
		found := false
		for _, c := range rootCmd.Commands() {
			if c.Use == "schema" {
				for _, sub := range c.Commands() {
					if strings.HasPrefix(use, "schema "+sub.Name()) {
						found = true
					}
				}
			}
		}
		if !found {
			t.Errorf("command %q not registered", use)
		}
	}
}

func TestSchemaValidateClean(t *testing.T) {
	vaultPath := writeSchemaVault(t, []vault.Entry{
		{Key: "APP_HOST", Value: "localhost"},
	})
	s := vault.Schema{
		Rules: []vault.SchemaRule{
			{Key: "APP_HOST", Required: true},
		},
	}
	if err := vault.SaveSchema(vaultPath, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"schema", "validate", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "valid") {
		t.Errorf("expected valid output, got: %s", buf.String())
	}
}

func TestSchemaValidateViolation(t *testing.T) {
	vaultPath := writeSchemaVault(t, []vault.Entry{
		{Key: "APP_HOST", Value: "localhost"},
	})
	s := vault.Schema{
		Rules: []vault.SchemaRule{
			{Key: "DB_URL", Required: true},
		},
	}
	if err := vault.SaveSchema(vaultPath, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	rootCmd.SetArgs([]string{"schema", "validate", vaultPath})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for schema violation")
	}
}

func TestSchemaAddRuleAndList(t *testing.T) {
	vaultPath := writeSchemaVault(t, nil)
	rootCmd.SetArgs([]string{"schema", "add-rule", vaultPath, "SECRET_KEY",
		"--required", "--pattern", `^[A-Z]+$`, "--desc", "must be uppercase"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("add-rule: %v", err)
	}
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"schema", "list", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(buf.String(), "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", buf.String())
	}
}

func TestSchemaListEmpty(t *testing.T) {
	vaultPath := writeSchemaVault(t, nil)
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"schema", "list", vaultPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(buf.String(), "no schema rules") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
