package worklog

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `worklog list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list <issue-key>",
		Short: "List worklogs on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listWorklogs(client, args[0], limit)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of worklogs to return")
	return cmd
}

func listWorklogs(client *api.Client, issueKey string, limit int) (map[string]any, error) {
	rb := api.NewRequestBuilder(client.JiraAPIPath("issue/" + issueKey + "/worklog"))
	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("listing worklogs for %s: %w", issueKey, err)
	}
	return result, nil
}
