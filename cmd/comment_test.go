package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envault/internal/vault"
)

func writeCommentVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vf := filepath.Join(dir, "test.vault")
	v, err := vault.LoadOrCreate(vf)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["API_KEY"] = vault.Entry{Value: "abc"}
	v.Entries["DB_PASS"] = vault.Entry{Value: "xyz"}
	if err := v.Save(vf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return vf
}

func TestCommentCommandRegistered(t *testing.T) {
	for _, sub := range []string{"set", "unset", "list"} {
		found := false
		for _, c := range rootCmd.Commands() {
			if c.Use == "comment" {
				for _, sc := range c.Commands() {
					if strings.HasPrefix(sc.Use, sub) {
						found = true
					}
				}
			}
		}
		if !found {
			t.Errorf("subcommand %q not registered", sub)
		}
	}
}

func TestCommentSetAndList(t *testing.T) {
	vf := writeCommentVault(t)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"comment", "set", vf, "API_KEY", "used for auth"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("set: %v", err)
	}
	buf.Reset()
	rootCmd.SetArgs([]string{"comment", "list", vf})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	cs, _ := vault.LoadComments(vf)
	if vault.GetComment(cs, "API_KEY") != "used for auth" {
		t.Error("expected comment to be stored")
	}
}

func TestCommentUnset(t *testing.T) {
	vf := writeCommentVault(t)
	_ = vault.SetComment(vf, "DB_PASS", "db password")
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetArgs([]string{"comment", "unset", vf, "DB_PASS"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unset: %v", err)
	}
	cs, _ := vault.LoadComments(vf)
	if vault.GetComment(cs, "DB_PASS") != "" {
		t.Error("expected comment to be removed")
	}
}

func TestCommentListEmpty(t *testing.T) {
	vf := writeCommentVault(t)
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"comment", "list", vf})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(buf.String(), "No comments found") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
