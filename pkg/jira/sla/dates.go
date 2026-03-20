package sla

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDates creates the `sla dates` command.
func NewCmdDates(f *cmdutil.Factory) *cobra.Command {
	var (
		includeStatusChanges bool
		includeStatusSummary bool
	)

	cmd := &cobra.Command{
		Use:   "dates <issue-key>",
		Short: "Get date-related fields for an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := getIssueDates(client, args[0], includeStatusChanges, includeStatusSummary)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().BoolVar(&includeStatusChanges, "include-status-changes", false, "Include status change history (adds expand=changelog)")
	cmd.Flags().BoolVar(&includeStatusSummary, "include-status-summary", false, "Include a summary of status durations")

	return cmd
}

func getIssueDates(client *api.Client, issueKey string, includeStatusChanges, includeStatusSummary bool) (map[string]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey)

	rb := api.NewRequestBuilder(path)
	if includeStatusChanges {
		rb.Query("expand", "changelog")
	}

	var issue map[string]any
	_, err := client.Get(rb.BuildPath(), &issue)
	if err != nil {
		return nil, fmt.Errorf("getting issue %s: %w", issueKey, err)
	}

	result := extractDates(issue, includeStatusChanges, includeStatusSummary)
	return result, nil
}

// dateFields are the standard Jira date-related field keys.
var dateFields = []string{
	"created",
	"updated",
	"resolutiondate",
	"duedate",
	"lastViewed",
	"statuscategorychangedate",
}

// extractDates pulls date-related fields from an issue response.
func extractDates(issue map[string]any, includeStatusChanges, includeStatusSummary bool) map[string]any {
	key, _ := issue["key"].(string)

	result := map[string]any{
		"key": key,
	}

	dates := make(map[string]any)

	fields, _ := issue["fields"].(map[string]any)
	if fields != nil {
		for _, f := range dateFields {
			if v, ok := fields[f]; ok && v != nil {
				dates[f] = v
			}
		}
	}

	result["dates"] = dates

	if includeStatusChanges {
		statusChanges := extractStatusChanges(issue)
		result["statusChanges"] = statusChanges

		if includeStatusSummary {
			result["statusSummary"] = buildStatusSummary(statusChanges)
		}
	}

	return result
}

// extractStatusChanges parses the changelog for status transitions.
func extractStatusChanges(issue map[string]any) []map[string]any {
	var changes []map[string]any

	changelog, _ := issue["changelog"].(map[string]any)
	if changelog == nil {
		return changes
	}

	histories, _ := changelog["histories"].([]any)
	for _, h := range histories {
		history, ok := h.(map[string]any)
		if !ok {
			continue
		}
		created, _ := history["created"].(string)
		items, _ := history["items"].([]any)
		for _, it := range items {
			item, ok := it.(map[string]any)
			if !ok {
				continue
			}
			field, _ := item["field"].(string)
			if field != "status" {
				continue
			}
			changes = append(changes, map[string]any{
				"timestamp":  created,
				"fromStatus": item["fromString"],
				"toStatus":   item["toString"],
			})
		}
	}

	return changes
}

// buildStatusSummary counts transitions per status.
func buildStatusSummary(changes []map[string]any) map[string]int {
	counts := make(map[string]int)
	for _, c := range changes {
		toStatus, _ := c["toStatus"].(string)
		if toStatus != "" {
			counts[toStatus]++
		}
	}
	return counts
}
