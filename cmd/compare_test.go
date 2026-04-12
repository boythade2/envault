package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"envault/internal/vault"
)

func writeCompareVault(t *testing.T, path, pass string, entries map[string]string) {
	t.Helper()
	v, err := vault.LoadOrCreate(path, pass)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	for k, val := range entries {
		v.Entries = append(v.Entries, vault.Entry{Key: k, Value: val})
	}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestCompareCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "compare <vaultA> <vaultB>" {
			return
		}
	}
	t.Error("compare command not registered")
}

func TestCompareCommandRequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "compare", "only-one")
	if err == nil {
		t.Error("expected error with one arg")
	}
}

func TestCompareIdenticalVaultsOutput(t *testing.T) {
	dir := t.TempDir()
	pA := filepath.Join(dir, "a.vault")
	pB := filepath.Join(dir, "b.vault")
	writeCompareVault(t, pA, "secret", map[string]string{"FOO": "bar"})
	writeCompareVault(t, pB, "secret", map[string]string{"FOO": "bar"})

	out, err := executeCommand(rootCmd, "compare", pA, pB, "--pass-a", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains([]byte(out), []byte("identical")) {
		t.Errorf("expected 'identical' in output, got: %s", out)
	}
}

func TestCompareDifferentVaultsOutput(t *testing.T) {
	dir := t.TempDir()
	pA := filepath.Join(dir, "a.vault")
	pB := filepath.Join(dir, "b.vault")
	writeCompareVault(t, pA, "secret", map[string]string{"FOO": "bar", "ONLY_A": "x"})
	writeCompareVault(t, pB, "secret", map[string]string{"FOO": "changed", "ONLY_B": "y"})

	out, err := executeCommand(rootCmd, "compare", pA, pB, "--pass-a", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"Only in A", "Only in B", "Changed"} {
		if !bytes.Contains([]byte(out), []byte(want)) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}
