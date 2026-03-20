package worklog

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdWorklog creates the `worklog` subcommand group.
func NewCmdWorklog(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worklog <command>",
		Short: "Manage issue worklogs",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdAdd(f))
	return cmd
}
