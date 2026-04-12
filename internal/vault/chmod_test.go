package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildChmodVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadPermissionsNoFile(t *testing.T) {
	vp := buildChmodVault(t)
	pm, err := LoadPermissions(vp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(pm.Permissions) != 0 {
		t.Errorf("expected empty permissions, got %d entries", len(pm.Permissions))
	}
}

func TestSetAndLoadPermission(t *testing.T) {
	vp := buildChmodVault(t)
	if err := SetPermission(vp, "DB_PASS", "alice", true); err != nil {
		t.Fatalf("SetPermission: %v", err)
	}
	pm, err := LoadPermissions(vp)
	if err != nil {
		t.Fatalf("LoadPermissions: %v", err)
	}
	p, ok := pm.Permissions["DB_PASS"]
	if !ok {
		t.Fatal("expected DB_PASS permission to exist")
	}
	if p.Owner != "alice" {
		t.Errorf("expected owner alice, got %s", p.Owner)
	}
	if !p.ReadOnly {
		t.Error("expected read_only=true")
	}
	if p.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestIsReadOnly(t *testing.T) {
	vp := buildChmodVault(t)
	_ = SetPermission(vp, "API_KEY", "", true)
	_ = SetPermission(vp, "LOG_LEVEL", "", false)

	ro, err := IsReadOnly(vp, "API_KEY")
	if err != nil || !ro {
		t.Errorf("expected API_KEY to be read-only, got ro=%v err=%v", ro, err)
	}
	ro, err = IsReadOnly(vp, "LOG_LEVEL")
	if err != nil || ro {
		t.Errorf("expected LOG_LEVEL to be writable, got ro=%v err=%v", ro, err)
	}
	ro, err = IsReadOnly(vp, "UNKNOWN")
	if err != nil || ro {
		t.Errorf("expected UNKNOWN to default to writable, got ro=%v err=%v", ro, err)
	}
}

func TestRemovePermission(t *testing.T) {
	vp)
	_ = SetPermission(vp, "SECRET", "bob", true)
	if err := RemovePermission(vp, "SECRET"); err != nil {
		t.Fatalf("RemovePermission: %v", err)
	}
	pm, _ := LoadPermissions(vp)
	if _, ok := pm.Permissions["SECRET"]; ok {
		t.")
	}
}

func TestPermissionFilePermissions(t *testing.T) {
	vp := buildChmodVault(t)
	_ = SetPermission(vp, "X", "", false)
	info, err := os.Stat(permPath(vp))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
