package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportFormat defines the output format for vault exports.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
)

// Export writes vault entries to a file in the specified format.
func (v *Vault) Export(path string, format ExportFormat) error {
	var content string
	var err error

	switch format {
	case FormatDotenv:
		content, err = v.toDotenv()
	case FormatJSON:
		content, err = v.toJSON()
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to serialize vault: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// toDotenv serializes vault entries as KEY=VALUE lines.
func (v *Vault) toDotenv() (string, error) {
	var sb strings.Builder
	for _, entry := range v.Entries {
		line := fmt.Sprintf("%s=%s\n", entry.Key, entry.Value)
		sb.WriteString(line)
	}
	return sb.String(), nil
}

// toJSON serializes vault entries as a flat JSON object.
func (v *Vault) toJSON() (string, error) {
	m := make(map[string]string, len(v.Entries))
	for _, entry := range v.Entries {
		m[entry.Key] = entry.Value
	}
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes) + "\n", nil
}
