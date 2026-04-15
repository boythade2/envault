package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildSquashVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vf := filepath.Join(dir, "vault.json")
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "APP_NAME", Value: "  myapp  "},
		{Key: "APP_ENV", Value: "Production"},
		{Key: "APP_DESC", Value: "hello   world"},
		{Key: "APP_CODE", Value: "ABC123"},
	}
	if err := v.Save(vf); err != nil {
		t.Fatalf("save: %v", err)
	}
	return vf
}

func TestSquashTrim(t *testing.T) {
	vf := buildSquashVault(t)
	results, err := SquashEntries(vf, SquashOptions{Keys: []string{"APP_NAME"}, Transform: "trim"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].NewVal != "myapp" {
		t.Errorf("expected 'myapp', got %q", results[0].NewVal)
	}
}

func TestSquashLower(t *testing.T) {
	vf := buildSquashVault(t)
	results, err := SquashEntries(vf, SquashOptions{Keys: []string{"APP_ENV"}, Transform: "lower"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].NewVal != "production" {
		t.Errorf("expected 'production', got %q", results[0].NewVal)
	}
}

func TestSquashUpper(t *testing.T) {
	vf := buildSquashVault(t)
	results, err := SquashEntries(vf, SquashOptions{Keys: []string{"APP_ENV"}, Transform: "upper"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].NewVal != "PRODUCTION" {
		t.Errorf("expected 'PRODUCTION', got %q", results[0].NewVal)
	}
}

func TestSquashCollapse(t *testing.T) {
	vf := buildSquashVault(t)
	results, err := SquashEntries(vf, SquashOptions{Keys: []string{"APP_DESC"}, Transform: "collapse"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].NewVal != "hello world" {
		t.Errorf("expected 'hello world', got %q", results[0].NewVal)
	}
}

func TestSquashDryRunDoesNotModify(t *testing.T) {
	vf := buildSquashVault(t)
	origData, _ := os.ReadFile(vf)
	_, err := SquashEntries(vf, SquashOptions{Keys: []string{"APP_NAME"}, Transform: "trim", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	afterData, _ := os.ReadFile(vf)
	if string(origData) != string(afterData) {
		t.Error("dry run must not modify the vault file")
	}
}

func TestSquashNoChangeSkipped(t *testing.T) {
	vf := buildSquashVault(t)
	// APP_CODE has no whitespace or case change for 'collapse'
	results, err := SquashEntries(vf, SquashOptions{Keys: []string{"APP_CODE"}, Transform: "collapse"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Skipped {
		t.Error("expected entry to be skipped when value unchanged")
	}
}

func TestSquashEmptyTransformReturnsError(t *testing.T) {
	vf := buildSquashVault(t)
	_, err := SquashEntries(Transform: ""})
	if err == nil {
		t.Error("expected error for empty transform")
	}
}
