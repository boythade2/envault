package vault

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// EnvDiffResult holds the comparison between vault entries and current OS environment.
type EnvDiffResult struct {
	OnlyInVault []string // keys present in vault but not in environment
	OnlyInEnv   []string // keys present in environment but not in vault (filtered by prefix)
	Mismatch    []string // keys present in both but with different values
	Match       []string // keys present in both with identical values
}

// DiffEnv compares vault entries against the current OS environment.
// If prefix is non-empty, only environment variables with that prefix are considered.
func DiffEnv(v *Vault, prefix string) EnvDiffResult {
	result := EnvDiffResult{}

	vaultKeys := make(map[string]string, len(v.Entries))
	for k, e := range v.Entries {
		vaultKeys[k] = e.Value
	}

	envKeys := make(map[string]string)
	for _, raw := range os.Environ() {
		parts := strings.SplitN(raw, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k, val := parts[0], parts[1]
		if prefix == "" || strings.HasPrefix(k, prefix) {
			key := k
			if prefix != "" {
				key = strings.TrimPrefix(k, prefix)
			}
			envKeys[key] = val
		}
	}

	for k, vv := range vaultKeys {
		ev, ok := envKeys[k]
		if !ok {
			result.OnlyInVault = append(result.OnlyInVault, k)
		} else if ev != vv {
			result.Mismatch = append(result.Mismatch, k)
		} else {
			result.Match = append(result.Match, k)
		}
	}

	for k := range envKeys {
		if _, ok := vaultKeys[k]; !ok {
			result.OnlyInEnv = append(result.OnlyInEnv, k)
		}
	}

	sort.Strings(result.OnlyInVault)
	sort.Strings(result.OnlyInEnv)
	sort.Strings(result.Mismatch)
	sort.Strings(result.Match)

	return result
}

// FormatEnvDiffResult returns a human-readable summary of an EnvDiffResult.
func FormatEnvDiffResult(r EnvDiffResult) string {
	var sb strings.Builder

	if len(r.Mismatch) > 0 {
		sb.WriteString("MISMATCH (vault vs env differ):\n")
		for _, k := range r.Mismatch {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
		}
	}
	if len(r.OnlyInVault) > 0 {
		sb.WriteString("ONLY IN VAULT (not set in environment):\n")
		for _, k := range r.OnlyInVault {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}
	if len(r.OnlyInEnv) > 0 {
		sb.WriteString("ONLY IN ENV (not in vault):\n")
		for _, k := range r.OnlyInEnv {
			sb.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}
	if len(r.Match) > 0 {
		sb.WriteString(fmt.Sprintf("MATCH: %d key(s) in sync\n", len(r.Match)))
	}
	if sb.Len() == 0 {
		return "vault and environment are identical\n"
	}
	return sb.String()
}
