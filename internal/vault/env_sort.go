package vault

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering direction for vault entries.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// SortField defines the field to sort vault entries by.
type SortField string

const (
	SortByKey       SortField = "key"
	SortByValue     SortField = "value"
	SortByUpdatedAt SortField = "updated_at"
)

// SortOptions configures how vault entries are sorted.
type SortOptions struct {
	Field SortField
	Order SortOrder
}

// SortedEntry pairs a key with its vault entry for sorting.
type SortedEntry struct {
	Key   string
	Value Entry
}

// SortEntries returns vault entries sorted according to the given options.
func SortEntries(v *Vault, opts SortOptions) []SortedEntry {
	entries := make([]SortedEntry, 0, len(v.Entries))
	for k, e := range v.Entries {
		entries = append(entries, SortedEntry{Key: k, Value: e})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		var less bool
		switch opts.Field {
		case SortByValue:
			less = strings.ToLower(entries[i].Value.Value) < strings.ToLower(entries[j].Value.Value)
		case SortByUpdatedAt:
			less = entries[i].Value.UpdatedAt.Before(entries[j].Value.UpdatedAt)
		default: // SortByKey
			less = strings.ToLower(entries[i].Key) < strings.ToLower(entries[j].Key)
		}
		if opts.Order == SortDesc {
			return !less
		}
		return less
	})

	return entries
}
