package vault

import (
	"testing"
)

func seedVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{Entries: make(map[string]Entry)}
	v.Entries["DATABASE_URL"] = Entry{Value: "postgres://localhost/mydb"}
	v.Entries["DATABASE_PASS"] = Entry{Value: "s3cr3t"}
	v.Entries["API_KEY"] = Entry{Value: "abc123"}
	v.Entries["APP_ENV"] = Entry{Value: "production"}
	return v
}

func TestSearchByKeyExact(t *testing.T) {
	v := seedVault(t)
	results := v.Search("API_KEY", false)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", results[0].Key)
	}
}

func TestSearchByKeyPartial(t *testing.T) {
	v := seedVault(t)
	results := v.Search("database", false)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// Results should be sorted: DATABASE_PASS before DATABASE_URL
	if results[0].Key != "DATABASE_PASS" {
		t.Errorf("expected DATABASE_PASS first, got %s", results[0].Key)
	}
}

func TestSearchByValueEnabled(t *testing.T) {
	v := seedVault(t)
	results := v.Search("abc123", true)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", results[0].Key)
	}
}

func TestSearchByValueDisabled(t *testing.T) {
	v := seedVault(t)
	// searching for a value substring with searchValues=false should return nothing
	results := v.Search("abc123", false)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchNoMatch(t *testing.T) {
	v := seedVault(t)
	results := v.Search("NONEXISTENT", true)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	v := seedVault(t)
	results := v.Search("app_env", false)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", results[0].Key)
	}
}
