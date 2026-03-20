package link

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdRemove creates the `link remove` command.
func NewCmdRemove(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <link-id>",
		Short: "Remove an issue link",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if err := removeLink(client, args[0]); err != nil {
				return err
			}

			return output.Print(cmd, map[string]any{
				"success": true,
				"message": "Issue link removed",
			})
		},
	}
	return cmd
}

func removeLink(client *api.Client, linkID string) error {
	path := client.JiraAPIPath("issueLink/" + linkID)
	_, err := client.Delete(path, nil)
	if err != nil {
		return fmt.Errorf("removing issue link %s: %w", linkID, err)
	}
	return nil
}
