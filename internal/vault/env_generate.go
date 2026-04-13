package vault

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	charsetLower   = "abcdefghijklmnopqrstuvwxyz"
	charsetUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigits  = "0123456789"
	charsetSymbols = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// GenerateOptions controls how a secret value is generated.
type GenerateOptions struct {
	Length      int
	UseUpper    bool
	UseDigits   bool
	UseSymbols  bool
	DryRun      bool
}

// GenerateResult holds the outcome of a generate operation.
type GenerateResult struct {
	Key     string
	Value   string
	Created bool
}

// GenerateSecret produces a random string based on GenerateOptions.
func GenerateSecret(opts GenerateOptions) (string, error) {
	if opts.Length <= 0 {
		return "", fmt.Errorf("length must be greater than zero")
	}
	charset := charsetLower
	if opts.UseUpper {
		charset += charsetUpper
	}
	if opts.UseDigits {
		charset += charsetDigits
	}
	if opts.UseSymbols {
		charset += charsetSymbols
	}
	var sb strings.Builder
	for i := 0; i < opts.Length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("random generation failed: %w", err)
		}
		sb.WriteByte(charset[idx.Int64()])
	}
	return sb.String(), nil
}

// GenerateAndStore generates a secret and stores it in the vault under key.
func GenerateAndStore(vaultFile, key string, opts GenerateOptions) (GenerateResult, error) {
	v, err := LoadOrCreate(vaultFile)
	if err != nil {
		return GenerateResult{}, err
	}
	secret, err := GenerateSecret(opts)
	if err != nil {
		return GenerateResult{}, err
	}
	result := GenerateResult{Key: key, Value: secret, Created: true}
	if opts.DryRun {
		return result, nil
	}
	v.Entries[key] = Entry{Value: secret, UpdatedAt: now()}
	return result, v.Save(vaultFile)
}

// FormatGenerateResult formats a GenerateResult for CLI display.
func FormatGenerateResult(r GenerateResult, dryRun bool) string {
	prefix := ""
	if dryRun {
		prefix = "[dry-run] "
	}
	return fmt.Sprintf("%sgenerated %s=%s", prefix, r.Key, r.Value)
}
