package dev

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdInfo creates the `dev info` command.
func NewCmdInfo(f *cmdutil.Factory) *cobra.Command {
	var (
		appType  string
		dataType string
	)

	cmd := &cobra.Command{
		Use:   "info <issue-key>",
		Short: "Get development information for an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := fetchDevInfo(client, args[0], appType, dataType)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&appType, "app-type", "", "Application type (e.g. stash, GitHub)")
	cmd.Flags().StringVar(&dataType, "data-type", "", "Data type (e.g. repository, pullrequest, branch)")

	return cmd
}

// resolveIssueID resolves a Jira issue key to its numeric ID.
func resolveIssueID(client *api.Client, issueKey string) (string, error) {
	path := api.NewRequestBuilder(client.JiraAPIPath("issue/" + issueKey)).
		Query("fields", "id").
		BuildPath()

	var result map[string]any
	_, err := client.Get(path, &result)
	if err != nil {
		return "", fmt.Errorf("resolving issue ID for %s: %w", issueKey, err)
	}

	id, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected id type for issue %s", issueKey)
	}
	return id, nil
}

func fetchDevInfo(client *api.Client, issueKey, appType, dataType string) (map[string]any, error) {
	issueID, err := resolveIssueID(client, issueKey)
	if err != nil {
		return nil, err
	}

	rb := api.NewRequestBuilder("/rest/dev-status/latest/issue/detail").
		Query("issueId", issueID).
		Query("applicationType", appType).
		Query("dataType", dataType)

	var result map[string]any
	_, err = client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("fetching dev info for %s: %w", issueKey, err)
	}
	return result, nil
}
