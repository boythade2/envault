package vault

import (
	"fmt"
	"strings"
)

// ReplaceResult describes the outcome of a single replacement operation.
type ReplaceResult struct {
	Key      string
	OldValue string
	NewValue string
	Changed  bool
}

// ReplaceOptions controls how value replacement is performed.
type ReplaceOptions struct {
	Keys     []string // if empty, all keys are considered
	Old      string
	New      string
	All      bool // replace all occurrences vs first only
	DryRun   bool
}

// ReplaceValues performs string replacement on vault entry values.
func ReplaceValues(v *Vault, path string, opts ReplaceOptions) ([]ReplaceResult, error) {
	if opts.Old == "" {
		return nil, fmt.Errorf("old value must not be empty")
	}

	var results []ReplaceResult

	for i, entry := range v.Entries {
		if len(opts.Keys) > 0 && !keyInSlice(entry.Key, opts.Keys) {
			continue
		}
		if !strings.Contains(entry.Value, opts.Old) {
			results = append(results, ReplaceResult{Key: entry.Key, OldValue: entry.Value, NewValue: entry.Value, Changed: false})
			continue
		}
		var newVal string
		if opts.All {
			newVal = strings.ReplaceAll(entry.Value, opts.Old, opts.New)
		} else {
			newVal = strings.Replace(entry.Value, opts.Old, opts.New, 1)
		}
		results = append(results, ReplaceResult{Key: entry.Key, OldValue: entry.Value, NewValue: newVal, Changed: true})
		if !opts.DryRun {
			v.Entries[i].Value = newVal
		}
	}

	if !opts.DryRun {
		if err := v.Save(path); err != nil {
			return nil, fmt.Errorf("save vault: %w", err)
		}
	}
	return results, nil
}

func keyInSlice(key string, list []string) bool {
	for _, k := range list {
		if k == key {
			return true
		}
	}
	return false
}

// FormatReplaceResults returns a human-readable summary.
func FormatReplaceResults(results []ReplaceResult) string {
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
			sb.WriteString(fmt.Sprintf("  ~ %s: %q -> %q\n", r.Key, r.OldValue, r.NewValue))
		}
	}
	if changed == 0 {
		return "No values matched the search string.\n"
	}
	return fmt.Sprintf("%d key(s) updated.\n", changed) + sb.String()
}
