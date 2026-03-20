package board

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdIssues creates the `board issues` command.
func NewCmdIssues(f *cmdutil.Factory) *cobra.Command {
	var (
		jql     string
		fields  string
		limit   int
		startAt int
	)

	cmd := &cobra.Command{
		Use:   "issues <board-id>",
		Short: "List issues on a board",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listBoardIssues(client, args[0], jql, fields, limit, startAt)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&jql, "jql", "", "JQL filter for issues")
	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated fields to return")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	cmd.Flags().IntVar(&startAt, "start-at", 0, "Starting index for pagination")
	return cmd
}

func listBoardIssues(client *api.Client, boardID, jql, fields string, limit, startAt int) (map[string]any, error) {
	rb := api.NewRequestBuilder(fmt.Sprintf("/rest/agile/1.0/board/%s/issue", boardID)).
		Query("jql", jql).
		Query("fields", fields)

	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}
	if startAt > 0 {
		rb.QueryInt("startAt", startAt)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("listing board %s issues: %w", boardID, err)
	}
	return result, nil
}
