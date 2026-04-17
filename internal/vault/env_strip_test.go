package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildStripVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "strip.vault")
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "KEY_A", Value: "  hello  "},
		{Key: "KEY_B", Value: "\tworld\t"},
		{Key: "KEY_C", Value: "clean"},
	}
	if err := v.Save(p); err != nil {
		t.Fatalf("save: %v", err)
	}
	return p
}

func TestStripAllKeys(t *testing.T) {
	p := buildStripVault(t)
	results, err := StripEntries(p, StripOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if !results[0].Changed || results[0].NewVal != "hello" {
		t.Errorf("KEY_A: expected stripped value 'hello', got %q", results[0].NewVal)
	}
	if !results[1].Changed || results[1].NewVal != "world" {
		t.Errorf("KEY_B: expected stripped value 'world', got %q", results[1].NewVal)
	}
	if results[2].Changed {
		t.Errorf("KEY_C should not have changed")
	}
}

func TestStripSelectedKey(t *testing.T) {
	p := buildStripVault(t)
	results, err := StripEntries(p, StripOptions{Keys: []string{"KEY_A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "KEY_A" || !results[0].Changed {
		t.Errorf("expected KEY_A to be stripped")
	}
}

func TestStripDryRunDoesNotWrite(t *testing.T) {
	p := buildStripVault(t)
	stat1, _ := os.Stat(p)
	_, err := StripEntries(p, StripOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stat2, _ := os.Stat(p)
	if stat1.ModTime() != stat2.ModTime() {
		t.Errorf("dry-run should not modify the file")
	}
}

func TestStripCustomChars(t *testing.T) {
	p := buildStripVault(t)
	v, _ := LoadOrCreate(p)
	v.Entries = []Entry{{Key: "K", Value: "***secret***"}}
	v.Save(p)

	results, err := StripEntries(p, StripOptions{Chars: "*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || results[0].NewVal != "secret" {
		t.Errorf("expected 'secret', got %q", results[0].NewVal)
	}
}

func TestFormatStripResultsNoChange(t *testing.T) {
	out := FormatStripResults([]StripResult{{Key: "K", Changed: false}})
	if out != "no values required stripping" {
		t.Errorf("unexpected output: %s", out)
	}
}
