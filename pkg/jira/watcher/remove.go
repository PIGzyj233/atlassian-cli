package watcher

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdRemove creates the `watcher remove` command.
func NewCmdRemove(f *cmdutil.Factory) *cobra.Command {
	var (
		username  string
		accountID string
	)

	cmd := &cobra.Command{
		Use:   "remove <issue-key>",
		Short: "Remove a watcher from an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if username == "" && accountID == "" {
				return fmt.Errorf("at least one of --username or --account-id must be provided")
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := removeWatcher(client, args[0], username, accountID)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "Username (for Server/DC)")
	cmd.Flags().StringVar(&accountID, "account-id", "", "Account ID (for Cloud)")
	return cmd
}

func removeWatcher(client *api.Client, issueKey, username, accountID string) (map[string]any, error) {
	rb := api.NewRequestBuilder(client.JiraAPIPath("issue/" + issueKey + "/watchers"))
	rb.Query("username", username)
	rb.Query("accountId", accountID)

	_, err := client.Delete(rb.BuildPath(), nil)
	if err != nil {
		return nil, fmt.Errorf("removing watcher from %s: %w", issueKey, err)
	}

	user := username
	if accountID != "" {
		user = accountID
	}

	return map[string]any{
		"success": true,
		"issue":   issueKey,
		"user":    user,
		"message": "Watcher removed",
	}, nil
}
