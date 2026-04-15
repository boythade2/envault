package vault

import (
	"fmt"
	"os"
	"strings"
)

// PlaceholderResult describes the outcome for a single key during placeholder resolution.
type PlaceholderResult struct {
	Key      string
	Original string
	Resolved string
	Changed  bool
	Missing  []string
}

// ResolvePlaceholders replaces {{KEY}} style placeholders in vault values with
// other vault values or OS environment variables. Keys listed in onlyKeys are
// processed; if onlyKeys is empty all entries are processed.
func ResolvePlaceholders(v *Vault, onlyKeys []string, dryRun bool) ([]PlaceholderResult, error) {
	if v == nil {
		return nil, fmt.Errorf("vault is nil")
	}

	index := make(map[string]string, len(v.Entries))
	for _, e := range v.Entries {
		index[e.Key] = e.Value
	}

	filter := make(map[string]bool, len(onlyKeys))
	for _, k := range onlyKeys {
		filter[k] = true
	}

	var results []PlaceholderResult

	for i, e := range v.Entries {
		if len(filter) > 0 && !filter[e.Key] {
			continue
		}
		resolved, missing := expandPlaceholders(e.Value, index)
		changed := resolved != e.Value
		results = append(results, PlaceholderResult{
			Key:      e.Key,
			Original: e.Value,
			Resolved: resolved,
			Changed:  changed,
			Missing:  missing,
		})
		if changed && !dryRun {
			v.Entries[i].Value = resolved
		}
	}
	return results, nil
}

// expandPlaceholders replaces all {{KEY}} tokens in s using the provided index
// and OS env as fallback. Returns the expanded string and any unresolved keys.
func expandPlaceholders(s string, index map[string]string) (string, []string) {
	var missing []string
	result := s
	for {
		start := strings.Index(result, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}}")
		if end == -1 {
			break
		}
		end += start
		token := result[start+2 : end]
		key := strings.TrimSpace(token)
		var replacement string
		if val, ok := index[key]; ok {
			replacement = val
		} else if val := os.Getenv(key); val != "" {
			replacement = val
		} else {
			missing = append(missing, key)
			replacement = "{{" + token + "}}"
			// advance past this token to avoid infinite loop
			result = result[:start] + replacement + result[end+2:]
			continue
		}
		result = result[:start] + replacement + result[end+2:]
	}
	return result, missing
}

// FormatPlaceholderResults returns a human-readable summary.
func FormatPlaceholderResults(results []PlaceholderResult) string {
	if len(results) == 0 {
		return "no entries processed\n"
	}
	var sb strings.Builder
	for _, r := range results {
		if !r.Changed && len(r.Missing) == 0 {
			continue
		}
		if len(r.Missing) > 0 {
			sb.WriteString(fmt.Sprintf("WARN  %s: unresolved placeholders: %s\n", r.Key, strings.Join(r.Missing, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("OK    %s: %q -> %q\n", r.Key, r.Original, r.Resolved))
		}
	}
	if sb.Len() == 0 {
		return "all placeholders already resolved\n"
	}
	return sb.String()
}
