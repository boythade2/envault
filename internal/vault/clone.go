package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CloneResult summarises what was copied during a clone operation.
type CloneResult struct {
	Source      string
	Destination string
	EntriesCount int
}

// CloneVault copies all entries from the vault at srcPath into a new vault
// file at dstPath, re-encrypting every value with newPassphrase.
// The destination file must not already exist.
func CloneVault(srcPath, dstPath, oldPassphrase, newPassphrase string) (CloneResult, error) {
	if _, err := os.Stat(dstPath); err == nil {
		return CloneResult{}, fmt.Errorf("destination already exists: %s", dstPath)
	}

	src, err := LoadOrCreate(srcPath, oldPassphrase)
	if err != nil {
		return CloneResult{}, fmt.Errorf("load source vault: %w", err)
	}

	// Decrypt all values from source, then re-encrypt with new passphrase.
	dst, err := LoadOrCreate(dstPath, newPassphrase)
	if err != nil {
		return CloneResult{}, fmt.Errorf("create destination vault: %w", err)
	}

	for key, entry := range src.Entries {
		plain, err := src.GetDecrypted(key, oldPassphrase)
		if err != nil {
			return CloneResult{}, fmt.Errorf("decrypt key %q: %w", key, err)
		}
		if err := dst.AddEncrypted(key, plain, newPassphrase); err != nil {
			return CloneResult{}, fmt.Errorf("encrypt key %q: %w", key, err)
		}
		// Preserve tags if any.
		if len(entry.Tags) > 0 {
			dst.Entries[key] = func() Entry {
				e := dst.Entries[key]
				e.Tags = append([]string(nil), entry.Tags...)
				return e
			}()
		}
	}

	if err := dst.Save(dstPath, newPassphrase); err != nil {
		return CloneResult{}, fmt.Errorf("save destination vault: %w", err)
	}

	return CloneResult{
		Source:       filepath.Clean(srcPath),
		Destination:  filepath.Clean(dstPath),
		EntriesCount: len(src.Entries),
	}, nil
}

// cloneMetaPath returns the sidecar JSON path used to record clone provenance.
func cloneMetaPath(vaultPath string) string {
	return vaultPath + ".clone.json"
}

// WriteCloneMeta persists a CloneResult as a small JSON sidecar next to the
// destination vault so users can trace the vault's origin.
func WriteCloneMeta(result CloneResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cloneMetaPath(result.Destination), data, 0600)
}
