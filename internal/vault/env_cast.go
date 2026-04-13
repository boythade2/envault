package vault

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents a supported type for casting an env value.
type CastType string

const (
	CastString CastType = "string"
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastBool   CastType = "bool"
)

// CastResult holds the outcome of a single cast operation.
type CastResult struct {
	Key     string
	OldVal  string
	NewVal  string
	Casted  bool
	Err     string
}

// CastEntries attempts to coerce entry values in the vault to the given type,
// normalising their string representation (e.g. "True" -> "true", "3.0" -> "3").
// If dryRun is true the vault file is not written.
func CastEntries(vaultPath, passphrase string, keys []string, castTo CastType, dryRun bool) ([]CastResult, error) {
	v, err := LoadOrCreate(vaultPath, passphrase)
	if err != nil {
		return nil, err
	}

	var results []CastResult

	for i, entry := range v.Entries {
		if !keyInList(entry.Key, keys) {
			continue
		}
		old := entry.Value
		newVal, castErr := castValue(old, castTo)
		r := CastResult{Key: entry.Key, OldVal: old}
		if castErr != nil {
			r.Err = castErr.Error()
		} else {
			r.NewVal = newVal
			r.Casted = newVal != old
			if !dryRun {
				v.Entries[i].Value = newVal
			}
		}
		results = append(results, r)
	}

	if !dryRun {
		if err := v.Save(vaultPath, passphrase); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func castValue(val string, to CastType) (string, error) {
	switch to {
	case CastInt:
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to int", val)
		}
		return strconv.FormatInt(int64(f), 10), nil
	case CastFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to float", val)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case CastBool:
		b, err := strconv.ParseBool(strings.TrimSpace(val))
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to bool", val)
		}
		return strconv.FormatBool(b), nil
	case CastString:
		return val, nil
	default:
		return "", fmt.Errorf("unknown cast type %q", to)
	}
}

func keyInList(key string, list []string) bool {
	for _, k := range list {
		if k == key {
			return true
		}
	}
	return false
}

// FormatCastResults returns a human-readable summary table.
func FormatCastResults(results []CastResult) string {
	if len(results) == 0 {
		return "no matching keys found\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-24s %-16s %-16s %s\n", "KEY", "OLD", "NEW", "STATUS"))
	for _, r := range results {
		status := "unchanged"
		newVal := r.OldVal
		if r.Err != "" {
			status = "error: " + r.Err
			newVal = "-"
		} else if r.Casted {
			status = "cast"
			newVal = r.NewVal
		}
		sb.WriteString(fmt.Sprintf("%-24s %-16s %-16s %s\n", r.Key, r.OldVal, newVal, status))
	}
	return sb.String()
}
