package vault

import (
	"fmt"
	"strings"
)

// DedupeResult holds the outcome of a deduplication pass.
type DedupeResult struct {
	Key      string
	Kept     string
	Removed  []string
	Strategy string
}

// DedupeOptions controls how duplicates are resolved.
type DedupeOptions struct {
	// Strategy: "first" keeps the earliest entry, "last" keeps the latest.
	Strategy string
	DryRun   bool
}

// DedupeEntries removes duplicate keys from the vault according to the chosen
// strategy. Duplicates are identified by case-insensitive key comparison.
func DedupeEntries(v *Vault, opts DedupeOptions) ([]DedupeResult, error) {
	if opts.Strategy == "" {
		opts.Strategy = "first"
	}
	if opts.Strategy != "first" && opts.Strategy != "last" {
		return nil, fmt.Errorf("unknown strategy %q: must be \"first\" or \"last\"", opts.Strategy)
	}

	// Group indices by normalised key.
	order := []string{}
	groups := map[string][]int{}
	for i, e := range v.Entries {
		norm := strings.ToUpper(e.Key)
		if _, seen := groups[norm]; !seen {
			order = append(order, norm)
		}
		groups[norm] = append(groups[norm], i)
	}

	var results []DedupeResult
	var kept []Entry

	for _, norm := range order {
		idxs := groups[norm]
		if len(idxs) == 1 {
			kept = append(kept, v.Entries[idxs[0]])
			continue
		}

		var keepIdx int
		if opts.Strategy == "last" {
			keepIdx = idxs[len(idxs)-1]
		} else {
			keepIdx = idxs[0]
		}

		var removed []string
		for _, i := range idxs {
			if i != keepIdx {
				removed = append(removed, v.Entries[i].Value)
			}
		}

		results = append(results, DedupeResult{
			Key:      v.Entries[keepIdx].Key,
			Kept:     v.Entries[keepIdx].Value,
			Removed:  removed,
			Strategy: opts.Strategy,
		})
		kept = append(kept, v.Entries[keepIdx])
	}

	if !opts.DryRun {
		v.Entries = kept
	}
	return results, nil
}

// FormatDedupeResults returns a human-readable summary.
func FormatDedupeResults(results []DedupeResult) string {
	if len(results) == 0 {
		return "No duplicate keys found."
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "[%s] kept %q, removed %d duplicate(s)\n",
			r.Key, r.Kept, len(r.Removed))
	}
	return strings.TrimRight(sb.String(), "\n")
}
