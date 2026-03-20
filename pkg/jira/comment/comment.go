package comment

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdComment creates the `comment` subcommand group.
func NewCmdComment(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment <command>",
		Short: "Manage issue comments",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdAdd(f))
	cmd.AddCommand(NewCmdEdit(f))
	return cmd
}
