package vault

import (
	"fmt"
	"strings"
)

// LowercaseResult holds the outcome of a single key lowercase operation.
type LowercaseResult struct {
	OldKey  string
	NewKey  string
	Skipped bool
	Reason  string
}

// LowercaseKeys renames all (or selected) keys to their lowercase equivalents.
// If keys is non-empty, only those keys are processed.
// If dryRun is true, the vault file is not written.
func LowercaseKeys(vaultPath, passphrase string, keys []string, dryRun bool) ([]LowercaseResult, error) {
	v, err := LoadOrCreate(vaultPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	targetAll := len(keys) == 0
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	existing := make(map[string]bool, len(v.Entries))
	for k := range v.Entries {
		existing[k] = true
	}

	var results []LowercaseResult

	for key, entry := range v.Entries {
		if !targetAll && !keySet[key] {
			continue
		}
		newKey := strings.ToLower(key)
		if newKey == key {
			results = append(results, LowercaseResult{OldKey: key, NewKey: key, Skipped: true, Reason: "already lowercase"})
			continue
		}
		if existing[newKey] {
			results = append(results, LowercaseResult{OldKey: key, NewKey: newKey, Skipped: true, Reason: "conflict with existing key"})
			continue
		}
		if !dryRun {
			v.Entries[newKey] = entry
			delete(v.Entries, key)
			existing[newKey] = true
			delete(existing, key)
		}
		results = append(results, LowercaseResult{OldKey: key, NewKey: newKey})
	}

	if !dryRun {
		if err := v.Save(vaultPath, passphrase); err != nil {
			return nil, fmt.Errorf("save vault: %w", err)
		}
	}
	return results, nil
}

// FormatLowercaseResults returns a human-readable summary.
func FormatLowercaseResults(results []LowercaseResult) string {
	if len(results) == 0 {
		return "no keys processed\n"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Skipped {
			sb.WriteString(fmt.Sprintf("SKIP  %s → %s (%s)\n", r.OldKey, r.NewKey, r.Reason))
		} else {
			sb.WriteString(fmt.Sprintf("OK    %s → %s\n", r.OldKey, r.NewKey))
		}
	}
	return sb.String()
}
