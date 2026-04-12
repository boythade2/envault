package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func buildTemplateVault() *Vault {
	v := &Vault{Entries: map[string]Entry{}}
	v.Entries["APP_HOST"] = Entry{Value: "localhost", UpdatedAt: time.Now()}
	v.Entries["APP_PORT"] = Entry{Value: "8080", UpdatedAt: time.Now()}
	v.Entries["DB_NAME"] = Entry{Value: "mydb", UpdatedAt: time.Now()}
	return v
}

func writeTemplate(t *testing.T, dir, content string) string {
	t.Helper()
	p := filepath.Join(dir, "tmpl.env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRenderTemplateAllPresent(t *testing.T) {
	v := buildTemplateVault()
	dir := t.TempDir()
	p := writeTemplate(t, dir, "HOST={{APP_HOST}}\nPORT={{ APP_PORT }}\n")

	res, err := RenderTemplate(p, v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
	if res.Rendered != "HOST=localhost\nPORT=8080\n" {
		t.Errorf("unexpected rendered output: %q", res.Rendered)
	}
}

func TestRenderTemplateMissingKey(t *testing.T) {
	v := buildTemplateVault()
	dir := t.TempDir()
	p := writeTemplate(t, dir, "HOST={{APP_HOST}}\nSECRET={{MISSING_KEY}}\n")

	res, err := RenderTemplate(p, v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING_KEY" {
		t.Errorf("expected [MISSING_KEY], got %v", res.Missing)
	}
	if res.Rendered != "HOST=localhost\nSECRET={{MISSING_KEY}}\n" {
		t.Errorf("unexpected rendered output: %q", res.Rendered)
	}
}

func TestRenderTemplateDuplicateMissingKey(t *testing.T) {
	v := buildTemplateVault()
	dir := t.TempDir()
	p := writeTemplate(t, dir, "A={{GHOST}}\nB={{GHOST}}\n")

	res, err := RenderTemplate(p, v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 unique missing key, got %v", res.Missing)
	}
}

func TestRenderTemplateFileNotFound(t *testing.T) {
	v := buildTemplateVault()
	_, err := RenderTemplate("/nonexistent/path/tmpl.env", v)
	if err == nil {
		t.Fatal("expected error for missing template file")
	}
}
