package vault

import (
	"testing"
)

func seedTagVault(t *testing.T) *Vault {
	t.Helper()
	v := &Vault{}
	v.Entries = []Entry{
		{Key: "DB_HOST", Value: "localhost", Tags: []string{"database", "prod"}},
		{Key: "DB_PASS", Value: "secret", Tags: []string{"database", "secret"}},
		{Key: "API_KEY", Value: "abc123", Tags: []string{"api", "prod"}},
		{Key: "LOG_LEVEL", Value: "info", Tags: []string{}},
	}
	return v
}

func TestTagFilterSingleTag(t *testing.T) {
	v := seedTagVault(t)
	results := v.TagFilter([]string{"database"})
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
}

func TestTagFilterMultipleTags(t *testing.T) {
	v := seedTagVault(t)
	results := v.TagFilter([]string{"database", "prod"})
	if len(results) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", results[0].Key)
	}
}

func TestTagFilterNoMatch(t *testing.T) {
	v := seedTagVault(t)
	results := v.TagFilter([]string{"nonexistent"})
	if len(results) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(results))
	}
}

func TestTagFilterEmptyTagsReturnsAll(t *testing.T) {
	v := seedTagVault(t)
	results := v.TagFilter([]string{})
	if len(results) != len(v.Entries) {
		t.Fatalf("expected all entries, got %d", len(results))
	}
}

func TestAllTagsSorted(t *testing.T) {
	v := seedTagVault(t)
	tags := v.AllTags()
	expected := []string{"api", "database", "prod", "secret"}
	if len(tags) != len(expected) {
		t.Fatalf("expected %d tags, got %d", len(expected), len(tags))
	}
	for i, tag := range tags {
		if tag != expected[i] {
			t.Errorf("expected tag %q at index %d, got %q", expected[i], i, tag)
		}
	}
}

func TestAddTag(t *testing.T) {
	v := seedTagVault(t)
	ok := v.AddTag("LOG_LEVEL", "ops")
	if !ok {
		t.Fatal("expected AddTag to return true")
	}
	results := v.TagFilter([]string{"ops"})
	if len(results) != 1 || results[0].Key != "LOG_LEVEL" {
		t.Error("expected LOG_LEVEL in ops tag results")
	}
}

func TestAddTagDuplicate(t *testing.T) {
	v := seedTagVault(t)
	v.AddTag("DB_HOST", "prod")
	var count int
	for _, t2 := range v.Entries[0].Tags {
		if t2 == "prod" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected tag to appear once, got %d", count)
	}
}

func TestRemoveTag(t *testing.T) {
	v := seedTagVault(t)
	ok := v.RemoveTag("DB_HOST", "prod")
	if !ok {
		t.Fatal("expected RemoveTag to return true")
	}
	results := v.TagFilter([]string{"prod"})
	for _, e := range results {
		if e.Key == "DB_HOST" {
			t.Error("DB_HOST should not appear in prod results after removal")
		}
	}
}

func TestAddTagMissingKey(t *testing.T) {
	v := seedTagVault(t)
	ok := v.AddTag("NONEXISTENT", "tag")
	if ok {
		t.Error("expected AddTag to return false for missing key")
	}
}
