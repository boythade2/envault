package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportDotenv(t *testing.T) {
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}

	tmpFile := filepath.Join(t.TempDir(), ".env")
	if err := v.Export(tmpFile, FormatDotenv); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("Expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("Expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestExportJSON(t *testing.T) {
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "API_KEY", Value: "secret"},
		{Key: "REGION", Value: "us-east-1"},
	}

	tmpFile := filepath.Join(t.TempDir(), "env.json")
	if err := v.Export(tmpFile, FormatJSON); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}
	if m["API_KEY"] != "secret" {
		t.Errorf("Expected API_KEY=secret, got %s", m["API_KEY"])
	}
	if m["REGION"] != "us-east-1" {
		t.Errorf("Expected REGION=us-east-1, got %s", m["REGION"])
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	v := &Vault{}
	tmpFile := filepath.Join(t.TempDir(), "out.txt")
	err := v.Export(tmpFile, ExportFormat("xml"))
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}

func TestExportFilePermissions(t *testing.T) {
	v := &Vault{}
	v.Entries = []Entry{{Key: "FOO", Value: "bar"}}

	tmpFile := filepath.Join(t.TempDir(), ".env")
	if err := v.Export(tmpFile, FormatDotenv); err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file mode 0600, got %v", info.Mode().Perm())
	}
}
