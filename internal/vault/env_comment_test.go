package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildCommentVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	vf := filepath.Join(dir, "test.vault")
	v, err := LoadOrCreate(vf)
	if err != nil {
		t.Fatalf("LoadOrCreate: %v", err)
	}
	v.Entries["API_KEY"] = Entry{Value: "abc123"}
	v.Entries["DB_PASS"] = Entry{Value: "secret"}
	if err := v.Save(vf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return vf
}

func TestLoadCommentsNoFile(t *testing.T) {
	dir := t.TempDir()
	vf := filepath.Join(dir, "test.vault")
	cs, err := LoadComments(vf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs.Comments) != 0 {
		t.Errorf("expected empty comments, got %d", len(cs.Comments))
	}
}

func TestSetAndGetComment(t *testing.T) {
	vf := buildCommentVault(t)
	if err := SetComment(vf, "API_KEY", "used for auth"); err != nil {
		t.Fatalf("SetComment: %v", err)
	}
	cs, err := LoadComments(vf)
	if err != nil {
		t.Fatalf("LoadComments: %v", err)
	}
	got := GetComment(cs, "API_KEY")
	if got != "used for auth" {
		t.Errorf("expected 'used for auth', got %q", got)
	}
}

func TestRemoveComment(t *testing.T) {
	vf := buildCommentVault(t)
	_ = SetComment(vf, "DB_PASS", "database password")
	if err := RemoveComment(vf, "DB_PASS"); err != nil {
		t.Fatalf("RemoveComment: %v", err)
	}
	cs, _ := LoadComments(vf)
	if GetComment(cs, "DB_PASS") != "" {
		t.Error("expected comment to be removed")
	}
}

func TestCommentFilePermissions(t *testing.T) {
	vf := buildCommentVault(t)
	_ = SetComment(vf, "API_KEY", "test")
	cp := commentPath(vf)
	info, err := os.Stat(cp)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestGetCommentMissingKeyReturnsEmpty(t *testing.T) {
	vf := buildCommentVault(t)
	cs, _ := LoadComments(vf)
	if got := GetComment(cs, "NONEXISTENT"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
