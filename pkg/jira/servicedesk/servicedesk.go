package servicedesk

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdServiceDesk creates the `servicedesk` subcommand group.
func NewCmdServiceDesk(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servicedesk <command>",
		Short: "Manage Jira Service Desk",
	}
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdQueues(f))
	cmd.AddCommand(NewCmdQueueIssues(f))
	return cmd
}
