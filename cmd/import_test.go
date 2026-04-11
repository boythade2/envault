package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "import [file]" {
			return
		}
	}
	t.Error("import command not registered")
}

func TestImportRequiresArg(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"import"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no file argument provided")
	}
}

func TestImportDotenvFormat(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	vaultFile := filepath.Join(dir, "vault.json")

	os.WriteFile(envFile, []byte("HELLO=world\nFOO=bar\n"), 0600)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"import", envFile, "--vault", vaultFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Imported 2 variable(s)") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestImportJSONFormat(t *testing.T) {
	dir := t.TempDir()
	jsonFile := filepath.Join(dir, "vars.json")
	vaultFile := filepath.Join(dir, "vault.json")

	os.WriteFile(jsonFile, []byte(`{"KEY1":"v1","KEY2":"v2"}`), 0600)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"import", jsonFile, "--vault", vaultFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Imported 2 variable(s)") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestImportUnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	vaultFile := filepath.Join(dir, "vault.json")
	srcFile := filepath.Join(dir, "vars.toml")
	os.WriteFile(srcFile, []byte("key = \"value\"\n"), 0600)

	rootCmd.SetArgs([]string{"import", srcFile, "--vault", vaultFile, "--format", "toml"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
