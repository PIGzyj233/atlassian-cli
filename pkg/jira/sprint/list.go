package sprint

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `sprint list` command.
// Ref: Python SprintsMixin.get_all_sprints_from_board()
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		state string
		limit int
	)

	cmd := &cobra.Command{
		Use:   "list <board-id>",
		Short: "List sprints for a board",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listSprints(client, args[0], state, limit)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&state, "state", "", "Filter by sprint state (active, future, closed)")
	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	return cmd
}

func listSprints(client *api.Client, boardID, state string, limit int) (map[string]any, error) {
	rb := api.NewRequestBuilder(fmt.Sprintf("/rest/agile/1.0/board/%s/sprint", boardID)).
		Query("state", state)

	if limit > 0 {
		rb.QueryInt("maxResults", limit)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("listing sprints for board %s: %w", boardID, err)
	}
	return result, nil
}
