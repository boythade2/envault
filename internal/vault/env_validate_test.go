package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildValidateVault(t *testing.T, entries map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vault.json")
	v, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	for k, val := range entries {
		v.Entries = append(v.Entries, Entry{Key: k, Value: val})
	}
	if err := v.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

func TestValidateEnvKeysClean(t *testing.T) {
	path := buildValidateVault(t, map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"APP_PORT":     "8080",
	})
	results, err := ValidateEnvKeys(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no issues, got %d", len(results))
	}
}

func TestValidateEnvKeysLowercaseWarn(t *testing.T) {
	path := buildValidateVault(t, map[string]string{"db_host": "localhost"})
	results, err := ValidateEnvKeys(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Warning == "" {
		t.Error("expected a warning for lowercase key")
	}
	if results[0].Error != "" {
		t.Errorf("unexpected error: %s", results[0].Error)
	}
}

func TestValidateEnvKeysSpaceError(t *testing.T) {
	path := buildValidateVault(t, map[string]string{"BAD KEY": "value"})
	results, err := ValidateEnvKeys(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Error == "" {
		t.Error("expected an error for key with space")
	}
}

func TestValidateEnvKeysMissingFile(t *testing.T) {
	_, err := ValidateEnvKeys(filepath.Join(t.TempDir(), "missing.json"))
	// LoadOrCreate creates a new file so no error expected
	if err != nil {
		t.Errorf("unexpected error for missing file: %v", err)
	}
}

func TestFormatEnvValidationResultsEmpty(t *testing.T) {
	out := FormatEnvValidationResults(nil)
	if out != "all keys are valid" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatEnvValidationResultsOutput(t *testing.T) {
	results := []EnvValidationResult{
		{Key: "bad key", Error: "key contains whitespace"},
		{Key: "lower", Warning: "key contains lowercase letters; consider uppercasing"},
	}
	out := FormatEnvValidationResults(results)
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !containsStr(out, "[ERROR]") {
		t.Error("expected [ERROR] tag in output")
	}
	if !containsStr(out, "[WARN]") {
		t.Error("expected [WARN] tag in output")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	return strings.Contains(s, sub)
}

func init() {
	_ = os.Getenv // ensure os imported
}
