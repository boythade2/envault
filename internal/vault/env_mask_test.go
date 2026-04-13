package vault

import (
	"strings"
	"testing"
)

func buildMaskVault() []Entry {
	return []Entry{
		{Key: "API_KEY", Value: "supersecret"},
		{Key: "DB_PASSWORD", Value: "hunter2"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestMaskAllByDefault(t *testing.T) {
	entries := buildMaskVault()
	masked, results := MaskValues(entries, MaskOptions{})
	for i, r := range results {
		if !r.Masked {
			t.Errorf("expected entry %s to be masked", r.Key)
		}
		if masked[i].Value == entries[i].Value {
			t.Errorf("expected value to differ from original for key %s", entries[i].Key)
		}
		if !strings.Contains(masked[i].Value, "*") {
			t.Errorf("expected masked value to contain '*' for key %s", entries[i].Key)
		}
	}
}

func TestMaskByPattern(t *testing.T) {
	entries := buildMaskVault()
	opts := MaskOptions{Patterns: []string{"(?i)(key|password)"}}
	masked, results := MaskValues(entries, opts)

	maskedKeys := map[string]bool{}
	for _, r := range results {
		if r.Masked {
			maskedKeys[r.Key] = true
		}
	}
	if !maskedKeys["API_KEY"] {
		t.Error("expected API_KEY to be masked")
	}
	if !maskedKeys["DB_PASSWORD"] {
		t.Error("expected DB_PASSWORD to be masked")
	}
	if maskedKeys["APP_ENV"] {
		t.Error("expected APP_ENV NOT to be masked")
	}
	// Unmasked entries should retain original value
	for i, e := range entries {
		if !maskedKeys[e.Key] && masked[i].Value != e.Value {
			t.Errorf("expected unmasked entry %s to keep original value", e.Key)
		}
	}
}

func TestMaskShowChars(t *testing.T) {
	entries := []Entry{{Key: "SECRET", Value: "abcdefgh"}}
	masked, _ := MaskValues(entries, MaskOptions{ShowChars: 3})
	if !strings.HasPrefix(masked[0].Value, "abc") {
		t.Errorf("expected value to start with 'abc', got %s", masked[0].Value)
	}
	if !strings.Contains(masked[0].Value[3:], "*") {
		t.Error("expected trailing stars after revealed chars")
	}
}

func TestMaskCustomChar(t *testing.T) {
	entries := []Entry{{Key: "TOKEN", Value: "hello"}}
	masked, _ := MaskValues(entries, MaskOptions{MaskChar: "#"})
	if !strings.Contains(masked[0].Value, "#") {
		t.Errorf("expected '#' mask char, got %s", masked[0].Value)
	}
}

func TestMaskDoesNotMutateOriginal(t *testing.T) {
	entries := buildMaskVault()
	orig := make([]string, len(entries))
	for i, e := range entries {
		orig[i] = e.Value
	}
	MaskValues(entries, MaskOptions{})
	for i, e := range entries {
		if e.Value != orig[i] {
			t.Errorf("original entry %s was mutated", e.Key)
		}
	}
}
