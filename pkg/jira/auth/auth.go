package auth

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAuth creates the `auth` subcommand group.
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Manage authentication",
	}

	cmd.AddCommand(NewCmdLogin(f))
	cmd.AddCommand(NewCmdLogout(f))
	cmd.AddCommand(NewCmdStatus(f))
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdSwitch(f))

	return cmd
}
