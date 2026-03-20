package page

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdPage creates the `page` subcommand group.
func NewCmdPage(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page <command>",
		Short: "Manage Confluence pages",
	}
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdMove(f))
	cmd.AddCommand(NewCmdChildren(f))
	cmd.AddCommand(NewCmdTree(f))
	cmd.AddCommand(NewCmdHistory(f))
	cmd.AddCommand(NewCmdDiff(f))
	return cmd
}
