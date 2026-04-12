package vault

import (
	"testing"
	"time"
)

func buildMergeVault(entries map[string]string) *Vault {
	v := &Vault{Entries: make(map[string]Entry)}
	for k, val := range entries {
		v.Entries[k] = Entry{Value: val, UpdatedAt: time.Now()}
	}
	return v
}

func TestMergeAddsNewKeys(t *testing.T) {
	dst := buildMergeVault(map[string]string{"A": "1"})
	src := buildMergeVault(map[string]string{"B": "2"})

	result, err := MergeVaults(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Added) != 1 || result.Added[0] != "B" {
		t.Errorf("expected B to be added, got %v", result.Added)
	}
	if dst.Entries["B"].Value != "2" {
		t.Errorf("expected B=2, got %s", dst.Entries["B"].Value)
	}
}

func TestMergeStrategyOursSkipsConflict(t *testing.T) {
	dst := buildMergeVault(map[string]string{"A": "original"})
	src := buildMergeVault(map[string]string{"A": "incoming"})

	result, err := MergeVaults(dst, src, MergeStrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %v", result.Skipped)
	}
	if dst.Entries["A"].Value != "original" {
		t.Errorf("expected original value preserved, got %s", dst.Entries["A"].Value)
	}
}

func TestMergeStrategyTheirsOverwrites(t *testing.T) {
	dst := buildMergeVault(map[string]string{"A": "original"})
	src := buildMergeVault(map[string]string{"A": "incoming"})

	result, err := MergeVaults(dst, src, MergeStrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Updated) != 1 {
		t.Errorf("expected 1 updated, got %v", result.Updated)
	}
	if dst.Entries["A"].Value != "incoming" {
		t.Errorf("expected incoming value, got %s", dst.Entries["A"].Value)
	}
}

func TestMergeStrategyErrorOnConflict(t *testing.T) {
	dst := buildMergeVault(map[string]string{"A": "1"})
	src := buildMergeVault(map[string]string{"A": "2"})

	result, err := MergeVaults(dst, src, MergeStrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
	if len(result.Conflict) != 1 || result.Conflict[0] != "A" {
		t.Errorf("expected conflict on A, got %v", result.Conflict)
	}
}

func TestMergeUnknownStrategyErrors(t *testing.T) {
	dst := buildMergeVault(map[string]string{"A": "1"})
	src := buildMergeVault(map[string]string{"A": "2"})

	_, err := MergeVaults(dst, src, MergeStrategy("unknown"))
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}
