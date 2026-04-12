package vault

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// InjectResult holds the outcome of an environment injection operation.
type InjectResult struct {
	Injected []string
	Skipped  []string
	Overridden []string
}

// InjectOptions controls behaviour of InjectEnv.
type InjectOptions struct {
	// Overwrite existing OS env vars when true.
	Overwrite bool
	// Prefix filters which vault keys are injected (empty = all).
	Prefix string
}

// InjectEnv sets OS environment variables from the vault entries.
// It returns a summary of what was injected, skipped, or overridden.
func InjectEnv(v *Vault, opts InjectOptions) (InjectResult, error) {
	var result InjectResult

	keys := make([]string, 0, len(v.Entries))
	for k := range v.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		_, exists := os.LookupEnv(k)
		if exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if err := os.Setenv(k, v.Entries[k].Value); err != nil {
			return result, fmt.Errorf("setenv %s: %w", k, err)
		}

		if exists {
			result.Overridden = append(result.Overridden, k)
		} else {
			result.Injected = append(result.Injected, k)
		}
	}

	return result, nil
}

// FormatInjectResult returns a human-readable summary of an InjectResult.
func FormatInjectResult(r InjectResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Injected: %d, Overridden: %d, Skipped: %d\n",
		len(r.Injected), len(r.Overridden), len(r.Skipped)))
	for _, k := range r.Injected {
		sb.WriteString(fmt.Sprintf("  + %s\n", k))
	}
	for _, k := range r.Overridden {
		sb.WriteString(fmt.Sprintf("  ~ %s (overridden)\n", k))
	}
	return sb.String()
}
