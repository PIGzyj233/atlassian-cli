package sprint

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdSprint creates the `sprint` subcommand group.
func NewCmdSprint(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sprint <command>",
		Short: "Manage agile sprints",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdIssues(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdAddIssues(f))
	return cmd
}
