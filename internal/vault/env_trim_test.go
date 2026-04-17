package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildTrimVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vault.json")
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "A", Value: "  hello  "},
		{Key: "B", Value: "\tworld\t"},
		{Key: "C", Value: "clean"},
	}
	if err := v.Save(p); err != nil {
		t.Fatal(err)
	}
	return v, p
}

func TestTrimAllKeys(t *testing.T) {
	v, p := buildTrimVault(t)
	results, err := TrimEntries(v, p, TrimOptions{Left: true, Right: true})
	if err != nil {
		t.Fatal(err)
	}
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}
	if changed != 2 {
		t.Fatalf("expected 2 changed, got %d", changed)
	}
	if v.Entries[0].Value != "hello" {
		t.Errorf("expected 'hello', got %q", v.Entries[0].Value)
	}
}

func TestTrimSelectedKey(t *testing.T) {
	v, p := buildTrimVault(t)
	results, err := TrimEntries(v, p, TrimOptions{Keys: []string{"A"}, Left: true, Right: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || !results[0].Changed {
		t.Fatal("expected exactly one changed result for key A")
	}
	if v.Entries[1].Value != "\tworld\t" {
		t.Error("key B should not have been modified")
	}
}

func TestTrimDryRunDoesNotWrite(t *testing.T) {
	v, p := buildTrimVault(t)
	origStat, _ := os.Stat(p)
	_, err := TrimEntries(v, p, TrimOptions{Left: true, Right: true, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	newStat, _ := os.Stat(p)
	if origStat.ModTime() != newStat.ModTime() {
		t.Error("dry run should not modify file")
	}
	if v.Entries[0].Value != "  hello  " {
		t.Error("dry run should not modify in-memory value")
	}
}

func TestFormatTrimResultsNoChange(t *testing.T) {
	v, p := buildTrimVault(t)
	v.Entries[0].Value = "clean"
	v.Entries[1].Value = "also_clean"
	v.Entries[2].Value = "clean"
	results, _ := TrimEntries(v, p, TrimOptions{Left: true, Right: true})
	out := FormatTrimResults(results, false)
	if out != "no entries required trimming\n" {
		t.Errorf("unexpected output: %s", out)
	}
}
