package vault

import (
	"testing"
	"time"
)

func buildSortVault(t *testing.T) *Vault {
	t.Helper()
	now := time.Now()
	v := &Vault{Entries: map[string]Entry{}}
	v.Entries["ZEBRA"] = Entry{Value: "alpha", UpdatedAt: now.Add(-2 * time.Minute)}
	v.Entries["APPLE"] = Entry{Value: "mango", UpdatedAt: now.Add(-1 * time.Minute)}
	v.Entries["MANGO"] = Entry{Value: "cherry", UpdatedAt: now}
	return v
}

func TestSortByKeyAsc(t *testing.T) {
	v := buildSortVault(t)
	results := SortEntries(v, SortOptions{Field: SortByKey, Order: SortAsc})
	if results[0].Key != "APPLE" || results[1].Key != "MANGO" || results[2].Key != "ZEBRA" {
		t.Fatalf("unexpected key order: %v", results)
	}
}

func TestSortByKeyDesc(t *testing.T) {
	v := buildSortVault(t)
	results := SortEntries(v, SortOptions{Field: SortByKey, Order: SortDesc})
	if results[0].Key != "ZEBRA" || results[1].Key != "MANGO" || results[2].Key != "APPLE" {
		t.Fatalf("unexpected key order: %v", results)
	}
}

func TestSortByValueAsc(t *testing.T) {
	v := buildSortVault(t)
	results := SortEntries(v, SortOptions{Field: SortByValue, Order: SortAsc})
	// alpha < cherry < mango
	if results[0].Value.Value != "alpha" || results[1].Value.Value != "cherry" || results[2].Value.Value != "mango" {
		t.Fatalf("unexpected value order: %v", results)
	}
}

func TestSortByUpdatedAtAsc(t *testing.T) {
	v := buildSortVault(t)
	results := SortEntries(v, SortOptions{Field: SortByUpdatedAt, Order: SortAsc})
	// ZEBRA is oldest
	if results[0].Key != "ZEBRA" {
		t.Fatalf("expected ZEBRA first (oldest), got %s", results[0].Key)
	}
	if results[2].Key != "MANGO" {
		t.Fatalf("expected MANGO last (newest), got %s", results[2].Key)
	}
}

func TestSortByUpdatedAtDesc(t *testing.T) {
	v := buildSortVault(t)
	results := SortEntries(v, SortOptions{Field: SortByUpdatedAt, Order: SortDesc})
	if results[0].Key != "MANGO" {
		t.Fatalf("expected MANGO first (newest), got %s", results[0].Key)
	}
}

func TestSortEmptyVault(t *testing.T) {
	v := &Vault{Entries: map[string]Entry{}}
	results := SortEntries(v, SortOptions{Field: SortByKey, Order: SortAsc})
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}
