package sprint

import (
	"fmt"
	"time"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdCreate creates the `sprint create` command.
// Ref: Python SprintsMixin.create_sprint() — validates dates.
func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		name      string
		startDate string
		endDate   string
		goal      string
	)

	cmd := &cobra.Command{
		Use:   "create <board-id>",
		Short: "Create a new sprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate start date if provided (ref: Python validates not in past)
			if startDate != "" {
				parsedStart, err := time.Parse(time.RFC3339, startDate)
				if err != nil {
					// Try date-only format
					parsedStart, err = time.Parse("2006-01-02", startDate)
					if err != nil {
						return fmt.Errorf("invalid --start-date format, use RFC3339 (2006-01-02T15:04:05Z) or date (2006-01-02)")
					}
				}
				if parsedStart.Before(time.Now()) {
					return fmt.Errorf("start date cannot be in the past")
				}

				// Validate end date > start date (ref: Python validates start < end)
				if endDate != "" {
					parsedEnd, err := time.Parse(time.RFC3339, endDate)
					if err != nil {
						parsedEnd, err = time.Parse("2006-01-02", endDate)
						if err != nil {
							return fmt.Errorf("invalid --end-date format, use RFC3339 (2006-01-02T15:04:05Z) or date (2006-01-02)")
						}
					}
					if !parsedStart.Before(parsedEnd) {
						return fmt.Errorf("start date must be before end date")
					}
				}
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := createSprint(client, args[0], name, startDate, endDate, goal)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Sprint name (required)")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Sprint start date (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&endDate, "end-date", "", "Sprint end date (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&goal, "goal", "", "Sprint goal")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func createSprint(client *api.Client, boardID, name, startDate, endDate, goal string) (map[string]any, error) {
	body := map[string]any{
		"name":          name,
		"originBoardId": boardID,
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

	var result map[string]any
	_, err := client.Post("/rest/agile/1.0/sprint", body, &result)
	if err != nil {
		return nil, fmt.Errorf("creating sprint: %w", err)
	}
	return result, nil
}
