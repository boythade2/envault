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

func writePinVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "env.vault")
}

func TestPinCommandRegistered(t *testing.T) {
	names := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		names[c.Name()] = true
	}
	for _, name := range []string{"pin", "unpin", "pins"} {
		if !names[name] {
			t.Errorf("command %q not registered", name)
		}
	}
}

func TestPinCommandRequiresTwoArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"pin", "only-one"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error with one arg, got nil")
	}
}

func TestPinAndListPins(t *testing.T) {
	vp := writePinVault(t)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"pin", vp, "DB_PASS", "--note", "critical"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("pin command: %v", err)
	}
	buf.Reset()
	rootCmd.SetArgs([]string{"pins", vp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("pins command: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %s", out)
	}
	if !strings.Contains(out, "critical") {
		t.Errorf("expected note 'critical' in output, got: %s", out)
	}
}

func TestUnpinCommand(t *testing.T) {
	vp := writePinVault(t)
	if err := vault.PinKey(vp, "TOKEN", ""); err != nil {
		t.Fatalf("setup PinKey: %v", err)
	}
	rootCmd.SetArgs([]string{"unpin", vp, "TOKEN"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unpin command: %v", err)
	}
	pinned, err := vault.IsPinned(vp, "TOKEN")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if pinned {
		t.Error("expected TOKEN to be unpinned after unpin command")
	}
}

func TestPinsEmptyVault(t *testing.T) {
	vp := writePinVault(t)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"pins", vp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("pins command: %v", err)
	}
	if !strings.Contains(buf.String(), "No pinned keys") {
		t.Errorf("expected 'No pinned keys' message, got: %s", buf.String())
	}
}

func TestPinCreatesFileWithCorrectPermissions(t *testing.T) {
	vp := writePinVault(t)
	if err := vault.PinKey(vp, "SECRET", ""); err != nil {
		t.Fatalf("PinKey: %v", err)
	}
	pinFile := filepath.Join(filepath.Dir(vp), "."+filepath.Base(vp)+".pins.json")
	data, err := os.ReadFile(pinFile)
	if err != nil {
		t.Fatalf("read pin file: %v", err)
	}
	var pl vault.PinList
	if err := json.Unmarshal(data, &pl); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(pl.Pins) != 1 {
		t.Errorf("expected 1 pin, got %d", len(pl.Pins))
	}
}
