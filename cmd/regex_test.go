package cmd

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/envault/internal/vault"
)

func writeRegexVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/test.vault"
	v := &vault.Vault{
		Path: path,
		Entries: []vault.Entry{
			{Key: "APP_HOST", Value: "localhost"},
			{Key: "APP_PORT", Value: "9000"},
			{Key: "DB_PASS", Value: "secret"},
		},
	}
	data, _ := json.Marshal(v)
	os.WriteFile(path, data, 0600)
	return path
}

func TestRegexCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "regex <vault-file> <pattern>" {
			found = true
		}
	}
	assert.True(t, found)
}

func TestRegexCommandRequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "regex", "only-one")
	assert.Error(t, err)
}

func TestRegexMatchesByKey(t *testing.T) {
	path := writeRegexVault(t)
	out, err := executeCommand(rootCmd, "regex", path, "^APP_")
	require.NoError(t, err)
	assert.Contains(t, out, "APP_HOST")
	assert.Contains(t, out, "APP_PORT")
	assert.NotContains(t, out, "DB_PASS")
}

func TestRegexNoMatchOutputsMessage(t *testing.T) {
	path := writeRegexVault(t)
	out, err := executeCommand(rootCmd, "regex", path, "^ZZZNOPE")
	require.NoError(t, err)
	assert.Contains(t, out, "no entries matched")
}

func TestRegexMatchValueFlag(t *testing.T) {
	path := writeRegexVault(t)
	out, err := executeCommand(rootCmd, "regex", path, "secret", "--match-value")
	require.NoError(t, err)
	assert.Contains(t, out, "DB_PASS")
}
