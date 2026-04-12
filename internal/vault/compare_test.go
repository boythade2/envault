package vault

import (
	"testing"
)

func buildCompareVault(entries map[string]string) *Vault {
	v := &Vault{}
	for k, val := range entries {
		v.Entries = append(v.Entries, Entry{Key: k, Value: val})
	}
	return v
}

func TestCompareIdenticalVaults(t *testing.T) {
	a := buildCompareVault(map[string]string{"FOO": "bar", "BAZ": "qux"})
	b := buildCompareVault(map[string]string{"FOO": "bar", "BAZ": "qux"})
	r := CompareVaults(a, b)
	if len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 || len(r.Changed) != 0 {
		t.Errorf("expected no differences, got %+v", r)
	}
	if len(r.Identical) != 2 {
		t.Errorf("expected 2 identical keys, got %d", len(r.Identical))
	}
}

func TestCompareOnlyInA(t *testing.T) {
	a := buildCompareVault(map[string]string{"FOO": "bar", "EXTRA": "val"})
	b := buildCompareVault(map[string]string{"FOO": "bar"})
	r := CompareVaults(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "EXTRA" {
		t.Errorf("expected EXTRA only in A, got %v", r.OnlyInA)
	}
}

func TestCompareOnlyInB(t *testing.T) {
	a := buildCompareVault(map[string]string{"FOO": "bar"})
	b := buildCompareVault(map[string]string{"FOO": "bar", "NEW": "val"})
	r := CompareVaults(a, b)
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "NEW" {
		t.Errorf("expected NEW only in B, got %v", r.OnlyInB)
	}
}

func TestCompareChangedValues(t *testing.T) {
	a := buildCompareVault(map[string]string{"FOO": "old"})
	b := buildCompareVault(map[string]string{"FOO": "new"})
	r := CompareVaults(a, b)
	if len(r.Changed) != 1 || r.Changed[0] != "FOO" {
		t.Errorf("expected FOO in changed, got %v", r.Changed)
	}
}

func TestCompareSummaryOutput(t *testing.T) {
	a := buildCompareVault(map[string]string{"A": "1", "B": "old"})
	b := buildCompareVault(map[string]string{"B": "new", "C": "3"})
	r := CompareVaults(a, b)
	summary := r.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	for _, want := range []string{"Only in A", "Only in B", "Changed"} {
		if !containsStr(summary, want) {
			t.Errorf("summary missing %q: %s", want, summary)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && (s[:len(sub)] == sub || containsStr(s[1:], sub)))
}
