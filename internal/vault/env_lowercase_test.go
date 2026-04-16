package vault

import (
	"path/filepath"
	"testing"
)

func buildLowercaseVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.env")
	pass := "testpass"
	v, err := LoadOrCreate(path, pass)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["MY_KEY"] = Entry{Value: "hello"}
	v.Entries["ANOTHER_VAR"] = Entry{Value: "world"}
	v.Entries["already_lower"] = Entry{Value: "fine"}
	v.Entries["CONFLICT"] = Entry{Value: "a"}
	v.Entries["conflict"] = Entry{Value: "b"}
	if err := v.Save(path, pass); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path, pass
}

func TestLowercaseAllKeys(t *testing.T) {
	path, pass := buildLowercaseVault(t)
	results, err := LowercaseKeys(path, pass, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	v, _ := LoadOrCreate(path, pass)
	if _, ok := v.Entries["MY_KEY"]; ok {
		t.Error("MY_KEY should have been renamed")
	}
	if _, ok := v.Entries["my_key"]; !ok {
		t.Error("my_key should exist after rename")
	}
}

func TestLowercaseSelectedKey(t *testing.T) {
	path, pass := buildLowercaseVault(t)
	_, err := LowercaseKeys(path, pass, []string{"ANOTHER_VAR"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := LoadOrCreate(path, pass)
	if _, ok := v.Entries["another_var"]; !ok {
		t.Error("another_var should exist")
	}
	if _, ok := v.Entries["MY_KEY"]; !ok {
		t.Error("MY_KEY should be untouched")
	}
}

func TestLowercaseDryRunDoesNotWrite(t *testing.T) {
	path, pass := buildLowercaseVault(t)
	_, err := LowercaseKeys(path, pass, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := LoadOrCreate(path, pass)
	if _, ok := v.Entries["MY_KEY"]; !ok {
		t.Error("dry run should not rename MY_KEY")
	}
}

func TestLowercaseSkipsConflict(t *testing.T) {
	path, pass := buildLowercaseVault(t)
	results, err := LowercaseKeys(path, pass, []string{"CONFLICT"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var skipped bool
	for _, r := range results {
		if r.OldKey == "CONFLICT" && r.Skipped {
			skipped = true
		}
	}
	if !skipped {
		t.Error("expected CONFLICT to be skipped due to existing lowercase key")
	}
}

func TestFormatLowercaseResultsEmpty(t *testing.T) {
	out := FormatLowercaseResults(nil)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
