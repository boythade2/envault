package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildSchemaVault(t *testing.T) (string, *Vault) {
	t.Helper()
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	v, err := LoadOrCreate(vaultPath)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries = append(v.Entries,
		Entry{Key: "APP_HOST", Value: "localhost"},
		Entry{Key: "APP_PORT", Value: "8080"},
	)
	return vaultPath, v
}

func TestLoadSchemaNoFile(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	s, err := LoadSchema(vaultPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(s.Rules) != 0 {
		t.Errorf("expected empty schema, got %d rules", len(s.Rules))
	}
}

func TestSaveAndLoadSchema(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	s := Schema{
		Rules: []SchemaRule{
			{Key: "APP_HOST", Required: true},
			{Key: "APP_PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	if err := SaveSchema(vaultPath, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	loaded, err := LoadSchema(vaultPath)
	if err != nil {
		t.Fatalf("LoadSchema: %v", err)
	}
	if len(loaded.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(loaded.Rules))
	}
}

func TestValidateSchemaAllPresent(t *testing.T) {
	vaultPath, v := buildSchemaVault(t)
	_ = vaultPath
	s := Schema{
		Rules: []SchemaRule{
			{Key: "APP_HOST", Required: true},
			{Key: "APP_PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	violations := ValidateSchema(v, s)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidateSchemaMissingRequired(t *testing.T) {
	_, v := buildSchemaVault(t)
	s := Schema{
		Rules: []SchemaRule{
			{Key: "DB_URL", Required: true},
		},
	}
	violations := ValidateSchema(v, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_URL" {
		t.Errorf("expected violation for DB_URL, got %s", violations[0].Key)
	}
}

func TestValidateSchemaPatternMismatch(t *testing.T) {
	_, v := buildSchemaVault(t)
	v.Entries = append(v.Entries, Entry{Key: "TOKEN", Value: "not-a-number"})
	s := Schema{
		Rules: []SchemaRule{
			{Key: "TOKEN", Pattern: `^\d+$`},
		},
	}
	violations := ValidateSchema(v, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestSchemaFilePermissions(t *testing.T) {
	dir := t.TempDir()
	vaultPath := filepath.Join(dir, "test.vault")
	s := Schema{Rules: []SchemaRule{{Key: "X", Required: true}}}
	if err := SaveSchema(vaultPath, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, "test.schema.json"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
