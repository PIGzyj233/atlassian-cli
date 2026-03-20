package field

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdSearch creates the `field search` command.
// Ref: Python FieldsMixin.search_fields() — fetches all fields then filters.
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var (
		keyword string
		limit   int
		refresh bool
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for Jira fields by keyword",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := searchFields(client, keyword, limit)
			if err != nil {
				return err
			}

			_ = refresh // reserved for future caching
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&keyword, "keyword", "", "Search keyword (case-insensitive substring match)")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to return")
	cmd.Flags().BoolVar(&refresh, "refresh", false, "Force refresh from server (bypass cache)")
	return cmd
}

func searchFields(client *api.Client, keyword string, limit int) ([]any, error) {
	// Fetch all fields via GET /rest/api/{v}/field
	path := client.JiraAPIPath("field")
	var allFields []any
	_, err := client.Get(path, &allFields)
	if err != nil {
		return nil, fmt.Errorf("fetching fields: %w", err)
	}

	// If no keyword, return first `limit` fields
	if keyword == "" {
		if limit > 0 && limit < len(allFields) {
			return allFields[:limit], nil
		}
		return allFields, nil
	}

	// Filter by keyword (case-insensitive substring match on id, key, name, clauseNames)
	// Ref: Python search_fields() uses fuzzy matching; we use simpler substring matching.
	lowerKeyword := strings.ToLower(keyword)
	var matched []any
	for _, f := range allFields {
		fm, ok := f.(map[string]any)
		if !ok {
			continue
		}
		if matchesKeyword(fm, lowerKeyword) {
			matched = append(matched, fm)
			if limit > 0 && len(matched) >= limit {
				break
			}
		}
	}

	return matched, nil
}

func matchesKeyword(field map[string]any, lowerKeyword string) bool {
	// Check id
	if id, _ := field["id"].(string); strings.Contains(strings.ToLower(id), lowerKeyword) {
		return true
	}
	// Check key
	if key, _ := field["key"].(string); strings.Contains(strings.ToLower(key), lowerKeyword) {
		return true
	}
	// Check name
	if name, _ := field["name"].(string); strings.Contains(strings.ToLower(name), lowerKeyword) {
		return true
	}
	// Check clauseNames array
	if clauseNames, ok := field["clauseNames"].([]any); ok {
		for _, cn := range clauseNames {
			if s, ok := cn.(string); ok && strings.Contains(strings.ToLower(s), lowerKeyword) {
				return true
			}
		}
	}
	return false
}
