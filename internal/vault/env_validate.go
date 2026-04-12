package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// EnvValidationResult holds the outcome of validating environment variable names.
type EnvValidationResult struct {
	Key     string
	Warning string
	Error   string
}

var validEnvKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)
var hasLowercase = regexp.MustCompile(`[a-z]`)
var hasSpace = regexp.MustCompile(`\s`)

// ValidateEnvKeys checks all keys in a vault file against common environment
// variable naming conventions. Returns a slice of results with warnings or
// errors for any key that does not conform.
func ValidateEnvKeys(path string) ([]EnvValidationResult, error) {
	v, err := LoadOrCreate(path)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	var results []EnvValidationResult

	for _, entry := range v.Entries {
		key := entry.Key
		res := EnvValidationResult{Key: key}

		switch {
		case strings.TrimSpace(key) == "":
			res.Error = "key is empty or blank"
		case hasSpace.MatchString(key):
			res.Error = "key contains whitespace"
		case strings.HasPrefix(key, "0") || (len(key) > 0 && key[0] >= '0' && key[0] <= '9'):
			res.Error = "key starts with a digit"
		case hasLowercase.MatchString(key):
			res.Warning = "key contains lowercase letters; consider uppercasing"
		case !validEnvKeyPattern.MatchString(key):
			res.Warning = "key contains characters outside [A-Z0-9_]"
		}

		if res.Warning != "" || res.Error != "" {
			results = append(results, res)
		}
	}

	return results, nil
}

// FormatEnvValidationResults returns a human-readable summary string.
func FormatEnvValidationResults(results []EnvValidationResult) string {
	if len(results) == 0 {
		return "all keys are valid"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Error != "" {
			fmt.Fprintf(&sb, "[ERROR] %s: %s\n", r.Key, r.Error)
		} else {
			fmt.Fprintf(&sb, "[WARN]  %s: %s\n", r.Key, r.Warning)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
