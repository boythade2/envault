package vault

import "sort"

// DiffResult holds the comparison between two vaults.
type DiffResult struct {
	Added   []string
	Removed []string
	Changed []string
	Unchanged []string
}

// Diff compares two vaults and returns the differences.
// It compares by key presence and value equality.
func Diff(base, other *Vault) DiffResult {
	result := DiffResult{}

	baseKeys := make(map[string]string, len(base.Entries))
	for _, e := range base.Entries {
		baseKeys[e.Key] = e.Value
	}

	otherKeys := make(map[string]string, len(other.Entries))
	for _, e := range other.Entries {
		otherKeys[e.Key] = e.Value
	}

	for key, baseVal := range baseKeys {
		otherVal, exists := otherKeys[key]
		if !exists {
			result.Removed = append(result.Removed, key)
		} else if baseVal != otherVal {
			result.Changed = append(result.Changed, key)
		} else {
			result.Unchanged = append(result.Unchanged, key)
		}
	}

	for key := range otherKeys {
		if _, exists := baseKeys[key]; !exists {
			result.Added = append(result.Added, key)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Changed)
	sort.Strings(result.Unchanged)

	return result
}

// HasChanges returns true if there are any differences between the vaults.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
