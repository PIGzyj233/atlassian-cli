package watcher

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAdd creates the `watcher add` command.
func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var user string

	cmd := &cobra.Command{
		Use:   "add <issue-key>",
		Short: "Add a watcher to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := addWatcher(client, args[0], user)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&user, "user", "", "User identifier (accountId for Cloud, username for Server)")
	_ = cmd.MarkFlagRequired("user")
	return cmd
}

func addWatcher(client *api.Client, issueKey, user string) (map[string]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey + "/watchers")

	// The Jira watchers POST endpoint expects a plain JSON string as the body,
	// not a JSON object. Passing the Go string directly to client.Post will
	// cause json.Marshal to produce the quoted string automatically (e.g. "user123").
	_, err := client.Post(path, user, nil)
	if err != nil {
		return nil, fmt.Errorf("adding watcher to %s: %w", issueKey, err)
	}

	return map[string]any{
		"success": true,
		"issue":   issueKey,
		"user":    user,
		"message": "Watcher added",
	}, nil
}
