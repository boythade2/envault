package vault

import (
	"fmt"
	"os"
	"strings"
)

// EnvRenameResult describes the outcome of a bulk rename operation.
type EnvRenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
	Reason  string
}

// BulkRenameOptions controls how bulk rename behaves.
type BulkRenameOptions struct {
	Prefix      string // strip or add a prefix
	AddPrefix   string
	StripPrefix string
	DryRun      bool
	Overwrite   bool
}

// BulkRenameKeys renames vault keys according to the provided options.
// It returns a list of results for each key considered.
func BulkRenameKeys(vaultPath string, opts BulkRenameOptions) ([]EnvRenameResult, error) {
	v, err := LoadOrCreate(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	var results []EnvRenameResult

	for key, entry := range v.Entries {
		newKey := key

		if opts.StripPrefix != "" && strings.HasPrefix(key, opts.StripPrefix) {
			newKey = strings.TrimPrefix(key, opts.StripPrefix)
		}

		if opts.AddPrefix != "" {
			newKey = opts.AddPrefix + newKey
		}

		if newKey == key {
			results = append(results, EnvRenameResult{OldKey: key, NewKey: key, Renamed: false, Reason: "no change"})
			continue
		}

		if _, exists := v.Entries[newKey]; exists && !opts.Overwrite {
			results = append(results, EnvRenameResult{OldKey: key, NewKey: newKey, Renamed: false, Reason: "target key already exists"})
			continue
		}

		results = append(results, EnvRenameResult{OldKey: key, NewKey: newKey, Renamed: true})

		if !opts.DryRun {
			v.Entries[newKey] = entry
			delete(v.Entries, key)
		}
	}

	if !opts.DryRun {
		if err := v.Save(vaultPath, os.FileMode(0600)); err != nil {
			return nil, fmt.Errorf("save vault: %w", err)
		}
	}

	return results, nil
}

// FormatBulkRenameResults returns a human-readable summary of rename results.
func FormatBulkRenameResults(results []EnvRenameResult) string {
	var sb strings.Builder
	renamed := 0
	for _, r := range results {
		if r.Renamed {
			sb.WriteString(fmt.Sprintf("  renamed: %s -> %s\n", r.OldKey, r.NewKey))
			renamed++
		} else {
			sb.WriteString(fmt.Sprintf("  skipped: %s (%s)\n", r.OldKey, r.Reason))
		}
	}
	sb.WriteString(fmt.Sprintf("\n%d key(s) renamed.\n", renamed))
	return sb.String()
}
