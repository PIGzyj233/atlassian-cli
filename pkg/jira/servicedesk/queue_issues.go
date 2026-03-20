package servicedesk

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdQueueIssues creates the `servicedesk queue-issues` command.
func NewCmdQueueIssues(f *cmdutil.Factory) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "queue-issues <desk-id> <queue-id>",
		Short: "List issues in a service desk queue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listQueueIssues(client, args[0], args[1], limit)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	return cmd
}

func listQueueIssues(client *api.Client, deskID, queueID string, limit int) (map[string]any, error) {
	rb := api.NewRequestBuilder(fmt.Sprintf("/rest/servicedeskapi/servicedesk/%s/queue/%s/issue", deskID, queueID))

	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}

	// ServiceDesk API requires X-ExperimentalApi header.
	headers := map[string]string{"X-ExperimentalApi": "opt-in"}

	var result map[string]any
	_, err := client.GetWithHeaders(rb.BuildPath(), &result, headers)
	if err != nil {
		return nil, fmt.Errorf("listing issues for queue %s in service desk %s: %w", queueID, deskID, err)
	}
	return result, nil
}
