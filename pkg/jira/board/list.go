package board

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `board list` command.
// Ref: Python BoardsMixin.get_all_agile_boards() — uses Agile API prefix.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		name      string
		project   string
		boardType string
		limit     int
		startAt   int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List agile boards",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listBoards(client, name, project, boardType, limit, startAt)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Filter boards by name (fuzzy search)")
	cmd.Flags().StringVar(&project, "project", "", "Filter boards by project key")
	cmd.Flags().StringVar(&boardType, "type", "", "Filter boards by type (scrum, kanban)")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	cmd.Flags().IntVar(&startAt, "start-at", 0, "Starting index for pagination")
	return cmd
}

// listBoards uses the Agile API which has a fixed /rest/agile/1.0/ prefix
// for both Cloud and Server/DC.
func listBoards(client *api.Client, name, project, boardType string, limit, startAt int) (map[string]any, error) {
	rb := api.NewRequestBuilder("/rest/agile/1.0/board").
		Query("name", name).
		Query("projectKeyOrId", project).
		Query("type", boardType)

	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}
	if startAt > 0 {
		rb.QueryInt("startAt", startAt)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("listing boards: %w", err)
	}
	return result, nil
}
