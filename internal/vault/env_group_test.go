package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildGroupVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadGroupsNoFile(t *testing.T) {
	vp := buildGroupVault(t)
	gs, err := LoadGroups(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gs.Groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(gs.Groups))
	}
}

func TestAddAndRemoveGroup(t *testing.T) {
	vp := buildGroupVault(t)
	if err := AddGroup(vp, "backend"); err != nil {
		t.Fatalf("AddGroup: %v", err)
	}
	gs, _ := LoadGroups(vp)
	if _, ok := gs.Groups["backend"]; !ok {
		t.Error("expected group 'backend' to exist")
	}
	if err := RemoveGroup(vp, "backend"); err != nil {
		t.Fatalf("RemoveGroup: %v", err)
	}
	gs, _ = LoadGroups(vp)
	if _, ok := gs.Groups["backend"]; ok {
		t.Error("expected group 'backend' to be removed")
	}
}

func TestAddDuplicateGroupReturnsError(t *testing.T) {
	vp := buildGroupVault(t)
	_ = AddGroup(vp, "frontend")
	if err := AddGroup(vp, "frontend"); err == nil {
		t.Error("expected error for duplicate group")
	}
}

func TestRemoveNonExistentGroupReturnsError(t *testing.T) {
	vp := buildGroupVault(t)
	if err := RemoveGroup(vp, "ghost"); err == nil {
		t.Error("expected error for non-existent group")
	}
}

func TestAssignAndUnassignKey(t *testing.T) {
	vp := buildGroupVault(t)
	_ = AddGroup(vp, "db")
	if err := AssignKeyToGroup(vp, "db", "DB_HOST"); err != nil {
		t.Fatalf("AssignKeyToGroup: %v", err)
	}
	if err := AssignKeyToGroup(vp, "db", "DB_PORT"); err != nil {
		t.Fatalf("AssignKeyToGroup: %v", err)
	}
	gs, _ := LoadGroups(vp)
	if len(gs.Groups["db"].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(gs.Groups["db"].Keys))
	}
	if err := UnassignKeyFromGroup(vp, "db", "DB_HOST"); err != nil {
		t.Fatalf("UnassignKeyFromGroup: %v", err)
	}
	gs, _ = LoadGroups(vp)
	if len(gs.Groups["db"].Keys) != 1 {
		t.Errorf("expected 1 key after unassign, got %d", len(gs.Groups["db"].Keys))
	}
}

func TestAssignDuplicateKeyReturnsError(t *testing.T) {
	vp := buildGroupVault(t)
	_ = AddGroup(vp, "api")
	_ = AssignKeyToGroup(vp, "api", "API_KEY")
	if err := AssignKeyToGroup(vp, "api", "API_KEY"); err == nil {
		t.Error("expected error for duplicate key assignment")
	}
}

func TestGroupFilePermissions(t *testing.T) {
	vp := buildGroupVault(t)
	_ = AddGroup(vp, "ops")
	p := groupPath(vp)
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
