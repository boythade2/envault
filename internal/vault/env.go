package vault

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// EnvOverrideResult holds the result of applying environment overrides to a vault.
type EnvOverrideResult struct {
	Applied  []string
	Skipped  []string
	NotFound []string
}

// ApplyEnvOverrides reads matching OS environment variables and overlays their
// values onto vault entries. Only keys already present in the vault are updated
// unless allowNew is true.
func ApplyEnvOverrides(v *Vault, prefix string, allowNew bool) EnvOverrideResult {
	result := EnvOverrideResult{}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		envKey, envVal := parts[0], parts[1]

		if prefix != "" && !strings.HasPrefix(envKey, prefix) {
			continue
		}

		vaultKey := envKey
		if prefix != "" {
			vaultKey = strings.TrimPrefix(envKey, prefix)
		}

		if _, exists := v.Entries[vaultKey]; exists {
			e := v.Entries[vaultKey]
			e.Value = envVal
			v.Entries[vaultKey] = e
			result.Applied = append(result.Applied, vaultKey)
		} else if allowNew {
			v.Entries[vaultKey] = Entry{Value: envVal}
			result.Applied = append(result.Applied, vaultKey)
		} else {
			result.NotFound = append(result.NotFound, vaultKey)
		}
	}

	sort.Strings(result.Applied)
	sort.Strings(result.NotFound)
	return result
}

// FormatEnvOverrideResult returns a human-readable summary of an override result.
func FormatEnvOverrideResult(r EnvOverrideResult) string {
	var sb strings.Builder
	if len(r.Applied) == 0 {
		sb.WriteString("No environment variables applied.\n")
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("Applied (%d):\n", len(r.Applied)))
	for _, k := range r.Applied {
		sb.WriteString(fmt.Sprintf("  + %s\n", k))
	}
	if len(r.NotFound) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped — key not in vault (%d):\n", len(r.NotFound)))
		for _, k := range r.NotFound {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
		}
	}
	return sb.String()
}
