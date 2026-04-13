package vault

import (
	"fmt"
	"os"
	"strings"
)

// PromoteResult describes the outcome of a single key promotion.
type PromoteResult struct {
	Key      string
	FromFile string
	ToFile   string
	Skipped  bool
	Reason   string
}

// PromoteOptions controls how promotion behaves.
type PromoteOptions struct {
	Keys      []string // if empty, promote all keys
	Overwrite bool
	DryRun    bool
}

// PromoteEntries copies matching entries from srcPath vault into dstPath vault.
func PromoteEntries(srcPath, dstPath string, passphrase string, opts PromoteOptions) ([]PromoteResult, error) {
	src, err := LoadOrCreate(srcPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("load source vault: %w", err)
	}

	dst, err := LoadOrCreate(dstPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("load destination vault: %w", err)
	}

	candidates := src.Entries
	if len(opts.Keys) > 0 {
		filtered := make(map[string]Entry)
		for _, k := range opts.Keys {
			upper := strings.ToUpper(k)
			if e, ok := src.Entries[upper]; ok {
				filtered[upper] = e
			} else {
				filtered[upper] = Entry{} // will be reported as missing
			}
		}
		candidates = filtered
	}

	var results []PromoteResult

	for key, entry := range candidates {
		res := PromoteResult{Key: key, FromFile: srcPath, ToFile: dstPath}

		if entry.Value == "" && entry.UpdatedAt.IsZero() {
			res.Skipped = true
			res.Reason = "key not found in source vault"
			results = append(results, res)
			continue
		}

		if _, exists := dst.Entries[key]; exists && !opts.Overwrite {
			res.Skipped = true
			res.Reason = "key already exists in destination (use --overwrite to replace)"
			results = append(results, res)
			continue
		}

		if !opts.DryRun {
			dst.Entries[key] = entry
		}
		results = append(results, res)
	}

	if !opts.DryRun && len(results) > 0 {
		if err := dst.Save(dstPath, passphrase); err != nil {
			return nil, fmt.Errorf("save destination vault: %w", err)
		}
	}

	return results, nil
}

// FormatPromoteResults renders a human-readable summary to stdout.
func FormatPromoteResults(results []PromoteResult, dryRun bool, w *os.File) {
	prefix := ""
	if dryRun {
		prefix = "[dry-run] "
	}
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(w, "%sSKIPPED  %s — %s\n", prefix, r.Key, r.Reason)
		} else {
			fmt.Fprintf(w, "%sPROMOTED %s\n", prefix, r.Key)
		}
	}
}
