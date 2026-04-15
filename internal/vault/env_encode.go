package vault

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EncodeFormat represents the encoding format to apply.
type EncodeFormat string

const (
	EncodeBase64    EncodeFormat = "base64"
	EncodeBase64URL EncodeFormat = "base64url"
	EncodeHex       EncodeFormat = "hex"
)

// EncodeResult holds the outcome of encoding a single entry.
type EncodeResult struct {
	Key      string
	Original string
	Encoded  string
	Skipped  bool
	Reason   string
}

// EncodeEntries encodes the values of the specified keys (or all keys if none
// specified) using the given format and writes the result back to the vault
// file unless DryRun is true.
func EncodeEntries(vaultPath, passphrase string, keys []string, format EncodeFormat, dryRun bool) ([]EncodeResult, error) {
	v, err := LoadOrCreate(vaultPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	targets := make(map[string]bool)
	for _, k := range keys {
		targets[k] = true
	}

	var results []EncodeResult
	for i, e := range v.Entries {
		if len(targets) > 0 && !targets[e.Key] {
			continue
		}
		encoded, err := encodeValue(e.Value, format)
		if err != nil {
			results = append(results, EncodeResult{Key: e.Key, Original: e.Value, Skipped: true, Reason: err.Error()})
			continue
		}
		results = append(results, EncodeResult{Key: e.Key, Original: e.Value, Encoded: encoded})
		if !dryRun {
			v.Entries[i].Value = encoded
		}
	}

	if !dryRun {
		if err := v.Save(vaultPath, passphrase); err != nil {
			return nil, fmt.Errorf("save vault: %w", err)
		}
	}
	return results, nil
}

func encodeValue(value string, format EncodeFormat) (string, error) {
	switch format {
	case EncodeBase64:
		return base64.StdEncoding.EncodeToString([]byte(value)), nil
	case EncodeBase64URL:
		return base64.URLEncoding.EncodeToString([]byte(value)), nil
	case EncodeHex:
		var sb strings.Builder
		for _, b := range []byte(value) {
			fmt.Fprintf(&sb, "%02x", b)
		}
		return sb.String(), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// FormatEncodeResults returns a human-readable summary of encode results.
func FormatEncodeResults(results []EncodeResult, dryRun bool) string {
	if len(results) == 0 {
		return "no entries matched"
	}
	var sb strings.Builder
	if dryRun {
		sb.WriteString("[dry-run] encoding preview:\n")
	}
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(&sb, "  SKIP  %s — %s\n", r.Key, r.Reason)
		} else {
			fmt.Fprintf(&sb, "  OK    %s: %s => %s\n", r.Key, r.Original, r.Encoded)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
