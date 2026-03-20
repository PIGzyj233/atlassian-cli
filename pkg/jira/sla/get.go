package sla

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdGet creates the `sla get` command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var metrics string

	cmd := &cobra.Command{
		Use:   "get <issue-key>",
		Short: "Get SLA information for a service desk request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := getSLA(client, args[0], metrics)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&metrics, "metrics", "", "Comma-separated metric names to filter")

	return cmd
}

func getSLA(client *api.Client, issueKey, metrics string) (map[string]any, error) {
	// Service Desk API uses a fixed path, not the versioned Jira API path.
	path := "/rest/servicedeskapi/request/" + issueKey + "/sla"

	// ServiceDesk API requires X-ExperimentalApi header.
	headers := map[string]string{"X-ExperimentalApi": "opt-in"}

	var result map[string]any
	_, err := client.GetWithHeaders(path, &result, headers)
	if err != nil {
		return nil, fmt.Errorf("getting SLA for %s: %w", issueKey, err)
	}

	// Client-side filtering by metric names if --metrics is specified.
	if metrics != "" {
		result = filterMetrics(result, metrics)
	}

	return result, nil
}

// filterMetrics filters the SLA values by the requested metric names.
func filterMetrics(result map[string]any, metrics string) map[string]any {
	wanted := make(map[string]bool)
	for _, m := range strings.Split(metrics, ",") {
		wanted[strings.TrimSpace(m)] = true
	}

	values, ok := result["values"]
	if !ok {
		return result
	}

	valSlice, ok := values.([]any)
	if !ok {
		return result
	}

	var filtered []any
	for _, v := range valSlice {
		entry, ok := v.(map[string]any)
		if !ok {
			continue
		}
		name, _ := entry["name"].(string)
		if wanted[name] {
			filtered = append(filtered, entry)
		}
	}

	result["values"] = filtered
	return result
}
