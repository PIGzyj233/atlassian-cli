package issue

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <issue-key>",
		Short: "Delete a Jira issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("use --confirm to delete issue %s", args[0])
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if err := deleteIssue(client, args[0]); err != nil {
				return err
			}

			return output.Print(cmd, map[string]any{
				"success": true,
				"key":     args[0],
				"message": "Issue deleted",
			})
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	return cmd
}

func deleteIssue(client *api.Client, issueKey string) error {
	path := client.JiraAPIPath("issue/" + issueKey)
	_, err := client.Delete(path, nil)
	if err != nil {
		return fmt.Errorf("deleting issue %s: %w", issueKey, err)
	}
	return nil
}
