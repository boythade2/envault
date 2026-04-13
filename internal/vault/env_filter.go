package vault

import (
	"strings"
)

// FilterResult holds the outcome of a filter operation.
type FilterResult struct {
	Key   string
	Value string
}

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	KeyPrefix    string
	KeySuffix    string
	ValueContains string
	TagFilter    []string
	InvertMatch  bool
}

// FilterEntries returns vault entries matching the given options.
func FilterEntries(v *Vault, opts FilterOptions) []FilterResult {
	var results []FilterResult

	for _, entry := range v.Entries {
		matched := matchesFilter(entry, opts)
		if opts.InvertMatch {
			matched = !matched
		}
		if matched {
			results = append(results, FilterResult{
				Key:   entry.Key,
				Value: entry.Value,
			})
		}
	}

	return results
}

func matchesFilter(entry Entry, opts FilterOptions) bool {
	if opts.KeyPrefix != "" && !strings.HasPrefix(entry.Key, opts.KeyPrefix) {
		return false
	}
	if opts.KeySuffix != "" && !strings.HasSuffix(entry.Key, opts.KeySuffix) {
		return false
	}
	if opts.ValueContains != "" && !strings.Contains(entry.Value, opts.ValueContains) {
		return false
	}
	if len(opts.TagFilter) > 0 && !entryHasAllTags(entry, opts.TagFilter) {
		return false
	}
	return true
}

// FormatFilterResults formats filter results as a human-readable string.
func FormatFilterResults(results []FilterResult) string {
	if len(results) == 0 {
		return "no entries matched the filter criteria\n"
	}
	var sb strings.Builder
	sb.WriteString("KEY\tVALUE\n")
	sb.WriteString("---\t-----\n")
	for _, r := range results {
		sb.WriteString(r.Key + "\t" + r.Value + "\n")
	}
	return sb.String()
}
