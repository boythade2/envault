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

func writeCastVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	pass := "testpass"
	v, err := vault.LoadOrCreate(vaultPath, pass)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries = []vault.Entry{
		{Key: "PORT", Value: "8080.0"},
		{Key: "DEBUG", Value: "True"},
	}
	if err := v.Save(vaultPath, pass); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return vaultPath, pass
}

func TestCastCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "cast <vault-file> <key>[,key...]" {
			return
		}
	}
	t.Error("cast command not registered")
}

func TestCastCommandRequiresTwoArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"cast", "only-one-arg"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error with single arg")
	}
}

func TestCastToInt(t *testing.T) {
	vaultPath, pass := writeCastVault(t)
	os.Setenv("ENVAULT_PASSPHRASE", pass)
	defer os.Unsetenv("ENVAULT_PASSPHRASE")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"cast", vaultPath, "PORT", "--type", "int"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	v, err := vault.LoadOrCreate(vaultPath, pass)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	for _, e := range v.Entries {
		if e.Key == "PORT" && e.Value != "8080" {
			t.Errorf("expected PORT=8080, got %s", e.Value)
		}
	}
}

func TestCastDryRunOutputContainsPreview(t *testing.T) {
	vaultPath, pass := writeCastVault(t)
	os.Setenv("ENVAULT_PASSPHRASE", pass)
	defer os.Unsetenv("ENVAULT_PASSPHRASE")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"cast", vaultPath, "DEBUG", "--type", "bool", "--dry-run"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "dry-run") {
		t.Errorf("expected dry-run notice in output, got: %s", output)
	}
}

func TestCastUnsupportedTypeReturnsError(t *testing.T) {
	vaultPath, pass := writeCastVault(t)
	os.Setenv("ENVAULT_PASSPHRASE", pass)
	defer os.Unsetenv("ENVAULT_PASSPHRASE")

	rootCmd.SetArgs([]string{"cast", vaultPath, "PORT", "--type", "xml"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

// Ensure JSON marshalling of CastResult works (used by potential --json flag).
func TestCastResultMarshal(t *testing.T) {
	r := vault.CastResult{Key: "PORT", OldVal: "8080.0", NewVal: "8080", Casted: true}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), "PORT") {
		t.Error("marshalled JSON missing key")
	}
}
