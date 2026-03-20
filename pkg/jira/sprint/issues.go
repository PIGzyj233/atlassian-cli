package sprint

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdIssues creates the `sprint issues` command.
func NewCmdIssues(f *cmdutil.Factory) *cobra.Command {
	var (
		fields string
		limit  int
	)

	cmd := &cobra.Command{
		Use:   "issues <sprint-id>",
		Short: "List issues in a sprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listSprintIssues(client, args[0], fields, limit)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated fields to return")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	return cmd
}

func listSprintIssues(client *api.Client, sprintID, fields string, limit int) (map[string]any, error) {
	rb := api.NewRequestBuilder(fmt.Sprintf("/rest/agile/1.0/sprint/%s/issue", sprintID)).
		Query("fields", fields)

	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("listing sprint %s issues: %w", sprintID, err)
	}
	return result, nil
}
