package vault

import (
	"fmt"
	"os"
	"strings"
)

// SquashResult holds the outcome of a single squash operation.
type SquashResult struct {
	Key     string
	OldVal  string
	NewVal  string
	Skipped bool
	Reason  string
}

// SquashOptions controls how SquashEntries behaves.
type SquashOptions struct {
	// Keys to squash; if empty all keys are considered.
	Keys []string
	// Transform is one of: "trim", "lower", "upper", "collapse" (collapse whitespace).
	Transform string
	DryRun    bool
}

// SquashEntries applies a normalisation transform to entry values in the vault.
func SquashEntries(vaultFile string, opts SquashOptions) ([]SquashResult, error) {
	if opts.Transform == "" {
		return nil, fmt.Errorf("transform must be one of: trim, lower, upper, collapse")
	}

	v, err := LoadOrCreate(vaultFile)
	if err != nil {
		return nil, fmt.Errorf("load vault: %w", err)
	}

	target := map[string]bool{}
	for _, k := range opts.Keys {
		target[k] = true
	}

	var results []SquashResult
	for i, e := range v.Entries {
		if len(target) > 0 && !target[e.Key] {
			continue
		}

		newVal, err := applySquashTransform(e.Value, opts.Transform)
		if err != nil {
			results = append(results, SquashResult{Key: e.Key, OldVal: e.Value, Skipped: true, Reason: err.Error()})
			continue
		}

		if newVal == e.Value {
			results = append(results, SquashResult{Key: e.Key, OldVal: e.Value, NewVal: newVal, Skipped: true, Reason: "no change"})
			continue
		}

		results = append(results, SquashResult{Key: e.Key, OldVal: e.Value, NewVal: newVal})
		if !opts.DryRun {
			v.Entries[i].Value = newVal
		}
	}

	if !opts.DryRun {
		if err := v.Save(vaultFile); err != nil {
			return nil, fmt.Errorf("save vault: %w", err)
		}
	}
	return results, nil
}

func applySquashTransform(val, transform string) (string, error) {
	switch transform {
	case "trim":
		return strings.TrimSpace(val), nil
	case "lower":
		return strings.ToLower(val), nil
	case "upper":
		return strings.ToUpper(val), nil
	case "collapse":
		fields := strings.Fields(val)
		return strings.Join(fields, " "), nil
	default:
		return "", fmt.Errorf("unknown transform %q", transform)
	}
}

// FormatSquashResults returns a human-readable summary.
func FormatSquashResults(results []SquashResult) string {
	if len(results) == 0 {
		return "no entries matched"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Skipped {
			fmt.Fprintf(&sb, "  skip  %-24(%s)\n", r.Key, r.Reason)
		} else {
			fmt.Fprintf(&sb, "  squash %-24s  %q -> %q\n", r.Key, r.OldVal, r.NewVal)
		}
	}
	_ = os.Stderr // suppress unused import
	return sb.String()
}
