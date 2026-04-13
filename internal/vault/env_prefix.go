package vault

import (
	"fmt"
	"strings"
)

// PrefixResult holds the outcome of a single key prefix operation.
type PrefixResult struct {
	OldKey  string
	NewKey  string
	Skipped bool
	Reason  string
}

// PrefixOptions controls the behaviour of AddKeyPrefix / RemoveKeyPrefix.
type PrefixOptions struct {
	Prefix    string
	DryRun    bool
	Overwrite bool
}

// AddKeyPrefix prepends Prefix to every key in the vault (or a selected subset).
// When keys is empty, all entries are processed.
func AddKeyPrefix(v *Vault, keys []string, opts PrefixOptions) ([]PrefixResult, error) {
	if opts.Prefix == "" {
		return nil, fmt.Errorf("prefix must not be empty")
	}
	targets := selectKeys(v, keys)
	var results []PrefixResult
	for _, k := range targets {
		newKey := opts.Prefix + k
		if _, exists := v.Entries[newKey]; exists && !opts.Overwrite {
			results = append(results, PrefixResult{OldKey: k, NewKey: newKey, Skipped: true, Reason: "key already exists"})
			continue
		}
		if !opts.DryRun {
			v.Entries[newKey] = v.Entries[k]
			deleten		}
		results = append(results, PrefixResult{OldKey: k, NewKey: newKey})
	}
	return results, nil
}

// RemoveKeyPrefix strips Prefix from every matching key in the vault.
// When keys is empty, all entries whose key starts with Prefix are processed.
func RemoveKeyPrefix(v *Vault, keys []string, opts PrefixOptions) ([]PrefixResult, error) {
	if opts.Prefix == "" {
		return nil, fmt.Errorf("prefix must not be empty")
	}
	targets := selectKeys(v, keys)
	var results []PrefixResult
	for _, k := range targets {
		if !strings.HasPrefix(k, opts.Prefix) {
			results = append(results, PrefixResult{OldKey: k, NewKey: k, Skipped: true, Reason: "prefix not present"})
			continue
		}
		newKey := strings.TrimPrefix(k, opts.Prefix)
		if newKey == "" {
			results = append(results, PrefixResult{OldKey: k, NewKey: k, Skipped: true, Reason: "stripping prefix yields empty key"})
			continue
		}
		if _, exists := v.Entries[newKey]; exists && !opts.Overwrite {
			results = append(results, PrefixResult{OldKey: k, NewKey: newKey, Skipped: true, Reason: "key already exists"})
			continue
		}
		if !opts.DryRun {
			v.Entries[newKey] = v.Entries[k]
			delete(v.Entries, k)
		}
		results = append(results, PrefixResult{OldKey: k, NewKey: newKey})
	}
	return results, nil
}

// FormatPrefixResults returns a human-readable summary.
func FormatPrefixResults(results []PrefixResult) string {
	if len(results) == 0 {
		return "no keys processed"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(&sb, "SKIP  %s  (%s)\n", r.OldKey, r.Reason)
		} else {
			fmt.Fprintf(&sb, "OK    %s  ->  %s\n", r.OldKey, r.NewKey)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// selectKeys returns the requested keys that exist in the vault.
// If the provided list is empty, all vault keys are returned.
func selectKeys(v *Vault, keys []string) []string {
	if len(keys) == 0 {
		out := make([]string, 0, len(v.Entries))
		for k := range v.Entries {
			out = append(out, k)
		}
		return out
	}
	var out []string
	for _, k := range keys {
		if _, ok := v.Entries[k]; ok {
			out = append(out, k)
		}
	}
	return out
}
