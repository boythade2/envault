package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildCastVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "cast.vault")
	passphrase := "castpass"
	v, err := LoadOrCreate(vaultPath, passphrase)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries = []Entry{
		{Key: "PORT", Value: "8080.0"},
		{Key: "DEBUG", Value: "True"},
		{Key: "RATIO", Value: "3"},
		{Key: "NAME", Value: "envault"},
	}
	if err := v.Save(vaultPath, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return vaultPath, passphrase
}

func TestCastToInt(t *testing.T) {
	vaultPath, pass := buildCastVault(t)
	results, err := CastEntries(vaultPath, pass, []string{"PORT"}, CastInt, false)
	if err != nil {
		t.Fatalf("CastEntries: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].NewVal != "8080" {
		t.Errorf("expected 8080, got %s", results[0].NewVal)
	}
	if !results[0].Casted {
		t.Error("expected Casted=true")
	}
}

func TestCastToBool(t *testing.T) {
	vaultPath, pass := buildCastVault(t)
	results, err := CastEntries(vaultPath, pass, []string{"DEBUG"}, CastBool, false)
	if err != nil {
		t.Fatalf("CastEntries: %v", err)
	}
	if results[0].NewVal != "true" {
		t.Errorf("expected true, got %s", results[0].NewVal)
	}
}

func TestCastToFloat(t *testing.T) {
	vaultPath, pass := buildCastVault(t)
	results, err := CastEntries(vaultPath, pass, []string{"RATIO"}, CastFloat, false)
	if err != nil {
		t.Fatalf("CastEntries: %v", err)
	}
	if results[0].NewVal != "3" {
		t.Errorf("expected 3, got %s", results[0].NewVal)
	}
}

func TestCastDryRunDoesNotWrite(t *testing.T) {
	vaultPath, pass := buildCastVault(t)
	stat1, _ := os.Stat(vaultPath)
	_, err := CastEntries(vaultPath, pass, []string{"PORT"}, CastInt, true)
	if err != nil {
		t.Fatalf("CastEntries dry-run: %v", err)
	}
	stat2, _ := os.Stat(vaultPath)
	if stat1.ModTime() != stat2.ModTime() {
		t.Error("dry-run should not modify the vault file")
	}
}

func TestCastInvalidValueReturnsError(t *testing.T) {
	vaultPath, pass := buildCastVault(t)
	results, err := CastEntries(vaultPath, pass, []string{"NAME"}, CastInt, false)
	if err != nil {
		t.Fatalf("CastEntries: %v", err)
	}
	if results[0].Err == "" {
		t.Error("expected error for non-numeric value")
	}
}

func TestCastNoMatchingKeys(t *testing.T) {
	vaultPath, pass := buildCastVault(t)
	results, err := CastEntries(vaultPath, pass, []string{"NONEXISTENT"}, CastBool, false)
	if err != nil {
		t.Fatalf("CastEntries: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFormatCastResultsEmpty(t *testing.T) {
	out := FormatCastResults(nil)
	if out != "no matching keys found\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
