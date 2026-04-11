package vault

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting warning or error for a vault entry.
type LintIssue struct {
	Key      string
	Severity string // "warn" or "error"
	Message  string
}

func (l LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(l.Severity), l.Key, l.Message)
}

// LintResult holds all issues found during a lint pass.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// LintVault inspects all entries in the vault and returns a LintResult
// containing warnings and errors based on common best-practice rules.
func LintVault(v *Vault) LintResult {
	result := LintResult{}

	for _, entry := range v.Entries {
		key := entry.Key

		// Error: empty key
		if strings.TrimSpace(key) == "" {
			result.Issues = append(result.Issues, LintIssue{
				Key:      "(empty)",
				Severity: "error",
				Message:  "key must not be empty",
			})
			continue
		}

		// Error: key contains spaces
		if strings.Contains(key, " ") {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "error",
				Message:  "key contains spaces, which are not valid in most shells",
			})
		}

		// Warn: key is not uppercase
		if key != strings.ToUpper(key) {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "warn",
				Message:  "key is not uppercase; convention recommends ALL_CAPS",
			})
		}

		// Warn: empty value
		if strings.TrimSpace(entry.Value) == "" {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "warn",
				Message:  "value is empty",
			})
		}

		// Warn: key starts with a digit
		if len(key) > 0 && key[0] >= '0' && key[0] <= '9' {
			result.Issues = append(result.Issues, LintIssue{
				Key:      key,
				Severity: "warn",
				Message:  "key starts with a digit, which is invalid in POSIX shells",
			})
		}
	}

	return result
}
