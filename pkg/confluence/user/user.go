package user

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdUser creates the `user` subcommand group.
func NewCmdUser(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user <command>",
		Short: "Manage Confluence users",
	}
	cmd.AddCommand(NewCmdSearch(f))
	return cmd
}
