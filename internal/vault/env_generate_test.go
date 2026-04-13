package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildGenerateVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "vault.json")
}

func TestGenerateSecretLength(t *testing.T) {
	secret, err := GenerateSecret(GenerateOptions{Length: 24, UseDigits: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secret) != 24 {
		t.Errorf("expected length 24, got %d", len(secret))
	}
}

func TestGenerateSecretZeroLengthReturnsError(t *testing.T) {
	_, err := GenerateSecret(GenerateOptions{Length: 0})
	if err == nil {
		t.Fatal("expected error for zero length")
	}
}

func TestGenerateSecretUsesUpperAndDigits(t *testing.T) {
	// Generate a long string to ensure all charsets appear with high probability.
	secret, err := GenerateSecret(GenerateOptions{Length: 200, UseUpper: true, UseDigits: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hasUpper := strings.ContainsAny(secret, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasDigit := strings.ContainsAny(secret, "0123456789")
	if !hasUpper {
		t.Error("expected uppercase characters in generated secret")
	}
	if !hasDigit {
		t.Error("expected digit characters in generated secret")
	}
}

func TestGenerateSecretIsUnique(t *testing.T) {
	a, _ := GenerateSecret(GenerateOptions{Length: 32, UseDigits: true})
	b, _ := GenerateSecret(GenerateOptions{Length: 32, UseDigits: true})
	if a == b {
		t.Error("two generated secrets should not be identical")
	}
}

func TestGenerateAndStoreWritesVault(t *testing.T) {
	vaultFile := buildGenerateVault(t)
	opts := GenerateOptions{Length: 16, UseDigits: true}
	r, err := GenerateAndStore(vaultFile, "MY_SECRET", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Key != "MY_SECRET" {
		t.Errorf("expected key MY_SECRET, got %s", r.Key)
	}
	if len(r.Value) != 16 {
		t.Errorf("expected value length 16, got %d", len(r.Value))
	}
	v, err := LoadOrCreate(vaultFile)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	entry, ok := v.Entries["MY_SECRET"]
	if !ok {
		t.Fatal("key MY_SECRET not found in vault after generate")
	}
	if entry.Value != r.Value {
		t.Errorf("stored value mismatch: got %s, want %s", entry.Value, r.Value)
	}
}

func TestGenerateAndStoreDryRunDoesNotWrite(t *testing.T) {
	vaultFile := buildGenerateVault(t)
	opts := GenerateOptions{Length: 12, DryRun: true}
	_, err := GenerateAndStore(vaultFile, "DRY_KEY", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(vaultFile); !os.IsNotExist(statErr) {
		t.Error("vault file should not exist after dry-run generate")
	}
}

func TestFormatGenerateResult(t *testing.T) {
	r := GenerateResult{Key: "TOKEN", Value: "abc123", Created: true}
	out := FormatGenerateResult(r, false)
	if !strings.Contains(out, "TOKEN") || !strings.Contains(out, "abc123") {
		t.Errorf("unexpected format output: %s", out)
	}
	dryOut := FormatGenerateResult(r, true)
	if !strings.HasPrefix(dryOut, "[dry-run]") {
		t.Errorf("expected dry-run prefix, got: %s", dryOut)
	}
}
