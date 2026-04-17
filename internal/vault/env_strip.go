package vault

import (
	"fmt"
	"strings"
)

// StripResult holds the outcome for a single key strip operation.
type StripResult struct {
	Key     string
	OldVal  string
	NewVal  string
	Changed bool
	Skipped bool
}

// StripOptions controls how StripEntries behaves.
type StripOptions struct {
	Keys     []string // empty = all keys
	Chars    string   // characters to strip (default: whitespace)
	DryRun   bool
}

// StripEntries removes leading/trailing characters from values in the vault.
func StripEntries(vaultPath string, opts StripOptions) ([]StripResult, error) {
	v, err := LoadOrCreate(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	targetAll := len(opts.Keys) == 0
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	var results []StripResult
	for i, e := range v.Entries {
		if !targetAll && !keySet[e.Key] {
			continue
		}
		var newVal string
		if opts.Chars == "" {
			newVal = strings.TrimSpace(e.Value)
		} else {
			newVal = strings.Trim(e.Value, opts.Chars)
		}
		changed := newVal != e.Value
		results = append(results, StripResult{
			Key:     e.Key,
			OldVal:  e.Value,
			NewVal:  newVal,
			Changed: changed,
		})
		if changed && !opts.DryRun {
			v.Entries[i].Value = newVal
		}
	}

	if !opts.DryRun {
		if err := v.Save(vaultPath); err != nil {
			return nil, fmt.Errorf("save vault: %w", err)
		}
	}
	return results, nil
}

// FormatStripResults returns a human-readable summary.
func FormatStripResults(results []StripResult) string {
	if len(results) == 0 {
		return "no entries matched"
	}
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
			fmt.Fprintf(&sb, "  stripped  %s\n", r.Key)
		}
	}
	if changed == 0 {
		return "no values required stripping"
	}
	return strings.TrimRight(sb.String(), "\n")
}
