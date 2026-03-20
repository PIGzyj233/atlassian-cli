package project

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdVersions creates the `project versions` command.
func NewCmdVersions(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions <project-key>",
		Short: "List versions for a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listVersions(client, args[0])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	return cmd
}

func listVersions(client *api.Client, projectKey string) (map[string]any, error) {
	path := client.JiraAPIPath(fmt.Sprintf("project/%s/versions", projectKey))

	var result []any
	_, err := client.Get(path, &result)
	if err != nil {
		return nil, fmt.Errorf("listing versions for project %s: %w", projectKey, err)
	}
	return map[string]any{"versions": result}, nil
}
