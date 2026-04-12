package vault

import "fmt"

// MergeStrategy defines how conflicting keys are resolved.
type MergeStrategy string

const (
	MergeStrategyOurs   MergeStrategy = "ours"   // keep value from base vault
	MergeStrategyTheirs MergeStrategy = "theirs" // take value from incoming vault
	MergeStrategyError  MergeStrategy = "error"  // return error on conflict
)

// MergeResult describes the outcome of a merge operation.
type MergeResult struct {
	Added    []string
	Updated  []string
	Skipped  []string
	Conflict []string
}

// MergeVaults merges entries from src into dst according to the given strategy.
// dst is modified in place and must be saved by the caller.
func MergeVaults(dst, src *Vault, strategy MergeStrategy) (MergeResult, error) {
	var result MergeResult

	for key, srcEntry := range src.Entries {
		if _, exists := dst.Entries[key]; !exists {
			dst.Entries[key] = srcEntry
			result.Added = append(result.Added, key)
			continue
		}

		// Conflict: key exists in both vaults.
		switch strategy {
		case MergeStrategyOurs:
			result.Skipped = append(result.Skipped, key)
		case MergeStrategyTheirs:
			dst.Entries[key] = srcEntry
			result.Updated = append(result.Updated, key)
		case MergeStrategyError:
			result.Conflict = append(result.Conflict, key)
		default:
			return result, fmt.Errorf("unknown merge strategy: %q", strategy)
		}
	}

	if len(result.Conflict) > 0 {
		return result, fmt.Errorf("merge conflict on keys: %v", result.Conflict)
	}

	return result, nil
}
