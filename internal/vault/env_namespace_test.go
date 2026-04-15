package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildNamespaceVault(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	vaultFile := filepath.Join(dir, "test.vault")
	v, err := LoadOrCreate(vaultFile)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["DB_HOST"] = Entry{Value: "localhost"}
	v.Entries["DB_PORT"] = Entry{Value: "5432"}
	v.Entries["API_KEY"] = Entry{Value: "secret"}
	if err := v.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return dir, vaultFile
}

func TestLoadNamespacesNoFile(t *testing.T) {
	_, vaultFile := buildNamespaceVault(t)
	store, err := LoadNamespaces(vaultFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Namespaces) != 0 {
		t.Errorf("expected empty namespaces, got %d", len(store.Namespaces))
	}
}

func TestAssignAndGetNamespace(t *testing.T) {
	_, vaultFile := buildNamespaceVault(t)
	if err := AssignNamespace(vaultFile, "database", "DB_HOST"); err != nil {
		t.Fatalf("AssignNamespace: %v", err)
	}
	if err := AssignNamespace(vaultFile, "database", "DB_PORT"); err != nil {
		t.Fatalf("AssignNamespace: %v", err)
	}
	keys, err := GetNamespaceKeys(vaultFile, "database")
	if err != nil {
		t.Fatalf("GetNamespaceKeys: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestAssignDuplicateKeyReturnsError(t *testing.T) {
	_, vaultFile := buildNamespaceVault(t)
	_ = AssignNamespace(vaultFile, "database", "DB_HOST")
	err := AssignNamespace(vaultFile, "database", "DB_HOST")
	if err == nil {
		t.Error("expected error for duplicate key, got nil")
	}
}

func TestUnassignNamespace(t *testing.T) {
	_, vaultFile := buildNamespaceVault(t)
	_ = AssignNamespace(vaultFile, "database", "DB_HOST")
	_ = AssignNamespace(vaultFile, "database", "DB_PORT")
	if err := UnassignNamespace(vaultFile, "database", "DB_HOST"); err != nil {
		t.Fatalf("UnassignNamespace: %v", err)
	}
	keys, _ := GetNamespaceKeys(vaultFile, "database")
	if len(keys) != 1 || keys[0] != "DB_PORT" {
		t.Errorf("expected [DB_PORT], got %v", keys)
	}
}

func TestUnassignNonExistentKeyReturnsError(t *testing.T) {
	_, vaultFile := buildNamespaceVault(t)
	_ = AssignNamespace(vaultFile, "database", "DB_HOST")
	err := UnassignNamespace(vaultFile, "database", "MISSING")
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestGetNamespaceNotFoundReturnsError(t *testing.T) {
	_, vaultFile := buildNamespaceVault(t)
	_, err := GetNamespaceKeys(vaultFile, "nonexistent")
	if err == nil {
		t.Error("expected error for missing namespace, got nil")
	}
}

func TestNamespaceFilePermissions(t *testing.T) {
	_, vaultFile := buildNamespaceVault(t)
	_ = AssignNamespace(vaultFile, "api", "API_KEY")
	nsFile := namespacePath(vaultFile)
	info, err := os.Stat(nsFile)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestFormatNamespaceListEmpty(t *testing.T) {
	store := &NamespaceStore{Namespaces: make(map[string][]string)}
	out := FormatNamespaceList(store)
	if out != "no namespaces defined" {
		t.Errorf("unexpected output: %q", out)
	}
}
