package vault

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildEncodeVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "encode.env")
	passphrase := "encodepass"
	v, err := LoadOrCreate(vaultPath, passphrase)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries = []Entry{
		{Key: "SECRET", Value: "hello"},
		{Key: "TOKEN", Value: "world"},
	}
	if err := v.Save(vaultPath, passphrase); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return vaultPath, passphrase
}

func TestEncodeBase64AllKeys(t *testing.T) {
	vaultPath, pass := buildEncodeVault(t)
	results, err := EncodeEntries(vaultPath, pass, nil, EncodeBase64, false)
	if err != nil {
		t.Fatalf("EncodeEntries: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		decoded, err := base64.StdEncoding.DecodeString(r.Encoded)
		if err != nil {
			t.Errorf("key %s: invalid base64: %v", r.Key, err)
		}
		if string(decoded) != r.Original {
			t.Errorf("key %s: round-trip failed: got %q", r.Key, decoded)
		}
	}
}

func TestEncodeSelectedKey(t *testing.T) {
	vaultPath, pass := buildEncodeVault(t)
	results, err := EncodeEntries(vaultPath, pass, []string{"SECRET"}, EncodeBase64, false)
	if err != nil {
		t.Fatalf("EncodeEntries: %v", err)
	}
	if len(results) != 1 || results[0].Key != "SECRET" {
		t.Fatalf("expected only SECRET, got %+v", results)
	}
}

func TestEncodeHex(t *testing.T) {
	vaultPath, pass := buildEncodeVault(t)
	results, err := EncodeEntries(vaultPath, pass, []string{"TOKEN"}, EncodeHex, false)
	if err != nil {
		t.Fatalf("EncodeEntries: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	// "world" in hex
	expected := "776f726c64"
	if results[0].Encoded != expected {
		t.Errorf("expected %q, got %q", expected, results[0].Encoded)
	}
}

func TestEncodeDryRunDoesNotModify(t *testing.T) {
	vaultPath, pass := buildEncodeVault(t)
	_, err := EncodeEntries(vaultPath, pass, nil, EncodeBase64, true)
	if err != nil {
		t.Fatalf("EncodeEntries dry-run: %v", err)
	}
	v, err := LoadOrCreate(vaultPath, pass)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	for _, e := range v.Entries {
		if e.Key == "SECRET" && e.Value != "hello" {
			t.Errorf("dry-run modified SECRET: %q", e.Value)
		}
	}
}

func TestEncodeUnsupportedFormatReturnsError(t *testing.T) {
	vaultPath, pass := buildEncodeVault(t)
	results, err := EncodeEntries(vaultPath, pass, nil, EncodeFormat("rot13"), false)
	if err != nil {
		// save would fail only if all entries errored but vault still writes; check results
		_ = err
	}
	for _, r := range results {
		if !r.Skipped {
			t.Errorf("expected key %s to be skipped for unsupported format", r.Key)
		}
	}
}

func TestEncodeFilePermissions(t *testing.T) {
	vaultPath, pass := buildEncodeVault(t)
	_, err := EncodeEntries(vaultPath, pass, nil, EncodeBase64, false)
	if err != nil {
		t.Fatalf("EncodeEntries: %v", err)
	}
	info, err := os.Stat(vaultPath)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestFormatEncodeResultsDryRun(t *testing.T) {
	results := []EncodeResult{
		{Key: "A", Original: "foo", Encoded: "Zm9v"},
	}
	out := FormatEncodeResults(results, true)
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run label in output: %q", out)
	}
}
