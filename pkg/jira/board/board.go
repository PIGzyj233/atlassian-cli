package board

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdBoard creates the `board` subcommand group.
func NewCmdBoard(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "board <command>",
		Short: "Manage agile boards",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdIssues(f))
	return cmd
}
