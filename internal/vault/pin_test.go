package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func buildPinVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.vault")
}

func TestLoadPinsNoFile(t *testing.T) {
	vp := buildPinVault(t)
	pl, err := LoadPins(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pl.Pins) != 0 {
		t.Errorf("expected empty pin list, got %d pins", len(pl.Pins))
	}
}

func TestPinAndUnpin(t *testing.T) {
	vp := buildPinVault(t)
	if err := PinKey(vp, "DB_PASSWORD", "do not change"); err != nil {
		t.Fatalf("PinKey: %v", err)
	}
	pinned, err := IsPinned(vp, "DB_PASSWORD")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if !pinned {
		t.Error("expected DB_PASSWORD to be pinned")
	}
	if err := UnpinKey(vp, "DB_PASSWORD"); err != nil {
		t.Fatalf("UnpinKey: %v", err)
	}
	pinned, err = IsPinned(vp, "DB_PASSWORD")
	if err != nil {
		t.Fatalf("IsPinned after unpin: %v", err)
	}
	if pinned {
		t.Error("expected DB_PASSWORD to be unpinned")
	}
}

func TestPinDuplicateKey(t *testing.T) {
	vp := buildPinVault(t)
	if err := PinKey(vp, "SECRET", ""); err != nil {
		t.Fatalf("first PinKey: %v", err)
	}
	if err := PinKey(vp, "SECRET", ""); err == nil {
		t.Error("expected error pinning duplicate key, got nil")
	}
}

func TestUnpinNonExistentKey(t *testing.T) {
	vp := buildPinVault(t)
	if err := UnpinKey(vp, "MISSING"); err == nil {
		t.Error("expected error unpinning non-existent key, got nil")
	}
}

func TestPinFilePermissions(t *testing.T) {
	vp := buildPinVault(t)
	if err := PinKey(vp, "API_KEY", "critical"); err != nil {
		t.Fatalf("PinKey: %v", err)
	}
	pp := pinPath(vp)
	info, err := os.Stat(pp)
	if err != nil {
		t.Fatalf("stat pin file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %o", info.Mode().Perm())
	}
}

func TestPinPreservesNote(t *testing.T) {
	vp := buildPinVault(t)
	const note = "managed by infra team"
	if err := PinKey(vp, "INFRA_TOKEN", note); err != nil {
		t.Fatalf("PinKey: %v", err)
	}
	pl, err := LoadPins(vp)
	if err != nil {
		t.Fatalf("LoadPins: %v", err)
	}
	if len(pl.Pins) != 1 || pl.Pins[0].Note != note {
		t.Errorf("expected note %q, got %q", note, pl.Pins[0].Note)
	}
}
