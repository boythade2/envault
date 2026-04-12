package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envault/internal/vault"
)

func TestProfileCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "profile" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("profile command not registered")
	}
}

func TestProfileAddAndList(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	var buf bytes.Buffer
	profileAddCmd.SetOut(&buf)
	err := profileAddCmd.RunE(profileAddCmd, []string{"dev", "dev.vault"})
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if !strings.Contains(buf.String(), "dev") {
		t.Errorf("expected 'dev' in output, got: %s", buf.String())
	}

	var listBuf bytes.Buffer
	profileListCmd.SetOut(&listBuf)
	if err := profileListCmd.RunE(profileListCmd, []string{}); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(listBuf.String(), "dev.vault") {
		t.Errorf("expected 'dev.vault' in list output, got: %s", listBuf.String())
	}
}

func TestProfileListEmpty(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	var buf bytes.Buffer
	profileListCmd.SetOut(&buf)
	if err := profileListCmd.RunE(profileListCmd, []string{}); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(buf.String(), "No profiles") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestProfileRemove(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	store, _ := vault.LoadProfiles(dir)
	_ = store.AddProfile("staging", "staging.vault")
	_ = store.Save(dir)

	var buf bytes.Buffer
	profileRemoveCmd.SetOut(&buf)
	if err := profileRemoveCmd.RunE(profileRemoveCmd, []string{"staging"}); err != nil {
		t.Fatalf("remove: %v", err)
	}

	reloaded, _ := vault.LoadProfiles(dir)
	if _, err := reloaded.GetProfile("staging"); err == nil {
		t.Fatal("expected profile to be removed")
	}
}

func TestProfileAddDuplicateReturnsError(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	_ = profileAddCmd.RunE(profileAddCmd, []string{"prod", "prod.vault"})
	err := profileAddCmd.RunE(profileAddCmd, []string{"prod", "other.vault"})
	if err == nil {
		t.Fatal("expected error for duplicate profile")
	}
}

func TestProfileFileCreated(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(old)

	_ = profileAddCmd.RunE(profileAddCmd, []string{"dev", "dev.vault"})
	if _, err := os.Stat(filepath.Join(dir, ".envault_profiles.json")); err != nil {
		t.Fatalf("profile file not created: %v", err)
	}
}
