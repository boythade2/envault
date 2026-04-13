package vault

import (
	"fmt"
	"strings"
)

// TruncateResult holds the outcome of a single truncation operation.
type TruncateResult struct {
	Key      string
	Original string
	Truncated string
	Skipped  bool
	Reason   string
}

// TruncateOptions controls how values are truncated.
type TruncateOptions struct {
	MaxLen  int
	Suffix  string
	Keys    []string // if empty, all keys are considered
	DryRun  bool
}

// TruncateEntries shortens vault entry values that exceed MaxLen.
func TruncateEntries(v *Vault, opts TruncateOptions) ([]TruncateResult, error) {
	if opts.MaxLen <= 0 {
		return nil, fmt.Errorf("max-len must be greater than zero")
	}
	suffix := opts.Suffix
	if suffix == "" {
		suffix = "..."
	}

	filterKeys := map[string]bool{}
	for _, k := range opts.Keys {
		filterKeys[k] = true
	}

	var results []TruncateResult
	for i, e := range v.Entries {
		if len(filterKeys) > 0 && !filterKeys[e.Key] {
			continue
		}
		if len(e.Value) <= opts.MaxLen {
			results = append(results, TruncateResult{
				Key:      e.Key,
				Original: e.Value,
				Truncated: e.Value,
				Skipped:  true,
				Reason:   "value within limit",
			})
			continue
		}
		cutAt := opts.MaxLen - len(suffix)
		if cutAt < 0 {
			cutAt = 0
		}
		truncated := e.Value[:cutAt] + suffix
		results = append(results, TruncateResult{
			Key:      e.Key,
			Original: e.Value,
			Truncated: truncated,
		})
		if !opts.DryRun {
			v.Entries[i].Value = truncated
		}
	}
	return results, nil
}

// FormatTruncateResults returns a human-readable summary.
func FormatTruncateResults(results []TruncateResult) string {
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if !r.Skipped {
			changed++
			sb.WriteString(fmt.Sprintf("  truncated  %s  (%d -> %d chars)\n",
				r.Key, len(r.Original), len(r.Truncated)))
		}
	}
	if changed == 0 {
		return "No values required truncation.\n"
	}
	return fmt.Sprintf("%d value(s) truncated.\n", changed) + sb.String()
}
