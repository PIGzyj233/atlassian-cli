package worklog

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAdd creates the `worklog add` command.
func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var (
		timeSpent string
		comment   string
		started   string
	)

	cmd := &cobra.Command{
		Use:   "add <issue-key>",
		Short: "Add a worklog entry to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := addWorklog(client, args[0], timeSpent, comment, started)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&timeSpent, "time-spent", "", "Time spent (e.g. \"2h 30m\") (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "Work description")
	cmd.Flags().StringVar(&started, "started", "", "Date/time work started (e.g. \"2024-01-15T09:00:00.000+0000\")")
	_ = cmd.MarkFlagRequired("time-spent")
	return cmd
}

func addWorklog(client *api.Client, issueKey, timeSpent, comment, started string) (map[string]any, error) {
	payload := map[string]any{
		"timeSpent": timeSpent,
	}
	if comment != "" {
		// Cloud (v3) requires ADF; Server/DC (v2) uses plain string.
		var commentValue any = comment
		if client.IsCloud() {
			commentValue = convert.MarkdownToADF(comment)
		}
		payload["comment"] = commentValue
	}
	if started != "" {
		payload["started"] = started
	}

	path := client.JiraAPIPath("issue/" + issueKey + "/worklog")
	var result map[string]any
	_, err := client.Post(path, payload, &result)
	if err != nil {
		return nil, fmt.Errorf("adding worklog to %s: %w", issueKey, err)
	}
	return result, nil
}
