package vault

import (
	"strings"
	"testing"
)

func buildEnvVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["DB_HOST"] = Entry{Value: "localhost"}
	v.Entries["DB_PORT"] = Entry{Value: "5432"}
	v.Entries["APP_SECRET"] = Entry{Value: "original"}
	return v
}

func TestApplyEnvOverridesNoPrefix(t *testing.T) {
	v := buildEnvVault(t)
	t.Setenv("DB_HOST", "prod-db.example.com")
	t.Setenv("DB_PORT", "5433")

	result := ApplyEnvOverrides(v, "", false)

	if v.Entries["DB_HOST"].Value != "prod-db.example.com" {
		t.Errorf("expected DB_HOST to be overridden, got %s", v.Entries["DB_HOST"].Value)
	}
	if v.Entries["DB_PORT"].Value != "5433" {
		t.Errorf("expected DB_PORT to be overridden, got %s", v.Entries["DB_PORT"].Value)
	}
	if len(result.Applied) < 2 {
		t.Errorf("expected at least 2 applied, got %d", len(result.Applied))
	}
}

func TestApplyEnvOverridesWithPrefix(t *testing.T) {
	v := buildEnvVault(t)
	t.Setenv("ENVAULT_DB_HOST", "prefixed-host")
	t.Setenv("ENVAULT_UNKNOWN", "should-be-skipped")

	result := ApplyEnvOverrides(v, "ENVAULT_", false)

	if v.Entries["DB_HOST"].Value != "prefixed-host" {
		t.Errorf("expected DB_HOST overridden via prefix, got %s", v.Entries["DB_HOST"].Value)
	}
	for _, k := range result.NotFound {
		if k == "UNKNOWN" {
			return
		}
	}
	t.Errorf("expected UNKNOWN in NotFound list, got %v", result.NotFound)
}

func TestApplyEnvOverridesAllowNew(t *testing.T) {
	v := buildEnvVault(t)
	t.Setenv("NEW_KEY", "new-value")

	result := ApplyEnvOverrides(v, "", true)

	if _, ok := v.Entries["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY to be added when allowNew=true")
	}
	found := false
	for _, k := range result.Applied {
		if k == "NEW_KEY" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected NEW_KEY in Applied list, got %v", result.Applied)
	}
}

func TestFormatEnvOverrideResultEmpty(t *testing.T) {
	r := EnvOverrideResult{}
	out := FormatEnvOverrideResult(r)
	if !strings.Contains(out, "No environment variables applied") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatEnvOverrideResultWithEntries(t *testing.T) {
	r := EnvOverrideResult{
		Applied:  []string{"DB_HOST", "DB_PORT"},
		NotFound: []string{"MISSING_KEY"},
	}
	out := FormatEnvOverrideResult(r)
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "MISSING_KEY") {
		t.Error("expected MISSING_KEY in skipped output")
	}
	if !strings.Contains(out, "Applied") {
		t.Error("expected 'Applied' header in output")
	}
}
