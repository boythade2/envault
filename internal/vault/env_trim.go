package vault

import (
	"fmt"
	"strings"
)

// TrimResult holds the result of trimming a single entry.
type TrimResult struct {
	Key      string
	OldValue string
	NewValue string
	Changed  bool
}

// TrimOptions controls how trimming is applied.
type TrimOptions struct {
	Keys    []string // empty means all keys
	Left    bool
	Right   bool
	DryRun  bool
}

// TrimEntries removes leading/trailing whitespace from vault entry values.
func TrimEntries(v *Vault, path string, opts TrimOptions) ([]TrimResult, error) {
	var results []TrimResult

	target := map[string]bool{}
	for _, k := range opts.Keys {
		target[k] = true
	}

	for i, e := range v.Entries {
		if len(target) > 0 && !target[e.Key] {
			continue
		}
		newVal := e.Value
		if opts.Left && opts.Right {
			newVal = strings.TrimSpace(e.Value)
		} else if opts.Left {
			newVal = strings.TrimLeft(e.Value, " \t")
		} else if opts.Right {
			newVal = strings.TrimRight(e.Value, " \t")
		} else {
			newVal = strings.TrimSpace(e.Value)
		}
		changed := newVal != e.Value
		results = append(results, TrimResult{
			Key:      e.Key,
			OldValue: e.Value,
			NewValue: newVal,
			Changed:  changed,
		})
		if changed && !opts.DryRun {
			v.Entries[i].Value = newVal
		}
	}

	if !opts.DryRun {
		if err := v.Save(path); err != nil {
			return nil, err
		}
	}
	return results, nil
}

// FormatTrimResults returns a human-readable summary of trim results.
func FormatTrimResults(results []TrimResult, dryRun bool) string {
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
			prefix := "trimmed"
			if dryRun {
				prefix = "would trim"
			}
			sb.WriteString(fmt.Sprintf("  %s: %s\n", prefix, r.Key))
		}
	}
	if changed == 0 {
		return "no entries required trimming\n"
	}
	return fmt.Sprintf("%d entr(ies) trimmed\n", changed) + sb.String()
}
