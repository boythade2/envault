package vault

import (
	"fmt"
	"regexp"

	"github.com/user/envault/internal/vault"
)

// RegexFilterResult holds the outcome for a single entry.
type RegexFilterResult struct {
	Key     string
	Value   string
	Matched bool
}

// RegexFilter returns entries whose key or value matches the given pattern.
func RegexFilter(v *Vault, pattern string, matchValue bool) ([]RegexFilterResult, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}
	var results []RegexFilterResult
	for _, e := range v.Entries {
		matched := re.MatchString(e.Key)
		if matchValue {
			matched = matched || re.MatchString(e.Value)
		}
		if matched {
			results = append(results, RegexFilterResult{Key: e.Key, Value: e.Value, Matched: true})
		}
	}
	return results, nil
}

// FormatRegexResults formats matched entries for CLI output.
func FormatRegexResults(results []RegexFilterResult) string {
	if len(results) == 0 {
		return "no entries matched"
	}
	out := fmt.Sprintf("%-30s %s\n", "KEY", "VALUE")
	for _, r := range results {
		out += fmt.Sprintf("%-30s %s\n", r.Key, r.Value)
	}
	return out
}
