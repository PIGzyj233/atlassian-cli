package sprint

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAddIssues creates the `sprint add-issues` command.
// Ref: Python SprintsMixin.add_issues_to_sprint()
func NewCmdAddIssues(f *cmdutil.Factory) *cobra.Command {
	var keys string

	cmd := &cobra.Command{
		Use:   "add-issues <sprint-id>",
		Short: "Add issues to a sprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse comma-separated keys
			issueKeys := parseKeys(keys)
			if len(issueKeys) == 0 {
				return fmt.Errorf("at least one issue key is required via --keys")
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// POST /rest/agile/1.0/sprint/{id}/issue
			// Body: {"issues": ["KEY-1", "KEY-2"]}
			path := fmt.Sprintf("/rest/agile/1.0/sprint/%s/issue", args[0])
			body := map[string]any{
				"issues": issueKeys,
			}

			_, err = client.Post(path, body, nil)
			if err != nil {
				return fmt.Errorf("adding issues to sprint %s: %w", args[0], err)
			}

			return output.Print(cmd, map[string]any{
				"success":  true,
				"sprintId": args[0],
				"added":    issueKeys,
				"count":    len(issueKeys),
			})
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "Comma-separated issue keys to add (required)")
	_ = cmd.MarkFlagRequired("keys")
	return cmd
}

func parseKeys(keys string) []string {
	var result []string
	for _, k := range strings.Split(keys, ",") {
		trimmed := strings.TrimSpace(k)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
