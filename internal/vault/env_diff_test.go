package vault

import (
	"testing"
)

func buildEnvDiffVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["APP_HOST"] = Entry{Value: "localhost"}
	v.Entries["APP_PORT"] = Entry{Value: "8080"}
	v.Entries["APP_SECRET"] = Entry{Value: "s3cr3t"}
	return v
}

func TestEnvDiffAllMatch(t *testing.T) {
	v := buildEnvDiffVault(t)
	t.Setenv("APP_HOST", "localhost")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_SECRET", "s3cr3t")

	r := DiffEnv(v, "")
	if len(r.Match) != 3 {
		t.Errorf("expected 3 matches, got %d", len(r.Match))
	}
	if len(r.Mismatch) != 0 || len(r.OnlyInVault) != 0 {
		t.Errorf("expected no mismatches or vault-only keys")
	}
}

func TestEnvDiffMismatch(t *testing.T) {
	v := buildEnvDiffVault(t)
	t.Setenv("APP_HOST", "localhost")
	t.Setenv("APP_PORT", "9999") // different
	t.Setenv("APP_SECRET", "s3cr3t")

	r := DiffEnv(v, "")
	if len(r.Mismatch) != 1 || r.Mismatch[0] != "APP_PORT" {
		t.Errorf("expected APP_PORT in mismatch, got %v", r.Mismatch)
	}
}

func TestEnvDiffOnlyInVault(t *testing.T) {
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["VAULT_ONLY_KEY"] = Entry{Value: "val"}

	r := DiffEnv(v, "")
	found := false
	for _, k := range r.OnlyInVault {
		if k == "VAULT_ONLY_KEY" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected VAULT_ONLY_KEY in OnlyInVault, got %v", r.OnlyInVault)
	}
}

func TestEnvDiffWithPrefix(t *testing.T) {
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["HOST"] = Entry{Value: "localhost"}

	t.Setenv("APP_HOST", "localhost")
	t.Setenv("OTHER_HOST", "remotehost")

	r := DiffEnv(v, "APP_")
	if len(r.Match) != 1 || r.Match[0] != "HOST" {
		t.Errorf("expected HOST to match via prefix, got match=%v", r.Match)
	}
	for _, k := range r.OnlyInEnv {
		if k == "HOST" {
			t.Errorf("OTHER_HOST should not appear after prefix strip as HOST conflict")
		}
	}
}

func TestFormatEnvDiffResultEmpty(t *testing.T) {
	r := EnvDiffResult{}
	out := FormatEnvDiffResult(r)
	if out != "vault and environment are identical\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatEnvDiffResultShowsSections(t *testing.T) {
	r := EnvDiffResult{
		Mismatch:    []string{"KEY_A"},
		OnlyInVault: []string{"KEY_B"},
		OnlyInEnv:   []string{"KEY_C"},
		Match:       []string{"KEY_D"},
	}
	out := FormatEnvDiffResult(r)
	for _, want := range []string{"KEY_A", "KEY_B", "KEY_C", "1 key(s) in sync"} {
		if !contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
