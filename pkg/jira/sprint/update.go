package sprint

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// validSprintStates are the allowed sprint state values.
// Ref: Python update_sprint() validates state ∈ {future, active, closed}.
var validSprintStates = map[string]bool{
	"future": true,
	"active": true,
	"closed": true,
}

// NewCmdUpdate creates the `sprint update` command.
// Ref: Python SprintsMixin.update_sprint() — partial update with state validation.
func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		name      string
		state     string
		startDate string
		endDate   string
		goal      string
	)

	cmd := &cobra.Command{
		Use:   "update <sprint-id>",
		Short: "Update an existing sprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate state if provided
			if state != "" && !validSprintStates[state] {
				return fmt.Errorf("invalid state %q; valid states are: future, active, closed", state)
			}

			// Ensure at least one update field is provided
			if name == "" && state == "" && startDate == "" && endDate == "" && goal == "" {
				return fmt.Errorf("at least one of --name, --state, --start-date, --end-date, or --goal is required")
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := updateSprint(client, args[0], name, state, startDate, endDate, goal)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New sprint name")
	cmd.Flags().StringVar(&state, "state", "", "New sprint state (future, active, closed)")
	cmd.Flags().StringVar(&startDate, "start-date", "", "New start date (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&endDate, "end-date", "", "New end date (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&goal, "goal", "", "New sprint goal")
	return cmd
}

// updateSprint performs a partial update. Only sends changed fields.
// Ref: Python update_sprint() uses update_partially_sprint().
func updateSprint(client *api.Client, sprintID, name, state, startDate, endDate, goal string) (map[string]any, error) {
	body := make(map[string]any)
	if name != "" {
		body["name"] = name
	}
	if state != "" {
		body["state"] = state
	}
	if startDate != "" {
		body["startDate"] = startDate
	}
	if endDate != "" {
		body["endDate"] = endDate
	}
	if goal != "" {
		body["goal"] = goal
	}

	path := fmt.Sprintf("/rest/agile/1.0/sprint/%s", sprintID)
	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("updating sprint %s: %w", sprintID, err)
	}
	return result, nil
}
