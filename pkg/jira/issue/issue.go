package issue

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdIssue creates the `issue` subcommand group.
func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage Jira issues",
	}

	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))

	return cmd
}
