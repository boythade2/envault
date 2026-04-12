package vault

import (
	"fmt"
	"time"
)

// CopyResult describes the outcome of a single key copy operation.
type CopyResult struct {
	SourceKey string
	DestKey   string
	Overwrote bool
}

// CopyEntry copies a key from one vault file to another.
// If destKey is empty, the source key name is used.
// If overwrite is false and the key already exists in the destination, an error is returned.
func CopyEntry(srcPath, dstPath, srcKey, destKey, passphrase string, overwrite bool) (CopyResult, error) {
	if destKey == "" {
		destKey = srcKey
	}

	src, err := LoadOrCreate(srcPath, passphrase)
	if err != nil {
		return CopyResult{}, fmt.Errorf("open source vault: %w", err)
	}

	var found *Entry
	for i := range src.Entries {
		if src.Entries[i].Key == srcKey {
			found = &src.Entries[i]
			break
		}
	}
	if found == nil {
		return CopyResult{}, fmt.Errorf("key %q not found in source vault", srcKey)
	}

	dst, err := LoadOrCreate(dstPath, passphrase)
	if err != nil {
		return CopyResult{}, fmt.Errorf("open destination vault: %w", err)
	}

	overwrote := false
	for i := range dst.Entries {
		if dst.Entries[i].Key == destKey {
			if !overwrite {
				return CopyResult{}, fmt.Errorf("key %q already exists in destination vault (use --overwrite to replace)", destKey)
			}
			dst.Entries[i].Value = found.Value
			dst.Entries[i].Tags = found.Tags
			dst.Entries[i].UpdatedAt = time.Now().UTC()
			overwrote = true
			break
		}
	}

	if !overwrote {
		dst.Entries = append(dst.Entries, Entry{
			Key:       destKey,
			Value:     found.Value,
			Tags:      found.Tags,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
	}

	if err := dst.Save(dstPath, passphrase); err != nil {
		return CopyResult{}, fmt.Errorf("save destination vault: %w", err)
	}

	return CopyResult{SourceKey: srcKey, DestKey: destKey, Overwrote: overwrote}, nil
}
