package vault

import (
	"strings"
)

// SearchResult holds a matched entry along with its key.
type SearchResult struct {
	Key   string
	Value string
}

// Search looks for entries in the vault whose keys or values contain the
// given query string (case-insensitive). If searchValues is false, only
// keys are searched.
func (v *Vault) Search(query string, searchValues bool) []SearchResult {
	lower := strings.ToLower(query)
	var results []SearchResult

	for key, entry := range v.Entries {
		keyMatch := strings.Contains(strings.ToLower(key), lower)
		valMatch := searchValues && strings.Contains(strings.ToLower(entry.Value), lower)

		if keyMatch || valMatch {
			results = append(results, SearchResult{
				Key:   key,
				Value: entry.Value,
			})
		}
	}

	// Sort results by key for deterministic output.
	sortSearchResults(results)
	return results
}

// sortSearchResults sorts a slice of SearchResult by Key alphabetically.
func sortSearchResults(results []SearchResult) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Key < results[j-1].Key; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
