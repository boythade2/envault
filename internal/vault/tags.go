package vault

import (
	"sort"
	"strings"
)

// TagFilter returns entries whose Tags field contains all of the given tags.
func (v *Vault) TagFilter(tags []string) []Entry {
	if len(tags) == 0 {
		return v.Entries
	}

	var results []Entry
	for _, e := range v.Entries {
		if entryHasAllTags(e, tags) {
			results = append(results, e)
		}
	}
	return results
}

// AllTags returns a sorted, deduplicated list of every tag used across all entries.
func (v *Vault) AllTags() []string {
	seen := make(map[string]struct{})
	for _, e := range v.Entries {
		for _, t := range e.Tags {
			seen[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
		}
	}

	tags := make([]string, 0, len(seen))
	for t := range seen {
		if t != "" {
			tags = append(tags, t)
		}
	}
	sort.Strings(tags)
	return tags
}

// AddTag appends a tag to the entry identified by key, if not already present.
// Returns false if the key does not exist.
func (v *Vault) AddTag(key, tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return false
	}
	for i, e := range v.Entries {
		if e.Key == key {
			for _, existing := range e.Tags {
				if existing == tag {
					return true
				}
			}
			v.Entries[i].Tags = append(v.Entries[i].Tags, tag)
			return true
		}
	}
	return false
}

// RemoveTag removes a tag from the entry identified by key.
// Returns false if the key does not exist.
func (v *Vault) RemoveTag(key, tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for i, e := range v.Entries {
		if e.Key == key {
			updated := e.Tags[:0]
			for _, t := range e.Tags {
				if t != tag {
					updated = append(updated, t)
				}
			}
			v.Entries[i].Tags = updated
			return true
		}
	}
	return false
}

func entryHasAllTags(e Entry, tags []string) bool {
	tagSet := make(map[string]struct{}, len(e.Tags))
	for _, t := range e.Tags {
		tagSet[strings.ToLower(t)] = struct{}{}
	}
	for _, required := range tags {
		if _, ok := tagSet[strings.ToLower(required)]; !ok {
			return false
		}
	}
	return true
}
