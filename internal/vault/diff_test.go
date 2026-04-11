package vault

import (
	"testing"
)

func buildVault(entries map[string]string) *Vault {
	v := &Vault{}
	for k, val := range entries {
		v.Entries = append(v.Entries, Entry{Key: k, Value: val})
	}
	return v
}

func TestDiffIdenticalVaults(t *testing.T) {
	a := buildVault(map[string]string{"FOO": "bar", "BAZ": "qux"})
	b := buildVault(map[string]string{"FOO": "bar", "BAZ": "qux"})

	result := Diff(a, b)

	if result.HasChanges() {
		t.Errorf("expected no changes, got added=%v removed=%v changed=%v", result.Added, result.Removed, result.Changed)
	}
	if len(result.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged keys, got %d", len(result.Unchanged))
	}
}

func TestDiffAdded(t *testing.T) {
	a := buildVault(map[string]string{"FOO": "bar"})
	b := buildVault(map[string]string{"FOO": "bar", "NEW_KEY": "value"})

	result := Diff(a, b)

	if len(result.Added) != 1 || result.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in Added, got %v", result.Added)
	}
}

func TestDiffRemoved(t *testing.T) {
	a := buildVault(map[string]string{"FOO": "bar", "OLD_KEY": "value"})
	b := buildVault(map[string]string{"FOO": "bar"})

	result := Diff(a, b)

	if len(result.Removed) != 1 || result.Removed[0] != "OLD_KEY" {
		t.Errorf("expected OLD_KEY in Removed, got %v", result.Removed)
	}
}

func TestDiffChanged(t *testing.T) {
	a := buildVault(map[string]string{"FOO": "old_value"})
	b := buildVault(map[string]string{"FOO": "new_value"})

	result := Diff(a, b)

	if len(result.Changed) != 1 || result.Changed[0] != "FOO" {
		t.Errorf("expected FOO in Changed, got %v", result.Changed)
	}
}

func TestDiffSortedOutput(t *testing.T) {
	a := buildVault(map[string]string{})
	b := buildVault(map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"})

	result := Diff(a, b)

	if len(result.Added) != 3 {
		t.Fatalf("expected 3 added, got %d", len(result.Added))
	}
	if result.Added[0] != "APPLE" || result.Added[1] != "MANGO" || result.Added[2] != "ZEBRA" {
		t.Errorf("expected sorted order, got %v", result.Added)
	}
}
