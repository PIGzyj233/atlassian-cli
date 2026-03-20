package project

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `project list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var includeArchived bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Jira projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listProjects(client, includeArchived)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().BoolVar(&includeArchived, "include-archived", false, "Include archived projects")
	return cmd
}

func listProjects(client *api.Client, includeArchived bool) (map[string]any, error) {
	path := client.JiraAPIPath("project")

	rb := api.NewRequestBuilder(path)
	if includeArchived {
		rb.Query("includeArchived", "true")
	}

	var result []any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}
	return map[string]any{"projects": result}, nil
}
