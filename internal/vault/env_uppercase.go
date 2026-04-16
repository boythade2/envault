package vault

import (
	"fmt"
	"os"
	"strings"
)

// UppercaseResult holds the result of a single key normalization.
type UppercaseResult struct {
	Key     string
	OldKey  string
	Skipped bool
	Reason  string
}

// UppercaseKeys renames all (or selected) keys to their uppercase form.
// If keys is non-empty, only those keys are affected.
// If dryRun is true, the vault file is not written.
func UppercaseKeys(vaultPath, passphrase string, keys []string, dryRun bool) ([]UppercaseResult, error) {
	v, err := LoadOrCreate(vaultPath, passphrase)
	if err != nil {
		return nil, err
	}

	target := map[string]bool{}
	for _, k := range keys {
		target[k] = true
	}

	var results []UppercaseResult

	for i, e := range v.Entries {
		if len(target) > 0 && !target[e.Key] {
			continue
		}
		upper := strings.ToUpper(e.Key)
		if upper == e.Key {
			results = append(results, UppercaseResult{Key: e.Key, OldKey: e.Key, Skipped: true, Reason: "already uppercase"})
			continue
		}
		// Check for conflict
		for _, other := range v.Entries {
			if other.Key == upper {
				results = append(results, UppercaseResult{Key: e.Key, OldKey: e.Key, Skipped: true, Reason: fmt.Sprintf("conflict with existing key %q", upper)})
				goto next
			}
		}
		v.Entries[i].Key = upper
		results = append(results, UppercaseResult{Key: upper, OldKey: e.Key})
	next:
	}

	if !dryRun {
		data, err := v.serialize(passphrase)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(vaultPath, data, 0600); err != nil {
			return nil, err
		}
	}
	return results, nil
}

// FormatUppercaseResults returns a human-readable summary.
func FormatUppercaseResults(results []UppercaseResult) string {
	if len(results) == 0 {
		return "no keys matched"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(&sb, "  SKIP  %s (%s)\n", r.OldKey, r.Reason)
		} else {
			fmt.Fprintf(&sb, "  OK    %s -> %s\n", r.OldKey, r.Key)
		}
	}
	return sb.String()
}
