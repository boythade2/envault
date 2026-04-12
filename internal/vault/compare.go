package vault

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two vault files.
type CompareResult struct {
	OnlyInA  []string
	OnlyInB  []string
	Changed  []string
	Identical []string
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	var sb strings.Builder
	if len(r.OnlyInA) > 0 {
		sb.WriteString(fmt.Sprintf("Only in A (%d): %s\n", len(r.OnlyInA), strings.Join(r.OnlyInA, ", ")))
	}
	if len(r.OnlyInB) > 0 {
		sb.WriteString(fmt.Sprintf("Only in B (%d): %s\n", len(r.OnlyInB), strings.Join(r.OnlyInB, ", ")))
	}
	if len(r.Changed) > 0 {
		sb.WriteString(fmt.Sprintf("Changed (%d): %s\n", len(r.Changed), strings.Join(r.Changed, ", ")))
	}
	if len(r.Identical) > 0 {
		sb.WriteString(fmt.Sprintf("Identical (%d): %s\n", len(r.Identical), strings.Join(r.Identical, ", ")))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// CompareVaults compares the entries of two vaults and returns a CompareResult.
func CompareVaults(a, b *Vault) *CompareResult {
	result := &CompareResult{}

	aKeys := make(map[string]string, len(a.Entries))
	for _, e := range a.Entries {
		aKeys[e.Key] = e.Value
	}

	bKeys := make(map[string]string, len(b.Entries))
	for _, e := range b.Entries {
		bKeys[e.Key] = e.Value
	}

	for k, va := range aKeys {
		if vb, ok := bKeys[k]; !ok {
			result.OnlyInA = append(result.OnlyInA, k)
		} else if va != vb {
			result.Changed = append(result.Changed, k)
		} else {
			result.Identical = append(result.Identical, k)
		}
	}

	for k := range bKeys {
		if _, ok := aKeys[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Changed)
	sort.Strings(result.Identical)

	return result
}
