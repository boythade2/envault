package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/user/envault/internal/vault"
)

func TestListCommandRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "list" {
			return
		}
	}
	t.Error("list command not registered on root")
}

func TestListEmptyVault(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(origDir)

	buf := &bytes.Buffer{}
	listCmd.SetOut(buf)

	if err := listCmd.RunE(listCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No entries found") {
		t.Errorf("expected empty vault message, got: %s", buf.String())
	}
}

func TestListShowsEntries(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(origDir)

	v, _ := vault.LoadOrCreate(dir)
	v.AddEntry("dev", ".env.dev.enc", true)
	v.AddEntry("prod", ".env.prod.enc", false)
	_ = v.Save(dir)

	buf := &bytes.Buffer{}
	listCmd.SetOut(buf)

	if err := listCmd.RunE(listCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"dev", ".env.dev.enc", "yes", "prod", ".env.prod.enc", "no"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestListOutputHasHeaders(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(origDir)

	v, _ := vault.LoadOrCreate(dir)
	v.AddEntry("ci", ".env.ci.enc", true)
	_ = v.Save(dir)

	buf := &bytes.Buffer{}
	listCmd.SetOut(buf)
	_ = listCmd.RunE(listCmd, []string{})

	out := buf.String()
	for _, header := range []string{"NAME", "FILE", "ENCRYPTED", "UPDATED"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output, got:\n%s", header, out)
		}
	}
}
