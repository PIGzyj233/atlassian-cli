package dev

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdBatchInfo creates the `dev batch-info` command.
func NewCmdBatchInfo(f *cmdutil.Factory) *cobra.Command {
	var (
		keys     string
		appType  string
		dataType string
	)

	cmd := &cobra.Command{
		Use:   "batch-info",
		Short: "Get development information for multiple issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			issueKeys := strings.Split(keys, ",")
			for i := range issueKeys {
				issueKeys[i] = strings.TrimSpace(issueKeys[i])
			}

			result, err := fetchBatchDevInfo(client, issueKeys, appType, dataType)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "Comma-separated list of issue keys (required)")
	_ = cmd.MarkFlagRequired("keys")
	cmd.Flags().StringVar(&appType, "app-type", "", "Application type (e.g. stash, GitHub)")
	cmd.Flags().StringVar(&dataType, "data-type", "", "Data type (e.g. repository, pullrequest, branch)")

	return cmd
}

func fetchBatchDevInfo(client *api.Client, issueKeys []string, appType, dataType string) (map[string]any, error) {
	results := make(map[string]any, len(issueKeys))

	for _, key := range issueKeys {
		issueID, err := resolveIssueID(client, key)
		if err != nil {
			results[key] = map[string]any{"error": fmt.Sprintf("resolving issue: %v", err)}
			continue
		}

		rb := api.NewRequestBuilder("/rest/dev-status/latest/issue/detail").
			Query("issueId", issueID).
			Query("applicationType", appType).
			Query("dataType", dataType)

		var result map[string]any
		_, err = client.Get(rb.BuildPath(), &result)
		if err != nil {
			results[key] = map[string]any{"error": fmt.Sprintf("fetching dev info: %v", err)}
			continue
		}

		results[key] = result
	}

	return results, nil
}
