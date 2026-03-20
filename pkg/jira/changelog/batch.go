package changelog

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdBatch creates the `changelog batch` command.
func NewCmdBatch(f *cmdutil.Factory) *cobra.Command {
	var (
		keys   string
		fields string
		limit  int
	)

	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Fetch changelogs for multiple issues (Cloud only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := batchChangelog(client, keys, fields, limit)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of issue keys (required)")
	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated field names to filter changelog items")
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum results per issue")
	_ = cmd.MarkFlagRequired("keys")

	return cmd
}

func batchChangelog(client *api.Client, keys, fields string, limit int) (map[string]any, error) {
	if !client.IsCloud() {
		return nil, fmt.Errorf("changelog batch is only supported on Jira Cloud")
	}

	issueKeys := parseCSV(keys)
	fieldFilter := parseCSV(fields)

	results := make(map[string]any, len(issueKeys))

	for _, key := range issueKeys {
		changelog, err := fetchChangelog(client, key, limit)
		if err != nil {
			results[key] = map[string]any{"error": fmt.Sprintf("fetching changelog: %v", err)}
			continue
		}

		if len(fieldFilter) > 0 {
			changelog = filterByFields(changelog, fieldFilter)
		}

		results[key] = changelog
	}

	return results, nil
}

func fetchChangelog(client *api.Client, issueKey string, limit int) (map[string]any, error) {
	path := client.JiraAPIPath(fmt.Sprintf("issue/%s/changelog", issueKey))

	rb := api.NewRequestBuilder(path)
	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("fetching changelog for %s: %w", issueKey, err)
	}
	return result, nil
}

// filterByFields filters each history entry's items to include only those
// matching the requested field names.
func filterByFields(changelog map[string]any, fieldFilter []string) map[string]any {
	wanted := make(map[string]bool, len(fieldFilter))
	for _, f := range fieldFilter {
		wanted[f] = true
	}

	values, ok := changelog["values"]
	if !ok {
		return changelog
	}

	valSlice, ok := values.([]any)
	if !ok {
		return changelog
	}

	var filteredHistories []any
	for _, v := range valSlice {
		history, ok := v.(map[string]any)
		if !ok {
			continue
		}

		items, ok := history["items"]
		if !ok {
			continue
		}

		itemSlice, ok := items.([]any)
		if !ok {
			continue
		}

		var filteredItems []any
		for _, item := range itemSlice {
			entry, ok := item.(map[string]any)
			if !ok {
				continue
			}
			fieldName, _ := entry["field"].(string)
			if wanted[fieldName] {
				filteredItems = append(filteredItems, entry)
			}
		}

		if len(filteredItems) > 0 {
			historyCopy := make(map[string]any, len(history))
			for k, v := range history {
				historyCopy[k] = v
			}
			historyCopy["items"] = filteredItems
			filteredHistories = append(filteredHistories, historyCopy)
		}
	}

	changelog["values"] = filteredHistories
	return changelog
}

// parseCSV splits a comma-separated string into trimmed, non-empty tokens.
func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
