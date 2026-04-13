package vault

import (
	"regexp"
	"strings"
)

// MaskResult holds the result of masking a single entry.
type MaskResult struct {
	Key    string
	Masked bool
}

// MaskOptions controls how values are masked.
type MaskOptions struct {
	// ShowChars reveals this many leading characters before masking.
	ShowChars int
	// MaskChar is the character used for masking (default '*').
	MaskChar string
	// Patterns is a list of key regex patterns to mask; empty means mask all.
	Patterns []string
}

// MaskValues returns a copy of entries with sensitive values masked.
// Entries whose keys match any of opts.Patterns are masked; if Patterns is
// empty every entry is masked.
func MaskValues(entries []Entry, opts MaskOptions) ([]Entry, []MaskResult) {
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}

	var compiled []*regexp.Regexp
	for _, p := range opts.Patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}

	masked := make([]Entry, len(entries))
	results := make([]MaskResult, len(entries))

	for i, e := range entries {
		masked[i] = e
		shouldMask := len(compiled) == 0
		for _, re := range compiled {
			if re.MatchString(e.Key) {
				shouldMask = true
				break
			}
		}
		results[i] = MaskResult{Key: e.Key, Masked: shouldMask}
		if shouldMask {
			masked[i].Value = maskString(e.Value, opts.ShowChars, opts.MaskChar)
		}
	}
	return masked, results
}

func maskString(s string, showChars int, maskChar string) string {
	if showChars <= 0 || showChars >= len(s) {
		return strings.Repeat(maskChar, max(len(s), 8))
	}
	return s[:showChars] + strings.Repeat(maskChar, len(s)-showChars)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
