package comment

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdEdit creates the `comment edit` command.
func NewCmdEdit(f *cmdutil.Factory) *cobra.Command {
	var (
		body       string
		visibility string
	)

	cmd := &cobra.Command{
		Use:   "edit <issue-key> <comment-id>",
		Short: "Edit an issue comment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := editComment(client, args[0], args[1], body, visibility)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&body, "body", "", "New comment body (required)")
	cmd.Flags().StringVar(&visibility, "visibility", "", "Visibility role")
	_ = cmd.MarkFlagRequired("body")
	return cmd
}

func editComment(client *api.Client, issueKey, commentID, body, visibility string) (map[string]any, error) {
	payload := map[string]any{"body": body}
	if visibility != "" {
		payload["visibility"] = map[string]any{
			"type":  "role",
			"value": visibility,
		}
	}

	path := client.JiraAPIPath("issue/" + issueKey + "/comment/" + commentID)
	var result map[string]any
	_, err := client.Put(path, payload, &result)
	if err != nil {
		return nil, fmt.Errorf("editing comment %s on %s: %w", commentID, issueKey, err)
	}
	return result, nil
}
