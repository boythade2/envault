package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envault/internal/vault"
)

func writeGroupVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestGroupCommandRegistered(t *testing.T) {
	for _, sub := range []string{"add", "remove", "assign", "unassign", "list"} {
		found := false
		for _, c := range rootCmd.Commands() {
			if c.Use == "group" {
				for _, sc := range c.Commands() {
					if strings.HasPrefix(sc.Use, sub) {
						found = true
					}
				}
			}
		}
		if !found {
			t.Errorf("subcommand %q not registered under group", sub)
		}
	}
}

func TestGroupAddAndList(t *testing.T) {
	vp := writeGroupVault(t)
	out := &bytes.Buffer{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"group", "add", vp, "backend"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("group add: %v", err)
	}
	out.Reset()
	rootCmd.SetArgs([]string{"group", "list", vp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("group list: %v", err)
	}
	if !strings.Contains(out.String(), "backend") {
		t.Errorf("expected 'backend' in output, got: %s", out.String())
	}
}

func TestGroupListEmpty(t *testing.T) {
	vp := writeGroupVault(t)
	out := &bytes.Buffer{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"group", "list", vp})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("group list: %v", err)
	}
	if !strings.Contains(out.String(), "No groups") {
		t.Errorf("expected empty message, got: %s", out.String())
	}
}

func TestGroupAssignAndUnassign(t *testing.T) {
	vp := writeGroupVault(t)
	_ = vault.AddGroup(vp, "db")
	out := &bytes.Buffer{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"group", "assign", vp, "db", "DB_HOST"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("group assign: %v", err)
	}
	gs, _ := vault.LoadGroups(vp)
	if len(gs.Groups["db"].Keys) != 1 {
		t.Errorf("expected 1 key after assign")
	}
	rootCmd.SetArgs([]string{"group", "unassign", vp, "db", "DB_HOST"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("group unassign: %v", err)
	}
	gs, _ = vault.LoadGroups(vp)
	if len(gs.Groups["db"].Keys) != 0 {
		t.Errorf("expected 0 keys after unassign")
	}
}

func TestGroupRemove(t *testing.T) {
	vp := writeGroupVault(t)
	_ = vault.AddGroup(vp, "ops")
	out := &bytes.Buffer{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"group", "remove", vp, "ops"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("group remove: %v", err)
	}
	gs, _ := vault.LoadGroups(vp)
	if _, ok := gs.Groups["ops"]; ok {
		t.Error("expected group 'ops' to be removed")
	}
	_ = json.Marshal
	_ = os.Stdout
}
