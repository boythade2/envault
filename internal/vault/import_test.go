package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportDotenv(t *testing.T) {
	content := `# comment
DB_HOST=localhost
DB_PORT=5432
SECRET="mysecret"
`
	tmp := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(tmp, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := ImportDotenv(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"SECRET":  "mysecret",
	}
	for k, v := range expect {
		if got[k] != v {
			t.Errorf("key %q: want %q, got %q", k, v, got[k])
		}
	}
}

func TestImportDotenvSkipsComments(t *testing.T) {
	content := "# this is a comment\nKEY=value\n"
	tmp := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(tmp, []byte(content), 0600)

	got, _ := ImportDotenv(tmp)
	if _, ok := got["# this is a comment"]; ok {
		t.Error("comment line should not be imported")
	}
	if got["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", got["KEY"])
	}
}

func TestImportDotenvMissingFile(t *testing.T) {
	_, err := ImportDotenv("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestImportJSON(t *testing.T) {
	content := `{"APP_ENV":"production","PORT":"8080"}`
	tmp := filepath.Join(t.TempDir(), "vars.json")
	if err := os.WriteFile(tmp, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := ImportJSON(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: want production, got %q", got["APP_ENV"])
	}
	if got["PORT"] != "8080" {
		t.Errorf("PORT: want 8080, got %q", got["PORT"])
	}
}

func TestImportJSONInvalidFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not json"), 0600)
	_, err := ImportJSON(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestImportJSONMissingFile(t *testing.T) {
	_, err := ImportJSON("/nonexistent/path/vars.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
