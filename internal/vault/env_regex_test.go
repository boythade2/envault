package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildRegexVault(t *testing.T) (*Vault, string) {
	t.Helper()
	dir := t.TempDir()
	v := &Vault{Path: dir + "/test.vault"}
	v.Entries = []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_URL", Value: "postgres://localhost/mydb"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
	return v, dir
}

func TestRegexFilterByKey(t *testing.T) {
	v, _ := buildRegexVault(t)
	results, err := RegexFilter(v, "^APP_", false)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	keys := []string{results[0].Key, results[1].Key}
	assert.Contains(t, keys, "APP_HOST")
	assert.Contains(t, keys, "APP_PORT")
}

func TestRegexFilterByValue(t *testing.T) {
	v, _ := buildRegexVault(t)
	results, err := RegexFilter(v, "localhost", true)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestRegexFilterNoMatch(t *testing.T) {
	v, _ := buildRegexVault(t)
	results, err := RegexFilter(v, "^NONEXISTENT", false)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestRegexFilterEmptyPatternReturnsError(t *testing.T) {
	v, _ := buildRegexVault(t)
	_, err := RegexFilter(v, "", false)
	assert.Error(t, err)
}

func TestRegexFilterInvalidPatternReturnsError(t *testing.T) {
	v, _ := buildRegexVault(t)
	_, err := RegexFilter(v, "[invalid", false)
	assert.Error(t, err)
}

func TestFormatRegexResultsEmpty(t *testing.T) {
	out := FormatRegexResults(nil)
	assert.Equal(t, "no entries matched", out)
}

func TestFormatRegexResultsHasHeader(t *testing.T) {
	results := []RegexFilterResult{{Key: "FOO", Value: "bar", Matched: true}}
	out := FormatRegexResults(results)
	assert.Contains(t, out, "KEY")
	assert.Contains(t, out, "FOO")
	assert.Contains(t, out, "bar")
}
