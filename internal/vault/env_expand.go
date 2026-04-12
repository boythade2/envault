package vault

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ExpandResult holds the result of a variable expansion operation.
type ExpandResult struct {
	Key      string
	Original string
	Expanded string
	Changed  bool
}

var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// ExpandVaultRefs expands ${KEY} references within vault entry values,
// resolving them against other entries in the same vault. If useOS is true,
// unresolved references fall back to OS environment variables.
func ExpandVaultRefs(v *Vault, useOS bool) ([]ExpandResult, error) {
	results := make([]ExpandResult, 0, len(v.Entries))

	for key, entry := range v.Entries {
		original := entry.Value
		expanded, err := expandValue(original, v, useOS)
		if err != nil {
			return nil, fmt.Errorf("expanding %q: %w", key, err)
		}
		results = append(results, ExpandResult{
			Key:      key,
			Original: original,
			Expanded: expanded,
			Changed:  expanded != original,
		})
	}

	return results, nil
}

func expandValue(value string, v *Vault, useOS bool) (string, error) {
	var expandErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		ref := refPattern.FindStringSubmatch(match)
		if len(ref) < 2 {
			return match
		}
		refKey := ref[1]
		if entry, ok := v.Entries[refKey]; ok {
			return entry.Value
		}
		if useOS {
			if val, ok := os.LookupEnv(refKey); ok {
				return val
			}
		}
		expandErr = fmt.Errorf("unresolved reference: %s", refKey)
		return match
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// FormatExpandResults returns a human-readable summary of expansion results.
func FormatExpandResults(results []ExpandResult) string {
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
			fmt.Fprintf(&sb, "  %s: %q -> %q\n", r.Key, r.Original, r.Expanded)
		}
	}
	if changed == 0 {
		sb.WriteString("  no references expanded\n")
	}
	return sb.String()
}
