package vault

import (
	"os"
	"strings"
	"testing"
)

func buildInjectVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{Entries: map[string]Entry{}}
	v.Entries["APP_HOST"] = Entry{Value: "localhost"}
	v.Entries["APP_PORT"] = Entry{Value: "8080"}
	v.Entries["DB_URL"] = Entry{Value: "postgres://localhost/db"}
	return v
}

func TestInjectEnvNoPrefix(t *testing.T) {
	v := buildInjectVault(t)
	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DB_URL")

	result, err := InjectEnv(v, InjectOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 3 {
		t.Errorf("expected 3 injected, got %d", len(result.Injected))
	}
	if got := os.Getenv("APP_HOST"); got != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %s", got)
	}
}

func TestInjectEnvWithPrefix(t *testing.T) {
	v := buildInjectVault(t)
	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DB_URL")

	result, err := InjectEnv(v, InjectOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(result.Injected))
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if os.Getenv("DB_URL") != "" {
		t.Error("DB_URL should not have been injected")
	}
}

func TestInjectEnvSkipsExistingWithoutOverwrite(t *testing.T) {
	v := buildInjectVault(t)
	os.Setenv("APP_HOST", "original")
	defer os.Unsetenv("APP_HOST")

	result, err := InjectEnv(v, InjectOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if os.Getenv("APP_HOST") != "original" {
		t.Error("APP_HOST should not have been overwritten")
	}
	if len(result.Skipped) == 0 {
		t.Error("expected APP_HOST in skipped")
	}
}

func TestInjectEnvOverwritesWhenEnabled(t *testing.T) {
	v := buildInjectVault(t)
	os.Setenv("APP_HOST", "original")
	defer os.Unsetenv("APP_HOST")

	result, err := InjectEnv(v, InjectOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if os.Getenv("APP_HOST") != "localhost" {
		t.Errorf("expected APP_HOST=localhost after overwrite, got %s", os.Getenv("APP_HOST"))
	}
	if len(result.Overridden) == 0 {
		t.Error("expected APP_HOST in overridden")
	}
}

func TestFormatInjectResult(t *testing.T) {
	r := InjectResult{
		Injected:   []string{"APP_PORT"},
		Overridden: []string{"APP_HOST"},
		Skipped:    []string{"DB_URL"},
	}
	out := FormatInjectResult(r)
	if !strings.Contains(out, "Injected: 1") {
		t.Error("expected injected count in output")
	}
	if !strings.Contains(out, "APP_HOST (overridden)") {
		t.Error("expected overridden key in output")
	}
}
