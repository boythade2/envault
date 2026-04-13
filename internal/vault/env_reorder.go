package vault

import (
	"fmt"
	"os"
	"strings"
)

// ReorderResult holds the outcome of a reorder operation.
type ReorderResult struct {
	Key      string
	Position int
	Skipped  bool
	Reason   string
}

// FormatReorderResults returns a human-readable summary of reorder results.
func FormatReorderResults(results []ReorderResult, dryRun bool) string {
	if len(results) == 0 {
		return "no entries reordered"
	}
	var sb strings.Builder
	if dryRun {
		sb.WriteString("[dry-run] reorder preview:\n")
	}
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(&sb, "  skip  %s: %s\n", r.Key, r.Reason)
		} else {
			fmt.Fprintf(&sb, "  move  %s -> position %d\n", r.Key, r.Position)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// ReorderEntries moves the specified keys to the front of the vault's entry
// list (in the order supplied), leaving remaining keys in their original
// relative order. When dryRun is true the vault file is not written.
func ReorderEntries(vaultPath string, keys []string, dryRun bool) ([]ReorderResult, error) {
	v, err := LoadOrCreate(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	// Index existing entries by key.
	index := make(map[string]Entry)
	for _, e := range v.Entries {
		index[e.Key] = e
	}

	var results []ReorderResult
	seen := make(map[string]bool)
	var front []Entry

	for pos, k := range keys {
		if seen[k] {
			results = append(results, ReorderResult{Key: k, Skipped: true, Reason: "duplicate in key list"})
			continue
		}
		seen[k] = true
		e, ok := index[k]
		if !ok {
			results = append(results, ReorderResult{Key: k, Skipped: true, Reason: "key not found"})
			continue
		}
		front = append(front, e)
		results = append(results, ReorderResult{Key: k, Position: pos + 1})
	}

	// Append remaining entries preserving original order.
	for _, e := range v.Entries {
		if !seen[e.Key] {
			front = append(front, e)
		}
	}
	v.Entries = front

	if !dryRun {
		data, err := v.marshal()
		if err != nil {
			return nil, fmt.Errorf("marshal vault: %w", err)
		}
		if err := os.WriteFile(vaultPath, data, 0600); err != nil {
			return nil, fmt.Errorf("write vault: %w", err)
		}
	}
	return results, nil
}
