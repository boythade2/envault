package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envault/internal/vault"
)

var (
	filterKeyPrefix     string
	filterKeySuffix     string
	filterValueContains string
	filterTags          string
	filterInvert        bool
)

var filterCmd = &cobra.Command{
	Use:   "filter <vault-file>",
	Short: "Filter vault entries by key prefix, suffix, value, or tag",
	Args:  cobra.ExactArgs(1),
	RunE:  runFilter,
}

func init() {
	filterCmd.Flags().StringVar(&filterKeyPrefix, "key-prefix", "", "Filter entries whose key starts with this prefix")
	filterCmd.Flags().StringVar(&filterKeySuffix, "key-suffix", "", "Filter entries whose key ends with this suffix")
	filterCmd.Flags().StringVar(&filterValueContains, "value-contains", "", "Filter entries whose value contains this substring")
	filterCmd.Flags().StringVar(&filterTags, "tags", "", "Comma-separated list of tags to filter by")
	filterCmd.Flags().BoolVar(&filterInvert, "invert", false, "Invert the filter match")
	rootCmd.AddCommand(filterCmd)
}

func runFilter(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]

	v, err := vault.LoadOrCreate(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	opts := vault.FilterOptions{
		KeyPrefix:     filterKeyPrefix,
		KeySuffix:     filterKeySuffix,
		ValueContains: filterValueContains,
		TagFilter:     parseTagList(filterTags),
		InvertMatch:   filterInvert,
	}

	if !hasFilterCriteria(opts) {
		return fmt.Errorf("at least one filter criterion must be specified (--key-prefix, --key-suffix, --value-contains, or --tags)")
	}

	results := vault.FilterEntries(v, opts)
	_, err = fmt.Fprint(os.Stdout, vault.FormatFilterResults(results))
	return err
}

// hasFilterCriteria reports whether at least one filter criterion has been set.
func hasFilterCriteria(opts vault.FilterOptions) bool {
	return opts.KeyPrefix != "" ||
		opts.KeySuffix != "" ||
		opts.ValueContains != "" ||
		len(opts.TagFilter) > 0
}

// parseTagList splits a comma-separated tag string into a trimmed, non-empty slice of tags.
func parseTagList(tags string) []string {
	if tags == "" {
		return nil
	}
	var tagList []string
	for _, t := range strings.Split(tags, ",") {
		if trimmed := strings.TrimSpace(t); trimmed != "" {
			tagList = append(tagList, trimmed)
		}
	}
	return tagList
}
