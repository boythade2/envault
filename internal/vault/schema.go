package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// SchemaRule defines a validation rule for a vault key.
type SchemaRule struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	Desc     string `json:"description,omitempty"`
}

// Schema holds a collection of rules for a vault.
type Schema struct {
	Rules []SchemaRule `json:"rules"`
}

// SchemaViolation describes a single rule violation.
type SchemaViolation struct {
	Key     string
	Message string
}

func schemaPath(vaultPath string) string {
	dir := filepath.Dir(vaultPath)
	base := filepath.Base(vaultPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	return filepath.Join(dir, name+".schema.json")
}

// LoadSchema reads the schema file associated with the given vault path.
// Returns an empty Schema and no error if the file does not exist.
func LoadSchema(vaultPath string) (Schema, error) {
	p := schemaPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return Schema{}, nil
	}
	if err != nil {
		return Schema{}, fmt.Errorf("read schema: %w", err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return Schema{}, fmt.Errorf("parse schema: %w", err)
	}
	return s, nil
}

// SaveSchema writes the schema to the file associated with the given vault path.
func SaveSchema(vaultPath string, s Schema) error {
	p := schemaPath(vaultPath)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	return os.WriteFile(p, data, 0600)
}

// ValidateSchema checks a vault's entries against the given schema and
// returns any violations found.
func ValidateSchema(v *Vault, s Schema) []SchemaViolation {
	var violations []SchemaViolation
	entryMap := make(map[string]string, len(v.Entries))
	for _, e := range v.Entries {
		entryMap[e.Key] = e.Value
	}
	for _, rule := range s.Rules {
		val, exists := entryMap[rule.Key]
		if rule.Required && !exists {
			violations = append(violations, SchemaViolation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}
		if exists && rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
			}
		}
	}
	return violations
}
