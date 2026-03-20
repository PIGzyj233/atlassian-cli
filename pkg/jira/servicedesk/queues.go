package servicedesk

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdQueues creates the `servicedesk queues` command.
func NewCmdQueues(f *cmdutil.Factory) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "queues <desk-id>",
		Short: "List queues for a service desk",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listQueues(client, args[0], limit)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	return cmd
}

func listQueues(client *api.Client, deskID string, limit int) (map[string]any, error) {
	rb := api.NewRequestBuilder(fmt.Sprintf("/rest/servicedeskapi/servicedesk/%s/queue", deskID))

	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}

	// ServiceDesk API requires X-ExperimentalApi header.
	headers := map[string]string{"X-ExperimentalApi": "opt-in"}

	var result map[string]any
	_, err := client.GetWithHeaders(rb.BuildPath(), &result, headers)
	if err != nil {
		return nil, fmt.Errorf("listing queues for service desk %s: %w", deskID, err)
	}
	return result, nil
}
