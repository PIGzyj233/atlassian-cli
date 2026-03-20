package watcher

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `watcher list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <issue-key>",
		Short: "List watchers of an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listWatchers(client, args[0])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	return cmd
}

func listWatchers(client *api.Client, issueKey string) (map[string]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey + "/watchers")

	var result map[string]any
	_, err := client.Get(path, &result)
	if err != nil {
		return nil, fmt.Errorf("listing watchers for %s: %w", issueKey, err)
	}
	return result, nil
}
