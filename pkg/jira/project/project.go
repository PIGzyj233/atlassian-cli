package project

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdProject creates the `project` subcommand group.
func NewCmdProject(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <command>",
		Short: "Manage Jira projects",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdVersions(f))
	cmd.AddCommand(NewCmdComponents(f))
	cmd.AddCommand(NewCmdVersionCreate(f))

	return cmd
}
