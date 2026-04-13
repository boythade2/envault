package vault

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// VaultStats holds aggregate statistics about a vault file.
type VaultStats struct {
	TotalKeys      int
	EmptyValues    int
	UniqueValues   int
	DuplicateKeys  int
	OldestUpdated  time.Time
	NewestUpdated  time.Time
	TopPrefixes    []PrefixCount
}

// PrefixCount holds a key prefix and how many keys share it.
type PrefixCount struct {
	Prefix string
	Count  int
}

// ComputeStats returns aggregate statistics for the given vault.
func ComputeStats(v *Vault) VaultStats {
	if len(v.Entries) == 0 {
		return VaultStats{}
	}

	valueSet := make(map[string]int)
	prefixMap := make(map[string]int)

	var oldest, newest time.Time
	empty := 0

	for _, e := range v.Entries {
		valueSet[e.Value]++
		if strings.TrimSpace(e.Value) == "" {
			empty++
		}
		if oldest.IsZero() || e.UpdatedAt.Before(oldest) {
			oldest = e.UpdatedAt
		}
		if newest.IsZero() || e.UpdatedAt.After(newest) {
			newest = e.UpdatedAt
		}
		if idx := strings.Index(e.Key, "_"); idx > 0 {
			prefixMap[e.Key[:idx]]++
		}
	}

	uniqueVals := 0
	for _, count := range valueSet {
		if count == 1 {
			uniqueVals++
		}
	}

	prefixes := make([]PrefixCount, 0, len(prefixMap))
	for p, c := range prefixMap {
		prefixes = append(prefixes, PrefixCount{Prefix: p, Count: c})
	}
	sort.Slice(prefixes, func(i, j int) bool {
		if prefixes[i].Count != prefixes[j].Count {
			return prefixes[i].Count > prefixes[j].Count
		}
		return prefixes[i].Prefix < prefixes[j].Prefix
	})
	if len(prefixes) > 5 {
		prefixes = prefixes[:5]
	}

	return VaultStats{
		TotalKeys:     len(v.Entries),
		EmptyValues:   empty,
		UniqueValues:  uniqueVals,
		OldestUpdated: oldest,
		NewestUpdated: newest,
		TopPrefixes:   prefixes,
	}
}

// FormatStats returns a human-readable summary of vault statistics.
func FormatStats(s VaultStats) string {
	if s.TotalKeys == 0 {
		return "vault is empty — no statistics available\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Total keys      : %d\n", s.TotalKeys))
	sb.WriteString(fmt.Sprintf("Empty values    : %d\n", s.EmptyValues))
	sb.WriteString(fmt.Sprintf("Unique values   : %d\n", s.UniqueValues))
	sb.WriteString(fmt.Sprintf("Oldest update   : %s\n", s.OldestUpdated.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Newest update   : %s\n", s.NewestUpdated.Format(time.RFC3339)))
	if len(s.TopPrefixes) > 0 {
		sb.WriteString("Top prefixes    :\n")
		for _, p := range s.TopPrefixes {
			sb.WriteString(fmt.Sprintf("  %-20s %d keys\n", p.Prefix+"_", p.Count))
		}
	}
	return sb.String()
}
