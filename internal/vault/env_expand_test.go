package vault

import (
	"os"
	"testing"
)

func buildExpandVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	v, err := LoadOrCreate(dir + "/test.vault")
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	return v, dir + "/test.vault"
}

func TestExpandNoRefs(t *testing.T) {
	v, _ := buildExpandVault(t)
	v.Entries["HOST"] = Entry{Value: "localhost"}
	v.Entries["PORT"] = Entry{Value: "5432"}

	results, err := ExpandVaultRefs(v, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no changes, got changed key %q", r.Key)
		}
	}
}

func TestExpandInternalRef(t *testing.T) {
	v, _ := buildExpandVault(t)
	v.Entries["HOST"] = Entry{Value: "localhost"}
	v.Entries["DB_URL"] = Entry{Value: "postgres://${HOST}/mydb"}

	results, err := ExpandVaultRefs(v, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "DB_URL" {
			if r.Expanded != "postgres://localhost/mydb" {
				t.Errorf("got %q, want %q", r.Expanded, "postgres://localhost/mydb")
			}
			if !r.Changed {
				t.Error("expected Changed=true")
			}
		}
	}
}

func TestExpandFallsBackToOS(t *testing.T) {
	v, _ := buildExpandVault(t)
	t.Setenv("OS_VAR", "from-env")
	v.Entries["COMBINED"] = Entry{Value: "prefix-${OS_VAR}"}

	results, err := ExpandVaultRefs(v, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "COMBINED" && r.Expanded != "prefix-from-env" {
			t.Errorf("got %q, want %q", r.Expanded, "prefix-from-env")
		}
	}
}

func TestExpandUnresolvedRefReturnsError(t *testing.T) {
	v, _ := buildExpandVault(t)
	v.Entries["BROKEN"] = Entry{Value: "${MISSING_KEY}"}

	_, err := ExpandVaultRefs(v, false)
	if err == nil {
		t.Fatal("expected error for unresolved reference")
	}
}

func TestExpandNoOSFallbackWhenDisabled(t *testing.T) {
	v, _ := buildExpandVault(t)
	os.Setenv("SOME_OS_KEY", "should-not-appear")
	t.Cleanup(func() { os.Unsetenv("SOME_OS_KEY") })
	v.Entries["VAL"] = Entry{Value: "${SOME_OS_KEY}"}

	_, err := ExpandVaultRefs(v, false)
	if err == nil {
		t.Fatal("expected error when OS fallback disabled")
	}
}

func TestFormatExpandResultsNoChanges(t *testing.T) {
	results := []ExpandResult{
		{Key: "A", Original: "val", Expanded: "val", Changed: false},
	}
	out := FormatExpandResults(results)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
