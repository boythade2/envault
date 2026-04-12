package vault

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateResult holds the rendered output and any missing keys.
type TemplateResult struct {
	Rendered string
	Missing  []string
}

// placeholder matches {{KEY}} or {{ KEY }} style tokens.
var placeholder = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_]+)\s*\}\}`)

// RenderTemplate reads a template file and substitutes {{KEY}} tokens with
// values from the vault. Missing keys are collected rather than causing an
// error so the caller can decide how to handle them.
func RenderTemplate(templatePath string, v *Vault) (*TemplateResult, error) {
	raw, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}

	seen := map[string]bool{}
	var missing []string

	rendered := placeholder.ReplaceAllStringFunc(string(raw), func(match string) string {
		subs := placeholder.FindStringSubmatch(match)
		if len(subs) < 2 {
			return match
		}
		key := strings.TrimSpace(subs[1])
		if entry, ok := v.Entries[key]; ok {
			return entry.Value
		}
		if !seen[key] {
			missing = append(missing, key)
			seen[key] = true
		}
		return match
	})

	return &TemplateResult{
		Rendered: rendered,
		Missing:  missing,
	}, nil
}
