package issue

import (
	"encoding/json"
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTransition(f *cmdutil.Factory) *cobra.Command {
	var (
		transitionID string
		comment      string
		fieldsJSON   string
	)

	cmd := &cobra.Command{
		Use:   "transition <issue-key>",
		Short: "Transition a Jira issue to a new status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			fields := make(map[string]any)
			if fieldsJSON != "" {
				if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
					return fmt.Errorf("parsing --fields-json: %w", err)
				}
			}

			if err := transitionIssue(client, args[0], transitionID, fields, comment); err != nil {
				return err
			}

			return output.Print(cmd, map[string]any{
				"success":      true,
				"key":          args[0],
				"transitionId": transitionID,
			})
		},
	}

	cmd.Flags().StringVar(&transitionID, "transition-id", "", "Target transition ID (required)")
	cmd.Flags().StringVar(&fieldsJSON, "fields-json", "", "Fields to set during transition as JSON")
	cmd.Flags().StringVar(&comment, "comment", "", "Comment to add during transition")
	_ = cmd.MarkFlagRequired("transition-id")

	return cmd
}

func transitionIssue(client *api.Client, issueKey, transitionID string, fields map[string]any, comment string) error {
	body := map[string]any{
		"transition": map[string]any{"id": transitionID},
	}

	if len(fields) > 0 {
		body["fields"] = fields
	}
	if comment != "" {
		body["update"] = map[string]any{
			"comment": []any{
				map[string]any{
					"add": map[string]any{"body": comment},
				},
			},
		}
	}

	path := client.JiraAPIPath("issue/" + issueKey + "/transitions")
	_, err := client.Post(path, body, nil)
	if err != nil {
		return fmt.Errorf("transitioning issue %s: %w", issueKey, err)
	}
	return nil
}
